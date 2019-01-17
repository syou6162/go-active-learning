package util

import (
	"os"

	"github.com/syou6162/go-active-learning/lib/model"
	"gopkg.in/vmarkovtsev/go-lcss.v1"
)

func FilterLabeledExamples(examples model.Examples) model.Examples {
	var result model.Examples
	for _, e := range examples {
		if e.IsLabeled() {
			result = append(result, e)
		}
	}
	return result
}

func FilterUnlabeledExamples(examples model.Examples) model.Examples {
	result := model.Examples{}

	alreadyLabeledByURL := make(map[string]bool)
	alreadyLabeledByTitle := make(map[string]bool)
	for _, e := range FilterLabeledExamples(examples) {
		alreadyLabeledByURL[e.Url] = true
		alreadyLabeledByURL[e.FinalUrl] = true
		alreadyLabeledByTitle[e.Title] = true
	}

	for _, e := range examples {
		if _, ok := alreadyLabeledByURL[e.Url]; ok {
			continue
		}
		if _, ok := alreadyLabeledByTitle[e.Title]; ok {
			continue
		}
		if !e.IsLabeled() {
			alreadyLabeledByURL[e.Url] = true
			alreadyLabeledByURL[e.FinalUrl] = true
			alreadyLabeledByTitle[e.Title] = true
			result = append(result, e)
		}
	}
	return result
}

func RemoveDuplicate(args []string) []string {
	results := make([]string, 0)
	encountered := map[string]bool{}
	for i := 0; i < len(args); i++ {
		if !encountered[args[i]] {
			encountered[args[i]] = true
			results = append(results, args[i])
		}
	}
	return results
}

func FilterStatusCodeOkExamples(examples model.Examples) model.Examples {
	result := model.Examples{}

	for _, e := range examples {
		if e.StatusCode == 200 {
			result = append(result, e)
		}
	}

	return result
}

func FilterStatusCodeNotOkExamples(examples model.Examples) model.Examples {
	result := model.Examples{}

	for _, e := range examples {
		if e.StatusCode != 200 {
			result = append(result, e)
		}
	}

	return result
}

func RemoveExample(examples model.Examples, toBeRemoved model.Example) model.Examples {
	result := model.Examples{}

	for _, e := range examples {
		if e.Url != toBeRemoved.Url {
			result = append(result, e)
		}
	}

	return result
}

func RemoveNegativeExamples(examples model.Examples) model.Examples {
	result := model.Examples{}
	for _, e := range examples {
		if e.Label != model.NEGATIVE {
			result = append(result, e)
		}
	}
	return result
}

func UniqueByFinalUrl(examples model.Examples) model.Examples {
	result := model.Examples{}
	m := make(map[string]bool)
	for _, e := range examples {
		if !m[e.FinalUrl] {
			m[e.FinalUrl] = true
			result = append(result, e)
		}
	}
	return result
}

func UniqueByTitle(examples model.Examples) model.Examples {
	result := model.Examples{}
	m := make(map[string]bool)
	for _, e := range examples {
		if !m[e.Title] {
			m[e.Title] = true
			result = append(result, e)
		}
	}
	return result
}

func SplitTrainAndDev(examples model.Examples) (train model.Examples, dev model.Examples) {
	Shuffle(examples)
	n := int(0.8 * float64(len(examples)))
	return examples[0:n], examples[n:]
}

func GetEnv(key, fallback string) string {
	value, ok := os.LookupEnv(key)
	if !ok {
		value = fallback
	}
	return value
}

func SortByCommentUsefulness(e model.Example, tweets model.ReferringTweets) model.ReferringTweets {
	usefulTweets := model.ReferringTweets{}
	otherwise := model.ReferringTweets{}

	for _, t := range tweets {
		origLen := len(t.FullText)
		commonLen := len(string(lcss.LongestCommonSubstring([]byte(e.Title), []byte(t.FullText))))
		if float64(commonLen)/float64(origLen) < 0.5 {
			usefulTweets = append(usefulTweets, t)
		} else {
			otherwise = append(otherwise, t)
		}
	}

	result := model.ReferringTweets{}
	result = append(result, usefulTweets...)
	result = append(result, otherwise...)
	return result
}
