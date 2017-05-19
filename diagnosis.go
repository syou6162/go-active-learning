package main

import (
	"fmt"
	"os"
	"sort"
	"strconv"

	"encoding/csv"

	"github.com/codegangsta/cli"
)

var commandDiagnose = cli.Command{
	Name:  "diagnose",
	Usage: "Diagnose label conflicts in training data",
	Description: `
Diagnose label conflicts in training data. 'conflict' means that an annotated label is '-1/1', but a predicted label by model is '1/-1'.
`,
	Action: doDiagnose,
	Flags: []cli.Flag{
		cli.StringFlag{Name: "input-filename"},
		cli.BoolFlag{Name: "filter-status-code-ok", Usage: "Use only examples with status code = 200"},
	},
}

func doDiagnose(c *cli.Context) error {
	inputFilename := c.String("input-filename")
	filterStatusCodeOk := c.Bool("filter-status-code-ok")

	if inputFilename == "" {
		_ = cli.ShowCommandHelp(c, "diagnose")
	}

	cache, _ := LoadCache(CacheFilename)
	examples, _ := ReadExamples(inputFilename)
	AttachMetaData(cache, examples)
	training := FilterLabeledExamples(examples)

	if filterStatusCodeOk {
		training = FilterStatusCodeOkExamples(training)
	}

	model := TrainedModel(training)

	wrongExamples := Examples{}
	correctExamples := Examples{}

	for _, e := range training {
		e.Score = model.PredictScore(e.Fv)
		if float64(e.Label)*e.Score < 0 {
			wrongExamples = append(wrongExamples, e)
		} else {
			correctExamples = append(correctExamples, e)
		}
	}

	sort.Sort(sort.Reverse(wrongExamples))
	sort.Sort(correctExamples)
	printResult(model, correctExamples, wrongExamples)

	cache.Save(CacheFilename)
	return nil
}

func printResult(model *Model, correctExamples Examples, wrongExamples Examples) error {
	fmt.Println("Index\tLabel\tScore\tURL\tTitle")
	result := append(wrongExamples, correctExamples...)

	w := csv.NewWriter(os.Stdout)
	w.Comma = '\t'

	for idx, e := range result {
		record := []string{
			strconv.Itoa(idx),
			strconv.Itoa(int(e.Label)),
			fmt.Sprintf("%0.03f", model.PredictScore(e.Fv)),
			e.Url,
			e.Title,
		}
		if err := w.Write(record); err != nil {
			return err
		}
	}

	w.Flush()
	if err := w.Error(); err != nil {
		return err
	}

	return nil
}
