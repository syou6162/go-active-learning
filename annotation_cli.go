package main

import (
	"fmt"
	"os"

	"math"
	"sort"

	"github.com/codegangsta/cli"
	"github.com/mattn/go-tty"
	"github.com/pkg/browser"
)

func input2ActionType() (ActionType, error) {
	t, err := tty.Open()
	defer t.Close()
	if err != nil {
		return EXIT, err
	}
	var r rune
	for r == 0 {
		r, err = t.ReadRune()
		if err != nil {
			return HELP, err
		}
	}
	return rune2ActionType(r), nil
}

func doAnnotate(c *cli.Context) error {
	inputFilename := c.String("input-filename")
	outputFilename := c.String("output-filename")
	openUrl := c.Bool("open-url")
	filterStatusCodeOk := c.Bool("filter-status-code-ok")
	showActiveFeatures := c.Bool("show-active-features")

	if inputFilename == "" {
		_ = cli.ShowCommandHelp(c, "cli")
		return cli.NewExitError("`input-filename` is a required field.", 1)
	}

	if outputFilename == "" {
		outputFilename = NewOutputFilename()
		fmt.Fprintln(os.Stderr, "'output-filename' is not specified. "+outputFilename+" is used as output-filename instead.")
	}

	cacheFilename := CacheFilename

	cache, err := LoadCache(cacheFilename)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
	}

	examples, err := ReadExamples(inputFilename)
	if err != nil {
		return err
	}

	stat := GetStat(examples)
	fmt.Fprintln(os.Stderr, fmt.Sprintf("Positive:%d, Negative:%d, Unlabeled:%d", stat["positive"], stat["negative"], stat["unlabeled"]))

	AttachMetaData(cache, examples)
	if filterStatusCodeOk {
		examples = FilterStatusCodeOkExamples(examples)
	}
	model := TrainedModel(examples)

annotationLoop:
	for {
		e := NextExampleToBeAnnotated(model, examples)
		fmt.Println("Label this example (Score: " + fmt.Sprintf("%+0.03f", e.Score) + "): " + e.Url + " (" + e.Title + ")")

		if openUrl {
			browser.OpenURL(e.Url)
		}
		if showActiveFeatures {
			ShowActiveFeatures(model, *e, 5)
		}

		act, err := input2ActionType()
		if err != nil {
			return err
		}
		switch act {
		case LABEL_AS_POSITIVE:
			fmt.Println("Labeled as positive")
			e.Annotate(POSITIVE)
		case LABEL_AS_NEGATIVE:
			fmt.Println("Labeled as negative")
			e.Annotate(NEGATIVE)
		case SKIP:
			fmt.Println("Skiped this example")
			continue
		case SAVE:
			fmt.Println("Saved labeld examples")
			WriteExamples(examples, outputFilename)
		case HELP:
			fmt.Println(ActionHelpDoc)
		case EXIT:
			fmt.Println("EXIT")
			break annotationLoop
		default:
			break annotationLoop
		}
		model = TrainedModel(examples)
	}

	WriteExamples(examples, outputFilename)
	cache.Save(cacheFilename)

	return nil
}

type FeatureWeightPair struct {
	Feature string
	Weight  float64
}

type FeatureWeightPairs []FeatureWeightPair

func SortedActiveFeatures(model *Model, example Example, n int) FeatureWeightPairs {
	pairs := FeatureWeightPairs{}
	for _, f := range example.Fv {
		pairs = append(pairs, FeatureWeightPair{f, model.GetAveragedWeight(f)})
	}
	sort.Sort(sort.Reverse(pairs))

	result := FeatureWeightPairs{}
	cnt := 0
	for _, pair := range pairs {
		if cnt >= n {
			break
		}
		if (example.Score > 0.0 && pair.Weight > 0.0) || (example.Score < 0.0 && pair.Weight < 0.0) {
			result = append(result, pair)
			cnt++
		}
	}
	return result
}

func ShowActiveFeatures(model *Model, example Example, n int) {
	for _, pair := range SortedActiveFeatures(model, example, n) {
		fmt.Println(fmt.Sprintf("%+0.1f %s", pair.Weight, pair.Feature))
	}
}

func (slice FeatureWeightPairs) Len() int {
	return len(slice)
}

func (slice FeatureWeightPairs) Less(i, j int) bool {
	return math.Abs(slice[i].Weight) < math.Abs(slice[j].Weight)
}

func (slice FeatureWeightPairs) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}
