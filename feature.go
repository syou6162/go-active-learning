package main

import (
	"net/url"
	"strings"
	"unicode"

	"github.com/ikawaha/kagome/tokenizer"
	"github.com/jdkato/prose/tag"
	"github.com/jdkato/prose/tokenize"
)

type FeatureVector []string

func isJapanese(str string) bool {
	flag := false
	for _, r := range str {
		if unicode.In(r, unicode.Hiragana) || unicode.In(r, unicode.Katakana) {
			flag = true
		}
	}

	if strings.ContainsAny(str, "。、") {
		flag = true
	}

	return flag
}

func extractEngNounFeatures(s string, prefix string) FeatureVector {
	var fv FeatureVector
	if s == "" {
		return fv
	}

	words := tokenize.NewTreebankWordTokenizer().Tokenize(s)
	tagger := tag.NewPerceptronTagger()
	for _, tok := range tagger.Tag(words) {
		switch tok.Tag {
		// https://www.ling.upenn.edu/courses/Fall_2003/ling001/penn_treebank_pos.html
		case "NN", "NNS", "NNP", "NNPS", "PRP", "PRP$":
			fv = append(fv, prefix+":"+strings.ToLower(tok.Text))
		}
	}

	return fv
}

func extractJpnNounFeatures(s string, prefix string) FeatureVector {
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

func ExtractNounFeatures(s string, prefix string) FeatureVector {
	if isJapanese(s) {
		return extractJpnNounFeatures(s, prefix)
	} else {
		return extractEngNounFeatures(s, prefix)
	}
}

func ExtractHostFeature(urlString string) string {
	prefix := "HOST"
	u, err := url.Parse(urlString)
	if err != nil {
		return prefix + ":INVALID_HOST"
	}
	return prefix + ":" + u.Host
}

func ExtractFeatures(e Example) FeatureVector {
	var fv FeatureVector
	fv = append(fv, "BIAS")
	fv = append(fv, ExtractHostFeature(e.FinalUrl))
	fv = append(fv, ExtractNounFeatures(e.Title, "TITLE")...)
	fv = append(fv, ExtractNounFeatures(e.Description, "DESCRIPTION")...)
	fv = append(fv, ExtractNounFeatures(e.Body, "BODY")...)
	return fv
}
