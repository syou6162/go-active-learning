package main

import (
	"fmt"
	"os"
)

type PerceptronClassifier struct {
	weight    map[string]float64
	cumWeight map[string]float64
	count     int
}

func newPerceptronClassifier() *PerceptronClassifier {
	return &PerceptronClassifier{make(map[string]float64), make(map[string]float64), 1}
}

func NewPerceptronClassifier(examples Examples) *PerceptronClassifier {
	train, dev := splitTrainAndDev(FilterLabeledExamples(examples))
	model := newPerceptronClassifier()
	for iter := 0; iter < 30; iter++ {
		shuffle(train)
		for _, example := range train {
			model.learn(*example)
		}

		devPredicts := make([]LabelType, len(dev))
		for i, example := range dev {
			devPredicts[i] = model.Predict(example.Fv)
		}
		accuracy := GetAccuracy(ExtractGoldLabels(dev), devPredicts)
		precision := GetPrecision(ExtractGoldLabels(dev), devPredicts)
		recall := GetRecall(ExtractGoldLabels(dev), devPredicts)
		f := (2 * recall * precision) / (recall + precision)
		fmt.Fprintln(os.Stderr, fmt.Sprintf("Iter:%d\tAccuracy:%0.03f\tPrecision:%0.03f\tRecall:%0.03f\tF-value:%0.03f", iter, accuracy, precision, recall, f))
	}
	return model
}

func (model *PerceptronClassifier) learn(example Example) {
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

func (model *PerceptronClassifier) predictForTraining(features FeatureVector) LabelType {
	result := 0.0
	for _, f := range features {
		w, ok := model.weight[f]
		if ok {
			result = result + w*1.0
		}
	}
	if result > 0 {
		return POSITIVE
	}
	return NEGATIVE
}

func (model PerceptronClassifier) PredictScore(features FeatureVector) float64 {
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

func (model PerceptronClassifier) Predict(features FeatureVector) LabelType {
	if model.PredictScore(features) > 0 {
		return POSITIVE
	}
	return NEGATIVE
}

func ExtractGoldLabels(examples Examples) []LabelType {
	golds := make([]LabelType, 0, 0)
	for _, e := range examples {
		golds = append(golds, e.Label)
	}
	return golds
}

func (model PerceptronClassifier) SortByScore(examples Examples) Examples {
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
