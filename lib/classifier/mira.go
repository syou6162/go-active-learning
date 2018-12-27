package classifier

import (
	"fmt"
	"math"
	"os"
	"runtime"
	"sort"
	"sync"
	"time"

	mkr "github.com/mackerelio/mackerel-client-go"
	"github.com/pkg/errors"
	"github.com/syou6162/go-active-learning/lib/evaluation"
	"github.com/syou6162/go-active-learning/lib/feature"
	"github.com/syou6162/go-active-learning/lib/model"
	"github.com/syou6162/go-active-learning/lib/util"
)

type MIRAClassifier struct {
	Weight map[string]float64 `json:"Weight"`
	C      float64            `json:"C"`
}

func newMIRAClassifier(c float64) *MIRAClassifier {
	return &MIRAClassifier{make(map[string]float64), c}
}

func NewMIRAClassifier(examples model.Examples, c float64) *MIRAClassifier {
	train := util.FilterLabeledExamples(examples)
	model := newMIRAClassifier(c)
	for iter := 0; iter < 30; iter++ {
		util.Shuffle(train)
		for _, example := range train {
			model.learn(*example)
		}
	}
	return model
}

func OverSamplingPositiveExamples(examples model.Examples) model.Examples {
	overSampled := model.Examples{}
	posExamples := model.Examples{}
	negExamples := model.Examples{}

	numNeg := 0

	for _, e := range examples {
		if e.Label == model.NEGATIVE {
			numNeg += 1
			negExamples = append(negExamples, e)
		} else if e.Label == model.POSITIVE {
			posExamples = append(posExamples, e)
		}
	}

	for len(overSampled) <= numNeg {
		util.Shuffle(posExamples)
		overSampled = append(overSampled, posExamples[0])
	}
	overSampled = append(overSampled, negExamples...)
	util.Shuffle(overSampled)

	return overSampled
}

func ExtractGoldLabels(examples model.Examples) []model.LabelType {
	golds := make([]model.LabelType, 0, 0)
	for _, e := range examples {
		golds = append(golds, e.Label)
	}
	return golds
}

type MIRAResult struct {
	mira   MIRAClassifier
	FValue float64
}

type MIRAResultList []MIRAResult

func (l MIRAResultList) Len() int           { return len(l) }
func (l MIRAResultList) Less(i, j int) bool { return l[i].FValue < l[j].FValue }
func (l MIRAResultList) Swap(i, j int)      { l[i], l[j] = l[j], l[i] }

func NewMIRAClassifierByCrossValidation(examples model.Examples) (*MIRAClassifier, error) {
	util.Shuffle(examples)
	train, dev := util.SplitTrainAndDev(util.FilterLabeledExamples(examples))
	train = OverSamplingPositiveExamples(train)

	params := []float64{1000, 500, 100, 50, 10.0, 5.0, 1.0, 0.5, 0.1, 0.05, 0.01, 0.005, 0.001}
	miraResults := MIRAResultList{}

	wg := &sync.WaitGroup{}
	cpus := runtime.NumCPU()
	runtime.GOMAXPROCS(cpus)

	models := make([]*MIRAClassifier, len(params))
	for idx, c := range params {
		wg.Add(1)
		go func(idx int, c float64) {
			defer wg.Done()
			model := NewMIRAClassifier(train, c)
			models[idx] = model
		}(idx, c)
	}
	wg.Wait()

	maxAccuracy := 0.0
	maxPrecision := 0.0
	maxRecall := 0.0
	maxFvalue := math.Inf(-1)
	for _, m := range models {
		c := m.C
		devPredicts := make([]model.LabelType, len(dev))
		for i, example := range dev {
			devPredicts[i] = m.Predict(example.Fv)
		}
		accuracy := evaluation.GetAccuracy(ExtractGoldLabels(dev), devPredicts)
		precision := evaluation.GetPrecision(ExtractGoldLabels(dev), devPredicts)
		recall := evaluation.GetRecall(ExtractGoldLabels(dev), devPredicts)
		f := (2 * recall * precision) / (recall + precision)
		if math.IsNaN(f) {
			continue
		}
		miraResults = append(miraResults, MIRAResult{*m, f})
		fmt.Fprintln(os.Stderr, fmt.Sprintf("C:%0.03f\tAccuracy:%0.03f\tPrecision:%0.03f\tRecall:%0.03f\tF-value:%0.03f", c, accuracy, precision, recall, f))
		if f >= maxFvalue {
			maxAccuracy = accuracy
			maxPrecision = precision
			maxRecall = recall
			maxFvalue = f
		}
	}
	if len(miraResults) == 0 {
		return nil, errors.New("Failed to learn MIRA")
	}
	err := postEvaluatedMetricsToMackerel(maxAccuracy, maxPrecision, maxRecall, maxFvalue)
	if err != nil {
		fmt.Println(err.Error())
	}

	sort.Sort(sort.Reverse(miraResults))
	bestModel := &miraResults[0].mira
	examples = OverSamplingPositiveExamples(examples)
	util.Shuffle(examples)
	return NewMIRAClassifier(util.FilterLabeledExamples(examples), bestModel.C), nil
}

func postEvaluatedMetricsToMackerel(accuracy float64, precision float64, recall float64, fvalue float64) error {
	apiKey := os.Getenv("MACKEREL_API_KEY")
	serviceName := os.Getenv("MACKEREL_SERVICE_NAME")
	if apiKey == "" || serviceName == "" {
		return nil
	}

	client := mkr.NewClient(apiKey)
	now := time.Now().Unix()
	err := client.PostServiceMetricValues(serviceName, []*mkr.MetricValue{
		{
			Name:  "evaluation.accuracy",
			Time:  now,
			Value: accuracy,
		},
		{
			Name:  "evaluation.precision",
			Time:  now,
			Value: precision,
		},
		{
			Name:  "evaluation.recall",
			Time:  now,
			Value: recall,
		},
		{
			Name:  "evaluation.fvalue",
			Time:  now,
			Value: fvalue,
		},
	})
	return err
}

func (m *MIRAClassifier) learn(example model.Example) {
	tmp := float64(example.Label) * m.PredictScore(example.Fv) // y w^T x
	loss := 0.0
	if tmp < 1.0 {
		loss = 1 - tmp
	}

	norm := float64(len(example.Fv) * len(example.Fv))
	// tau := math.Min(m.C, loss/norm) // update by PA-I
	tau := loss / (norm + 1.0/m.C) // update by PA-II

	if tau != 0.0 {
		for _, f := range example.Fv {
			w, _ := m.Weight[f]
			m.Weight[f] = w + tau*float64(example.Label)
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
