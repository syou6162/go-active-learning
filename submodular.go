package main

import (
	"math"
)

func SelectSubExamplesBySubModular(model *Model, whole Examples, sizeConstraint int, alpha float64, r float64) Examples {
	selected := Examples{}
	remainings := whole
	simMat := GetSimilarityMatrix(model, whole)
	for {
		if len(selected) >= sizeConstraint || len(remainings) == 0 {
			break
		}
		argmax := SelectBestExample(simMat, remainings, selected, whole, alpha, r)
		selected = append(selected, remainings[argmax])
		remainings = append(remainings[:argmax], remainings[argmax+1:]...)
	}
	// (1 - 1/e)/2の保証を与えるためにはもうちょっと頑張る必要があるが、省略している
	// http://www.anthology.aclweb.org/E/E09/E09-1089.pdf
	return selected
}

func SelectBestExample(mat SimilarityMatrix, remainings Examples, selected Examples, whole Examples, alpha float64, r float64) int {
	maxScore := math.Inf(-1)
	argmax := 0
	for idx, example := range remainings {
		subset := Examples{}
		for _, e := range selected {
			subset = append(subset, e)
		}
		subset = append(subset, example)
		c1 := CoverageFunction(mat, subset, whole, alpha)
		c2 := CoverageFunction(mat, selected, whole, alpha)
		score := (c1 - c2) / math.Pow(float64(len(example.Fv)), r)
		if score >= maxScore {
			argmax = idx
			maxScore = score
		}
	}
	return argmax
}

func CoverageFunction(mat SimilarityMatrix, subset Examples, whole Examples, alpha float64) float64 {
	sum := 0.0
	for _, e := range whole {
		sum += math.Min(
			coverageFunction(mat, e, subset),
			alpha*coverageFunction(mat, e, whole),
		)
	}
	return sum
}

func coverageFunction(mat SimilarityMatrix, example *Example, examples Examples) float64 {
	sum := 0.0
	for _, e := range examples {
		sum += GetCosineSimilarity(mat, e, example)
	}
	return sum
}

type SimilarityMatrix map[string]float64

func GetSimilarityMatrix(model *Model, examples Examples) SimilarityMatrix {
	mat := SimilarityMatrix{}
	for _, e1 := range examples {
		for _, e2 := range examples {
			mat[e1.Url+"+"+e2.Url] = cosineSimilarity(model, e1, e2)
		}
	}
	return mat
}

func GetCosineSimilarity(mat SimilarityMatrix, e1 *Example, e2 *Example) float64 {
	return mat[e1.Url+"+"+e2.Url]
}

func cosineSimilarity(model *Model, e1 *Example, e2 *Example) float64 {
	sum := 0.0

	// Find features that exist in both e1 and e2
	existBoth := make(map[string]bool)
	for _, f := range e1.Fv {
		existBoth[f] = true
	}
	for _, f := range e2.Fv {
		if _, ok := existBoth[f]; !ok {
			delete(existBoth, f)
		}
	}

	for k := range existBoth {
		w := model.GetAveragedWeight(k)
		sum += w * w
	}
	return sum / (Norm(model, e1) * Norm(model, e2))
}

func Norm(model *Model, e *Example) float64 {
	sum := 0.0
	for _, f := range e.Fv {
		w := model.GetAveragedWeight(f)
		sum += w * w
	}
	return math.Sqrt(sum)
}
