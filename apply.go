package main

import (
	"fmt"
	"os"
	"strings"

	"encoding/json"

	"github.com/codegangsta/cli"
)

func doApply(c *cli.Context) error {
	inputFilename := c.String("input-filename")
	filterStatusCodeOk := c.Bool("filter-status-code-ok")
	jsonOutput := c.Bool("json-output")
	subsetSelection := c.Bool("subset-selection")
	sizeConstraint := c.Int("size-constraint")
	alpha := c.Float64("alpha")
	r := c.Float64("r")

	if inputFilename == "" {
		_ = cli.ShowCommandHelp(c, "apply")
		return cli.NewExitError("`input-filename` is a required field.", 1)
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

	AttachMetaData(cache, examples)
	if filterStatusCodeOk {
		examples = FilterStatusCodeOkExamples(examples)
	}
	model := TrainedModel(examples)

	result := Examples{}
	for _, e := range FilterUnlabeledExamples(examples) {
		e.Score = model.PredictScore(e.Fv)
		e.Title = strings.Replace(e.Title, "\n", " ", -1)
		if e.Score > 0.0 {
			result = append(result, e)
		}
	}

	if subsetSelection {
		result = SelectSubExamplesBySubModular(model, result, sizeConstraint, alpha, r)
	}

	for _, e := range result {
		e.Score = model.PredictScore(e.Fv)
		e.Title = strings.Replace(e.Title, "\n", " ", -1)
		if jsonOutput {
			b, err := json.Marshal(e)
			if err != nil {
				return err
			}
			fmt.Println(string(b))
		} else {
			fmt.Println(fmt.Sprintf("%0.03f\t%s", e.Score, e.Url))
		}
	}

	cache.Save(cacheFilename)
	return nil
}

var commandApply = cli.Command{
	Name:  "apply",
	Usage: "apply classifier to unlabeled examples",
	Description: `
Apply classifier to unlabeled examples, and print a pair of score and url.
`,
	Action: doApply,
	Flags: []cli.Flag{
		cli.StringFlag{Name: "input-filename"},
		cli.BoolFlag{Name: "filter-status-code-ok", Usage: "Use only examples with status code = 200"},
		cli.BoolFlag{Name: "subset-selection"},
		cli.BoolFlag{Name: "json-output"},
		cli.Int64Flag{Name: "size-constraint", Value: 10},
		cli.Float64Flag{Name: "alpha", Value: 1.0},
		cli.Float64Flag{Name: "r", Value: 1.0},
	},
}
