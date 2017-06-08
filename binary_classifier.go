package main

import "sort"

type BinaryClassifier interface {
	PredictScore(FeatureVector) float64
	Predict(FeatureVector) LabelType
	SortByScore(Examples) Examples
	GetWeight(string) float64
}

func SortByScore(model BinaryClassifier, examples Examples) Examples {
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
