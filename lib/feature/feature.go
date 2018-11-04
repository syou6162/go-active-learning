package feature

import (
	"encoding/json"
	"net/url"
	"strings"
	"sync"
	"unicode"

	"github.com/ikawaha/kagome/tokenizer"
	"github.com/jdkato/prose/tag"
	"github.com/jdkato/prose/tokenize"
)

type FeatureVector []string

func (fv *FeatureVector) MarshalBinary() ([]byte, error) {
	json, err := json.Marshal(fv)
	if err != nil {
		return nil, err
	}
	return []byte(json), nil
}

func (fv *FeatureVector) UnmarshalBinary(data []byte) error {
	err := json.Unmarshal(data, fv)
	if err != nil {
		return err
	}
	return nil
}

var (
	japaneseTokenizer     *tokenizer.Tokenizer
	japaneseTokenizerOnce sync.Once
	englishTokenizer      *tokenize.TreebankWordTokenizer
	englishTokenizerOnce  sync.Once
	englishTagger         *tag.PerceptronTagger
	englishTaggerOnce     sync.Once
)

func GetJapaneseTokenizer() *tokenizer.Tokenizer {
	japaneseTokenizerOnce.Do(func() {
		t := tokenizer.New()
		japaneseTokenizer = &t
	})

	return japaneseTokenizer
}

func GetEnglishTokenizer() *tokenize.TreebankWordTokenizer {
	englishTokenizerOnce.Do(func() {
		englishTokenizer = tokenize.NewTreebankWordTokenizer()
	})
	return englishTokenizer
}

func GetEnglishTagger() *tag.PerceptronTagger {
	englishTaggerOnce.Do(func() {
		englishTagger = tag.NewPerceptronTagger()
	})
	return englishTagger
}

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

func extractEngNounFeaturesWithoutPrefix(s string) FeatureVector {
	var fv FeatureVector
	if s == "" {
		return fv
	}

	words := GetEnglishTokenizer().Tokenize(s)
	tagger := GetEnglishTagger()
	for _, tok := range tagger.Tag(words) {
		switch tok.Tag {
		// https://www.ling.upenn.edu/courses/Fall_2003/ling001/penn_treebank_pos.html
		case "NN", "NNS", "NNP", "NNPS", "PRP", "PRP$":
			fv = append(fv, strings.ToLower(tok.Text))
		}
	}

	return fv
}

func extractEngNounFeatures(s string, prefix string) FeatureVector {
	var fv FeatureVector
	for _, surface := range extractEngNounFeaturesWithoutPrefix(s) {
		fv = append(fv, prefix+":"+surface)
	}
	return fv
}

func ExtractJpnNounFeaturesWithoutPrefix(s string) FeatureVector {
	var fv FeatureVector
	if s == "" {
		return fv
	}
	t := GetJapaneseTokenizer()
	tokens := t.Tokenize(strings.ToLower(s))
	for _, token := range tokens {
		if token.Pos() == "名詞" {
			surface := token.Surface
			if len(token.Features()) >= 2 && token.Features()[1] == "数" {
				surface = "NUM"
			}
			fv = append(fv, surface)
		}
	}
	return fv
}

func ExtractJpnNounFeatures(s string, prefix string) FeatureVector {
	var fv FeatureVector
	for _, surface := range ExtractJpnNounFeaturesWithoutPrefix(s) {
		fv = append(fv, prefix+":"+surface)
	}
	return fv
}

func ExtractNounFeatures(s string, prefix string) FeatureVector {
	if isJapanese(s) {
		return ExtractJpnNounFeatures(s, prefix)
	} else {
		return extractEngNounFeatures(s, prefix)
	}
}

func ExtractNounFeaturesWithoutPrefix(s string) FeatureVector {
	if isJapanese(s) {
		return ExtractJpnNounFeaturesWithoutPrefix(s)
	} else {
		return extractEngNounFeaturesWithoutPrefix(s)
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

func ExtractPath(urlString string) string {
	path := ""
	u, err := url.Parse(urlString)
	if err != nil {
		return path
	}
	return u.Path
}
