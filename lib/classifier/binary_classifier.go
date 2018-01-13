package classifier

import (
	"sort"

	"github.com/syou6162/go-active-learning/lib/example"
	"github.com/syou6162/go-active-learning/lib/feature"
	"github.com/syou6162/go-active-learning/lib/util"
)

type BinaryClassifier interface {
	PredictScore(feature.FeatureVector) float64
	Predict(feature.FeatureVector) example.LabelType
	SortByScore(example.Examples) example.Examples
	GetWeight(string) float64
	GetActiveFeatures() []string
}

func NewBinaryClassifier(examples example.Examples) BinaryClassifier {
	// return NewPerceptronClassifier(examples)
	return NewMIRAClassifierByCrossValidation(examples)
}

func SortByScore(model BinaryClassifier, examples example.Examples) example.Examples {
	var unlabeledExamples example.Examples
	for _, e := range util.FilterUnlabeledExamples(examples) {
		e.Score = model.PredictScore(e.Fv)
		if !e.IsLabeled() && e.Score != 0.0 {
			unlabeledExamples = append(unlabeledExamples, e)
		}
	}

	sort.Sort(unlabeledExamples)
	return unlabeledExamples
}
