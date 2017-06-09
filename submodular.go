package main

import (
	"math"
)

func SelectSubExamplesBySubModular(whole Examples, sizeConstraint int, alpha float64, r float64) Examples {
	selected := Examples{}
	remainings := whole
	simMat := GetSimilarityMatrixByTFIDF(whole)
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

func GetSimilarityMatrixByTFIDF(examples Examples) SimilarityMatrix {
	idf := GetIDF(examples)

	dfByURL := make(map[string]map[string]float64)
	sumByUrl := make(map[string]float64)
	for _, e := range examples {
		df := GetDF(*e)
		dfByURL[e.Url] = df

		sum := 0.0
		for k, v := range df {
			sum += v * v * idf[k] * idf[k]
		}
		sumByUrl[e.Url] = sum
	}

	mat := SimilarityMatrix{}
	for _, e1 := range examples {
		df1 := dfByURL[e1.Url]
		s1 := math.Sqrt(sumByUrl[e1.Url])

		for _, e2 := range examples {
			df2 := dfByURL[e2.Url]
			s2 := math.Sqrt(sumByUrl[e1.Url])

			s := 0.0
			for k, v := range df2 {
				s += v * df1[k] * idf[k] * idf[k]
			}
			mat[e1.Url+"+"+e2.Url] = s / (s1 * s2)
		}
	}
	return mat
}

func GetCosineSimilarity(mat SimilarityMatrix, e1 *Example, e2 *Example) float64 {
	return mat[e1.Url+"+"+e2.Url]
}

func GetDF(example Example) map[string]float64 {
	df := make(map[string]float64)
	n := 0.0
	fv := ExtractNounFeatures(example.Body, "BODY")

	for _, f := range fv {
		df[f]++
		n++
	}

	for k, v := range df {
		df[k] = v/ n
	}
	return df
}

func GetIDF(examples Examples) map[string]float64 {
	idf := make(map[string]float64)
	cnt := make(map[string]float64)
	n := float64(len(examples))

	for _, e := range examples {
		fv := ExtractNounFeatures(e.Body, "BODY")
		for _, f := range fv {
			cnt[f]++
		}
	}

	for k, v := range cnt {
		idf[k] = math.Log(n / v) + 1
	}
	return idf
}
