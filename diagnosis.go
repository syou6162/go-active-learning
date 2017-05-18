package main

import (
	"fmt"
	"sort"
	"strconv"

	"github.com/codegangsta/cli"
)

var commandDiagnose = cli.Command{
	Name:  "diagnose",
	Usage: "Diagnose URLs",
	Description: `
Diagnose URLs.
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

	fmt.Println("Index\tLabel\tScore\tURL\tTitle")

	sort.Sort(sort.Reverse(wrongExamples))

	idx := 0
	for _, e := range wrongExamples {
		fmt.Println(strconv.Itoa(idx) + "\t" + strconv.Itoa(int(e.Label)) + "\t" + fmt.Sprintf("%0.03f", model.PredictScore(e.Fv)) + "\t" + e.Url + "\t" + e.Title)
		idx++
	}

	sort.Sort(correctExamples)
	for _, e := range correctExamples {
		fmt.Println(strconv.Itoa(idx) + "\t" + strconv.Itoa(int(e.Label)) + "\t" + fmt.Sprintf("%0.03f", model.PredictScore(e.Fv)) + "\t" + e.Url + "\t" + e.Title)
	}

	cache.Save(CacheFilename)
	return nil
}
