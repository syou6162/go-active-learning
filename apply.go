package main

import (
	"fmt"
	"os"

	"github.com/codegangsta/cli"
)

func doApply(c *cli.Context) error {
	inputFilename := c.String("input-filename")
	filterStatusCodeOk := c.Bool("filter-status-code-ok")

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

	alreadyLabeled := make(map[string]bool)
	for _, e := range FilterLabeledExamples(examples) {
		alreadyLabeled[e.Url] = true
	}
	for _, e := range examples {
		if _, ok := alreadyLabeled[e.Url]; ok {
			continue
		}
		if !e.IsLabeled() {
			e.Score = model.PredictScore(e.Fv)
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
	},
}
