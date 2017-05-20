package main

import (
	"strings"

	"github.com/ikawaha/kagome/tokenizer"
)

type FeatureVector []string

func ExtractNounFeatures(s string, prefix string) FeatureVector {
	var fv FeatureVector
	if s == "" {
		return fv
	}
	t := tokenizer.New()
	tokens := t.Tokenize(strings.ToLower(s))
	for _, token := range tokens {
		if token.Pos() == "名詞" {
			surface := token.Surface
			if len(token.Features()) >= 2 && token.Features()[1] == "数" {
				surface = "NUM"
			}
			fv = append(fv, prefix+":"+surface)
		}
	}
	return fv
}

func ExtractFeatures(e Example) FeatureVector {
	var fv FeatureVector
	fv = append(fv, "BIAS")
	fv = append(fv, ExtractNounFeatures(e.Title, "TITLE")...)
	fv = append(fv, ExtractNounFeatures(e.Description, "DESCRIPTION")...)
	fv = append(fv, ExtractNounFeatures(e.Body, "BODY")...)
	return fv
}
