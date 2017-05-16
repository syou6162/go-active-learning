package main

import (
	"github.com/ikawaha/kagome/tokenizer"
	"strings"

	"unicode/utf8"
)

type FeatureVector []string

func ExtractNounFeatures(s string, prefix string) FeatureVector {
	var fv FeatureVector
	if s == "" {
		return fv
	}
	t := tokenizer.New()
	tokens := t.Tokenize(s)
	for _, token := range tokens {
		if token.Pos() == "名詞" {
			fv = append(fv, prefix+":"+token.Surface)
		}
	}
	return fv
}

func ExtractFeatures(e Example) FeatureVector {
	var fv FeatureVector

	html := strings.ToLower(strings.Replace(e.RawHTML, " ", "", -1))
	if !utf8.ValidString(html) {
		return fv
	}

	if !utf8.ValidString(e.Title) {
		return fv
	}
	fv = append(fv, ExtractNounFeatures(e.Title, "TITLE")...)

	if !utf8.ValidString(e.Description) {
		return fv
	}
	fv = append(fv, ExtractNounFeatures(e.Description, "DESCRIPTION")...)

	if !utf8.ValidString(e.Body) {
		return fv
	}
	fv = append(fv, ExtractNounFeatures(e.Body, "BODY")...)

	return fv
}
