package main

type BinaryClassifier interface {
	Learn(Example)
	PredictScore(FeatureVector) float64
	Predict(FeatureVector) LabelType
	SortByScore(Examples) Examples
	GetWeight(string) float64
}
