package main

import (
	"github.com/ikawaha/kagome/tokenizer"
)

type FeatureVector []string

func ExtractFeatures(title string) FeatureVector {
	var fv FeatureVector
	t := tokenizer.New()
	tokens := t.Tokenize(title)
	for _, token := range tokens {
		if token.Pos() == "名詞" {
			fv = append(fv, token.Surface)
		}
	}
	return fv
}
