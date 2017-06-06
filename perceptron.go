package main

import (
	"fmt"
	"os"
	"sort"
)

type Model struct {
	weight    map[string]float64
	cumWeight map[string]float64
	count     int
}

func NewModel() *Model {
	return &Model{make(map[string]float64), make(map[string]float64), 1}
}

func GetAccuracy(gold []LabelType, predict []LabelType) float64 {
	if len(gold) != len(predict) {
		return 0.0
	}
	sum := 0.0
	for i, v := range gold {
		if v == predict[i] {
			sum += 1.0
		}
	}
	return sum / float64(len(gold))
}

func GetPrecision(gold []LabelType, predict []LabelType) float64 {
	tp := 0.0
	fp := 0.0
	for i, v := range gold {
		if v == predict[i] {
			tp += 1.0
		}
		if predict[i] == POSITIVE && v == NEGATIVE {
			fp += 1.0
		}
	}
	return tp / (tp + fp)
}

func GetRecall(gold []LabelType, predict []LabelType) float64 {
	tp := 0.0
	fn := 0.0
	for i, v := range gold {
		if v == predict[i] {
			tp += 1.0
		}
		if predict[i] == NEGATIVE && v == POSITIVE {
			fn += 1.0
		}
	}
	return tp / (tp + fn)
}

func (model *Model) Learn(example Example) {
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

func (model *Model) predictForTraining(features FeatureVector) LabelType {
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

func (model Model) PredictScore(features FeatureVector) float64 {
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

func (model Model) Predict(features FeatureVector) LabelType {
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

func (model Model) SortByScore(examples Examples) Examples {
	var unlabeledExamples Examples
	for _, e := range FilterUnlabeledExamples(examples) {
		e.Score = model.PredictScore(e.Fv)
		if !e.IsLabeled() && e.Score != 0.0 {
			unlabeledExamples = append(unlabeledExamples, e)
		}
	}

	sort.Sort(unlabeledExamples)
	return unlabeledExamples
}

func TrainedModel(examples Examples) *Model {
	train := FilterLabeledExamples(examples)
	model := NewModel()
	for iter := 0; iter < 30; iter++ {
		shuffle(train)
		for _, example := range train {
			model.Learn(*example)
		}

		trainPredicts := make([]LabelType, len(train))
		for i, example := range train {
			trainPredicts[i] = model.Predict(example.Fv)
		}
		accuracy := GetAccuracy(ExtractGoldLabels(train), trainPredicts)
		precision := GetPrecision(ExtractGoldLabels(train), trainPredicts)
		recall := GetRecall(ExtractGoldLabels(train), trainPredicts)
		f := (2 * recall * precision) / (recall + precision)
		fmt.Fprintln(os.Stderr, fmt.Sprintf("Iter:%d\tAccuracy:%0.03f\tPrecision:%0.03f\tRecall:%0.03f\tF-value:%0.03f", iter, accuracy, precision, recall, f))
	}
	return model
}

func (model Model) GetAveragedWeight(f string) float64 {
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
