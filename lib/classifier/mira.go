package classifier

import (
	"fmt"
	"math"
	"math/rand"
	"os"
	"runtime"
	"sort"
	"sync"

	"github.com/pkg/errors"
	"github.com/syou6162/go-active-learning/lib/evaluation"
	"github.com/syou6162/go-active-learning/lib/feature"
	"github.com/syou6162/go-active-learning/lib/model"
	"github.com/syou6162/go-active-learning/lib/util"
)

type ModelType int

const (
	EXAMPLE ModelType = 0
	TWITTER ModelType = 1
)

type MIRAClassifier struct {
	ModelType ModelType          `json:"ModelType"`
	Weight    map[string]float64 `json:"Weight"`
	C         float64            `json:"C"`
	Accuracy  float64            `json:"Accuracy"`
	Precision float64            `json:"Precision"`
	Recall    float64            `json:"Recall"`
	Fvalue    float64            `json:"Fvalue"`
}

type LearningInstance interface {
	GetFeatureVector() feature.FeatureVector
	GetLabel() model.LabelType
}

type LearningInstances []LearningInstance

var errNoTrainingInstances = errors.New("Empty training set")
var errNoDevelopmentInstances = errors.New("Empty development set")
var errNoMIRAModelLearned = errors.New("Fail to learn MIRA models")
var errModelEvaluationFailure = errors.New("Failed to evaluate best MIRA")

func newMIRAClassifier(modelType ModelType, c float64) *MIRAClassifier {
	return &MIRAClassifier{
		ModelType: modelType,
		Weight:    make(map[string]float64),
		C:         c,
		Accuracy:  0.0,
		Precision: 0.0,
		Recall:    0.0,
		Fvalue:    0.0,
	}
}

func filterLabeledInstances(instances LearningInstances) LearningInstances {
	var result LearningInstances
	for _, i := range instances {
		if i.GetLabel() != 0 {
			result = append(result, i)
		}
	}
	return result
}

func shuffle(instances LearningInstances) {
	n := len(instances)
	for i := n - 1; i >= 0; i-- {
		j := rand.Intn(i + 1)
		instances[i], instances[j] = instances[j], instances[i]
	}
}

func splitTrainAndDev(instances LearningInstances) (train LearningInstances, dev LearningInstances) {
	shuffle(instances)
	n := int(0.8 * float64(len(instances)))
	return instances[0:n], instances[n:]
}

func NewMIRAClassifier(modelType ModelType, instances LearningInstances, c float64) *MIRAClassifier {
	train := filterLabeledInstances(instances)
	model := newMIRAClassifier(modelType, c)
	for iter := 0; iter < 30; iter++ {
		shuffle(train)
		for _, example := range train {
			model.learn(example)
		}
	}
	return model
}

func overSamplingPositiveExamples(instances LearningInstances) LearningInstances {
	overSampled := LearningInstances{}
	posInstances := LearningInstances{}
	negInstances := LearningInstances{}

	numNeg := 0

	for _, i := range instances {
		if i.GetLabel() == model.NEGATIVE {
			numNeg += 1
			negInstances = append(negInstances, i)
		} else if i.GetLabel() == model.POSITIVE {
			posInstances = append(posInstances, i)
		}
	}

	for len(overSampled) <= numNeg {
		shuffle(posInstances)
		overSampled = append(overSampled, posInstances[0])
	}
	overSampled = append(overSampled, negInstances...)
	shuffle(overSampled)

	return overSampled
}

func extractGoldLabels(instances LearningInstances) []model.LabelType {
	golds := make([]model.LabelType, 0, 0)
	for _, i := range instances {
		golds = append(golds, i.GetLabel())
	}
	return golds
}

type MIRAClassifierList []MIRAClassifier

func (l MIRAClassifierList) Len() int           { return len(l) }
func (l MIRAClassifierList) Less(i, j int) bool { return l[i].Fvalue < l[j].Fvalue }
func (l MIRAClassifierList) Swap(i, j int)      { l[i], l[j] = l[j], l[i] }

func isValidTrainAndDevelopmentInstances(train LearningInstances, dev LearningInstances) (bool, error) {
	if len(train) == 0 {
		return false, errNoTrainingInstances
	}
	if len(dev) == 0 {
		return false, errNoDevelopmentInstances
	}
	return true, nil
}

func NewMIRAClassifierByCrossValidation(modelType ModelType, instances LearningInstances) (*MIRAClassifier, error) {
	shuffle(instances)
	train, dev := splitTrainAndDev(filterLabeledInstances(instances))
	if valid, err := isValidTrainAndDevelopmentInstances(train, dev); !valid {
		return nil, err
	}

	train = overSamplingPositiveExamples(train)

	params := []float64{1000, 500, 100, 50, 10.0, 5.0, 1.0, 0.5, 0.1, 0.05, 0.01, 0.005, 0.001}
	miraResults := MIRAClassifierList{}

	wg := &sync.WaitGroup{}
	cpus := runtime.NumCPU()
	runtime.GOMAXPROCS(cpus)

	models := make([]*MIRAClassifier, len(params))
	for idx, c := range params {
		wg.Add(1)
		go func(idx int, c float64) {
			defer wg.Done()
			model := NewMIRAClassifier(modelType, train, c)
			models[idx] = model
		}(idx, c)
	}
	wg.Wait()

	if len(models) == 0 {
		return nil, errNoMIRAModelLearned
	}

	maxFvalue := math.Inf(-1)
	for _, m := range models {
		devPredicts := make([]model.LabelType, len(dev))
		for i, instance := range dev {
			devPredicts[i] = m.Predict(instance.GetFeatureVector())
		}
		m.Accuracy = evaluation.GetAccuracy(extractGoldLabels(dev), devPredicts)
		m.Precision = evaluation.GetPrecision(extractGoldLabels(dev), devPredicts)
		m.Recall = evaluation.GetRecall(extractGoldLabels(dev), devPredicts)
		m.Fvalue = (2 * m.Recall * m.Precision) / (m.Recall + m.Precision)
		fmt.Fprintln(os.Stderr, fmt.Sprintf("C:%0.03f\tAccuracy:%0.03f\tPrecision:%0.03f\tRecall:%0.03f\tF-value:%0.03f", m.C, m.Accuracy, m.Precision, m.Recall, m.Fvalue))
		if math.IsNaN(m.Fvalue) {
			continue
		}
		miraResults = append(miraResults, *m)
		if m.Fvalue >= maxFvalue {
			maxFvalue = m.Fvalue
		}
	}
	if len(miraResults) == 0 {
		return nil, errModelEvaluationFailure
	}

	sort.Sort(sort.Reverse(miraResults))
	bestModel := &miraResults[0]
	instances = overSamplingPositiveExamples(instances)
	shuffle(instances)
	result := NewMIRAClassifier(modelType, filterLabeledInstances(instances), bestModel.C)
	result.Accuracy = bestModel.Accuracy
	result.Precision = bestModel.Precision
	result.Recall = bestModel.Recall
	result.Fvalue = bestModel.Fvalue
	return result, nil
}

func (m *MIRAClassifier) learn(instance LearningInstance) {
	tmp := float64(instance.GetLabel()) * m.PredictScore(instance.GetFeatureVector()) // y w^T x
	loss := 0.0
	if tmp < 1.0 {
		loss = 1 - tmp
	}

	norm := float64(len(instance.GetFeatureVector()) * len(instance.GetFeatureVector()))
	// tau := math.Min(m.C, loss/norm) // update by PA-I
	tau := loss / (norm + 1.0/m.C) // update by PA-II

	if tau != 0.0 {
		for _, f := range instance.GetFeatureVector() {
			w, _ := m.Weight[f]
			m.Weight[f] = w + tau*float64(instance.GetLabel())
		}
	}
}

func (m MIRAClassifier) PredictScore(features feature.FeatureVector) float64 {
	result := 0.0
	for _, f := range features {
		w, ok := m.Weight[f]
		if ok {
			result = result + w*1.0
		}
	}
	return result
}

func (m MIRAClassifier) Predict(features feature.FeatureVector) model.LabelType {
	if m.PredictScore(features) > 0 {
		return model.POSITIVE
	}
	return model.NEGATIVE
}

func (m MIRAClassifier) SortByScore(examples model.Examples) model.Examples {
	var unlabeledExamples model.Examples
	for _, e := range util.FilterUnlabeledExamples(examples) {
		e.Score = m.PredictScore(e.Fv)
		if !e.IsLabeled() && e.Score != 0.0 {
			unlabeledExamples = append(unlabeledExamples, e)
		}
	}

	sort.Sort(unlabeledExamples)
	return unlabeledExamples
}

func (m MIRAClassifier) GetWeight(f string) float64 {
	w, ok := m.Weight[f]
	if ok {
		return w
	}
	return 0.0
}

func (m MIRAClassifier) GetActiveFeatures() []string {
	result := make([]string, 0)
	for f := range m.Weight {
		result = append(result, f)
	}
	return result
}
