package classifier

import (
	"fmt"
	"os"

	"github.com/syou6162/go-active-learning/lib/evaluation"
	"github.com/syou6162/go-active-learning/lib/feature"
	"github.com/syou6162/go-active-learning/lib/model"
	"github.com/syou6162/go-active-learning/lib/util"
)

type PerceptronClassifier struct {
	weight    map[string]float64
	cumWeight map[string]float64
	count     int
}

func newPerceptronClassifier() *PerceptronClassifier {
	return &PerceptronClassifier{make(map[string]float64), make(map[string]float64), 1}
}

func NewPerceptronClassifier(examples model.Examples) *PerceptronClassifier {
	train, dev := util.SplitTrainAndDev(util.FilterLabeledExamples(examples))
	m := newPerceptronClassifier()
	for iter := 0; iter < 30; iter++ {
		util.Shuffle(train)
		for _, example := range train {
			m.learn(*example)
		}

		devPredicts := make([]model.LabelType, len(dev))
		for i, example := range dev {
			devPredicts[i] = m.Predict(example.Fv)
		}
		accuracy := evaluation.GetAccuracy(ExtractGoldLabels(dev), devPredicts)
		precision := evaluation.GetPrecision(ExtractGoldLabels(dev), devPredicts)
		recall := evaluation.GetRecall(ExtractGoldLabels(dev), devPredicts)
		f := (2 * recall * precision) / (recall + precision)
		fmt.Fprintln(os.Stderr, fmt.Sprintf("Iter:%d\tAccuracy:%0.03f\tPrecision:%0.03f\tRecall:%0.03f\tF-value:%0.03f", iter, accuracy, precision, recall, f))
	}
	return m
}

func (m *PerceptronClassifier) learn(example model.Example) {
	predict := m.predictForTraining(example.Fv)
	if example.Label != predict {
		for _, f := range example.Fv {
			w, _ := m.weight[f]
			cumW, _ := m.cumWeight[f]
			m.weight[f] = w + float64(example.Label)*1.0
			m.cumWeight[f] = cumW + float64(m.count)*float64(example.Label)*1.0
		}
		m.count += 1
	}
}

func (m *PerceptronClassifier) predictForTraining(features feature.FeatureVector) model.LabelType {
	result := 0.0
	for _, f := range features {
		w, ok := m.weight[f]
		if ok {
			result = result + w*1.0
		}
	}
	if result > 0 {
		return model.POSITIVE
	}
	return model.NEGATIVE
}

func (m PerceptronClassifier) PredictScore(features feature.FeatureVector) float64 {
	result := 0.0
	for _, f := range features {
		w, ok := m.weight[f]
		if ok {
			result = result + w*1.0
		}

		w, ok = m.cumWeight[f]
		if ok {
			result = result - w*1.0/float64(m.count)
		}

	}
	return result
}

func (m PerceptronClassifier) Predict(features feature.FeatureVector) model.LabelType {
	if m.PredictScore(features) > 0 {
		return model.POSITIVE
	}
	return model.NEGATIVE
}

func ExtractGoldLabels(examples model.Examples) []model.LabelType {
	golds := make([]model.LabelType, 0, 0)
	for _, e := range examples {
		golds = append(golds, e.Label)
	}
	return golds
}

func (m PerceptronClassifier) SortByScore(examples model.Examples) model.Examples {
	return SortByScore(m, examples)
}

func (m PerceptronClassifier) GetWeight(f string) float64 {
	result := 0.0
	w, ok := m.weight[f]
	if ok {
		result = result + w*1.0
	}

	w, ok = m.cumWeight[f]
	if ok {
		result = result - w*1.0/float64(m.count)
	}
	return result
}

func (m PerceptronClassifier) GetActiveFeatures() []string {
	result := make([]string, 0)
	for f := range m.weight {
		result = append(result, f)
	}
	return result
}
