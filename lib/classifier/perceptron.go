package classifier

import (
	"fmt"
	"os"

	"github.com/syou6162/go-active-learning/lib/example"
	"github.com/syou6162/go-active-learning/lib/feature"
	"github.com/syou6162/go-active-learning/lib/util"
	"github.com/syou6162/go-active-learning/lib/evaluation"
)

type PerceptronClassifier struct {
	weight    map[string]float64
	cumWeight map[string]float64
	count     int
}

func newPerceptronClassifier() *PerceptronClassifier {
	return &PerceptronClassifier{make(map[string]float64), make(map[string]float64), 1}
}

func NewPerceptronClassifier(examples example.Examples) *PerceptronClassifier {
	train, dev := util.SplitTrainAndDev(util.FilterLabeledExamples(examples))
	model := newPerceptronClassifier()
	for iter := 0; iter < 30; iter++ {
		util.Shuffle(train)
		for _, example := range train {
			model.learn(*example)
		}

		devPredicts := make([]example.LabelType, len(dev))
		for i, example := range dev {
			devPredicts[i] = model.Predict(example.Fv)
		}
		accuracy := evaluation.GetAccuracy(ExtractGoldLabels(dev), devPredicts)
		precision := evaluation.GetPrecision(ExtractGoldLabels(dev), devPredicts)
		recall := evaluation.GetRecall(ExtractGoldLabels(dev), devPredicts)
		f := (2 * recall * precision) / (recall + precision)
		fmt.Fprintln(os.Stderr, fmt.Sprintf("Iter:%d\tAccuracy:%0.03f\tPrecision:%0.03f\tRecall:%0.03f\tF-value:%0.03f", iter, accuracy, precision, recall, f))
	}
	return model
}

func (model *PerceptronClassifier) learn(example example.Example) {
	predict := model.predictForTraining(example.Fv)
	if example.Label != predict {
		for _, f := range example.Fv {
			w, _ := model.weight[f]
			cumW, _ := model.cumWeight[f]
			model.weight[f] = w + float64(example.Label)*1.0
			model.cumWeight[f] = cumW + float64(model.count)*float64(example.Label)*1.0
		}
		model.count += 1
	}
}

func (model *PerceptronClassifier) predictForTraining(features feature.FeatureVector) example.LabelType {
	result := 0.0
	for _, f := range features {
		w, ok := model.weight[f]
		if ok {
			result = result + w*1.0
		}
	}
	if result > 0 {
		return example.POSITIVE
	}
	return example.NEGATIVE
}

func (model PerceptronClassifier) PredictScore(features feature.FeatureVector) float64 {
	result := 0.0
	for _, f := range features {
		w, ok := model.weight[f]
		if ok {
			result = result + w*1.0
		}

		w, ok = model.cumWeight[f]
		if ok {
			result = result - w*1.0/float64(model.count)
		}

	}
	return result
}

func (model PerceptronClassifier) Predict(features feature.FeatureVector) example.LabelType {
	if model.PredictScore(features) > 0 {
		return example.POSITIVE
	}
	return example.NEGATIVE
}

func ExtractGoldLabels(examples example.Examples) []example.LabelType {
	golds := make([]example.LabelType, 0, 0)
	for _, e := range examples {
		golds = append(golds, e.Label)
	}
	return golds
}

func (model PerceptronClassifier) SortByScore(examples example.Examples) example.Examples {
	return SortByScore(model, examples)
}

func (model PerceptronClassifier) GetWeight(f string) float64 {
	result := 0.0
	w, ok := model.weight[f]
	if ok {
		result = result + w*1.0
	}

	w, ok = model.cumWeight[f]
	if ok {
		result = result - w*1.0/float64(model.count)
	}
	return result
}

func (model PerceptronClassifier) GetActiveFeatures() []string {
	result := make([]string, 0)
	for f := range model.weight {
		result = append(result, f)
	}
	return result
}
