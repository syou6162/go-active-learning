package util

import (
	"fmt"
	"os"
	"time"

	"github.com/syou6162/go-active-learning/lib/example"
)

func FilterLabeledExamples(examples example.Examples) example.Examples {
	var result example.Examples
	for _, e := range examples {
		if e.IsLabeled() {
			result = append(result, e)
		}
	}
	return result
}

func FilterUnlabeledExamples(examples example.Examples) example.Examples {
	result := example.Examples{}

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

func FilterStatusCodeOkExamples(examples example.Examples) example.Examples {
	result := example.Examples{}

	for _, e := range examples {
		if e.StatusCode == 200 {
			result = append(result, e)
		}
	}

	return result
}

func RemoveExample(examples example.Examples, toBeRemoved example.Example) example.Examples {
	result := example.Examples{}

	for _, e := range examples {
		if e.Url != toBeRemoved.Url {
			result = append(result, e)
		}
	}

	return result
}

func RemoveNegativeExamples(examples example.Examples) example.Examples {
	result := example.Examples{}
	for _, e := range examples {
		if e.Label != example.NEGATIVE {
			result = append(result, e)
		}
	}
	return result
}

func SplitTrainAndDev(examples example.Examples) (train example.Examples, dev example.Examples) {
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
