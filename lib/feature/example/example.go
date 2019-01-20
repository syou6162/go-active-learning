package example_feature

import (
	"net/url"
	"strings"
	"sync"
	"unicode"

	"github.com/ikawaha/kagome/tokenizer"
	"github.com/jdkato/prose/tag"
	"github.com/jdkato/prose/tokenize"
	"github.com/syou6162/go-active-learning/lib/feature"
)

var excludingWordList = []string{
	`:`, `;`,
	`,`, `.`,
	`"`, `''`,
	`+`, `-`, `*`, `/`, `|`, `++`, `--`,
	`[`, `]`,
	`{`, `}`,
	`(`, `)`,
	`<`, `>`,
	`「`, `」`,
	`／`,
	`@`, `#`, `~`, `%`, `$`, `^`,
}

var (
	japaneseTokenizer     *tokenizer.Tokenizer
	japaneseTokenizerOnce sync.Once
	englishTokenizer      *tokenize.TreebankWordTokenizer
	englishTokenizerOnce  sync.Once
	englishTagger         *tag.PerceptronTagger
	englishTaggerOnce     sync.Once
	excludingWordMapOnce  sync.Once
)

var excludingWordMap = make(map[string]bool)

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
	for _, r := range str {
		if unicode.In(r, unicode.Hiragana) || unicode.In(r, unicode.Katakana) || unicode.In(r, unicode.Han) {
			return true
		}
	}

	if strings.ContainsAny(str, "。、") {
		return true
	}

	return false
}

func IsExcludingWord(w string) bool {
	excludingWordMapOnce.Do(func() {
		for _, w := range excludingWordList {
			excludingWordMap[w] = true
		}
	})
	if _, ok := excludingWordMap[w]; ok {
		return true
	}
	return false
}

func extractEngNounFeaturesWithoutPrefix(s string) feature.FeatureVector {
	var fv feature.FeatureVector
	if s == "" {
		return fv
	}

	words := GetEnglishTokenizer().Tokenize(s)
	tagger := GetEnglishTagger()
	for _, tok := range tagger.Tag(words) {
		if IsExcludingWord(tok.Text) {
			continue
		}
		switch tok.Tag {
		// https://www.ling.upenn.edu/courses/Fall_2003/ling001/penn_treebank_pos.html
		case "NN", "NNS", "NNP", "NNPS", "PRP", "PRP$":
			fv = append(fv, strings.ToLower(tok.Text))
		}
	}

	return fv
}

func extractEngNounFeatures(s string, prefix string) feature.FeatureVector {
	var fv feature.FeatureVector
	for _, surface := range extractEngNounFeaturesWithoutPrefix(s) {
		fv = append(fv, prefix+":"+surface)
	}
	return fv
}

func ExtractJpnNounFeaturesWithoutPrefix(s string) feature.FeatureVector {
	var fv feature.FeatureVector
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
			if IsExcludingWord(surface) {
				continue
			}
			fv = append(fv, surface)
		}
	}
	return fv
}

func ExtractJpnNounFeatures(s string, prefix string) feature.FeatureVector {
	var fv feature.FeatureVector
	for _, surface := range ExtractJpnNounFeaturesWithoutPrefix(s) {
		fv = append(fv, prefix+":"+surface)
	}
	return fv
}

func ExtractNounFeatures(s string, prefix string) feature.FeatureVector {
	if isJapanese(s) {
		return ExtractJpnNounFeatures(s, prefix)
	} else {
		return extractEngNounFeatures(s, prefix)
	}
}

func ExtractNounFeaturesWithoutPrefix(s string) feature.FeatureVector {
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
