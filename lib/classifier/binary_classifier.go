package classifier

import (
	"sort"

	"github.com/syou6162/go-active-learning/lib/feature"
	"github.com/syou6162/go-active-learning/lib/model"
	"github.com/syou6162/go-active-learning/lib/util"
)

type BinaryClassifier interface {
	PredictScore(feature.FeatureVector) float64
	Predict(feature.FeatureVector) model.LabelType
	SortByScore(model.Examples) model.Examples
	GetWeight(string) float64
	GetActiveFeatures() []string
}

func NewBinaryClassifier(examples model.Examples) BinaryClassifier {
	// return NewPerceptronClassifier(examples)
	return NewMIRAClassifierByCrossValidation(examples)
}

func SortByScore(m BinaryClassifier, examples model.Examples) model.Examples {
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
