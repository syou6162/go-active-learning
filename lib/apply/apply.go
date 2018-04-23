package apply

import (
	"fmt"
	"strings"

	"encoding/json"

	"github.com/codegangsta/cli"
	"github.com/syou6162/go-active-learning/lib/cache"
	"github.com/syou6162/go-active-learning/lib/classifier"
	"github.com/syou6162/go-active-learning/lib/db"
	"github.com/syou6162/go-active-learning/lib/example"
	"github.com/syou6162/go-active-learning/lib/submodular"
	"github.com/syou6162/go-active-learning/lib/util"
	"github.com/syou6162/go-active-learning/lib/util/file"
)

func doApply(c *cli.Context) error {
	inputFilename := c.String("input-filename")
	filterStatusCodeOk := c.Bool("filter-status-code-ok")
	jsonOutput := c.Bool("json-output")
	subsetSelection := c.Bool("subset-selection")
	sizeConstraint := c.Int("size-constraint")
	alpha := c.Float64("alpha")
	r := c.Float64("r")
	scoreThreshold := c.Float64("score-threshold")

	if inputFilename == "" {
		_ = cli.ShowCommandHelp(c, "apply")
		return cli.NewExitError("`input-filename` is a required field.", 1)
	}

	cache, err := cache.NewCache()
	if err != nil {
		return err
	}
	defer cache.Close()

	conn, err := db.CreateDBConnection()
	if err != nil {
		return err
	}

	examples, err := db.ReadExamples(conn)
	if err != nil {
		return err
	}

	cache.AttachMetaData(examples)
	if filterStatusCodeOk {
		examples = util.FilterStatusCodeOkExamples(examples)
	}
	model := classifier.NewBinaryClassifier(examples)

	examplesFromFile, err := file.ReadExamples(inputFilename)
	if err != nil {
		return err
	}
	cache.AttachMetaData(examplesFromFile)

	result := example.Examples{}
	for _, e := range util.FilterUnlabeledExamples(examplesFromFile) {
		e.Score = model.PredictScore(e.Fv)
		e.Title = strings.Replace(e.Title, "\n", " ", -1)
		if e.Score > scoreThreshold {
			result = append(result, e)
		}
	}

	if subsetSelection {
		result = submodular.SelectSubExamplesBySubModular(result, sizeConstraint, alpha, r)
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

	return nil
}

var CommandApply = cli.Command{
	Name:  "apply",
	Usage: "Apply classifier to unlabeled examples",
	Description: `
Apply classifier to unlabeled examples, and print a pair of score and url.
`,
	Action: doApply,
	Flags: []cli.Flag{
		cli.StringFlag{Name: "input-filename"},
		cli.BoolFlag{Name: "filter-status-code-ok", Usage: "Use only examples with status code = 200"},
		cli.BoolFlag{Name: "json-output", Usage: "Make output with json format or not (tsv format)."},
		cli.BoolFlag{Name: "subset-selection", Usage: "Use subset selection algorithm (maximizing submodular function) to filter entries"},
		cli.Int64Flag{Name: "size-constraint", Value: 10, Usage: "Budget constraint. Max number of entries to be contained"},
		cli.Float64Flag{Name: "alpha", Value: 1.0},
		cli.Float64Flag{Name: "r", Value: 1.0, Usage: "Scaling factor for number of words"},
		cli.Float64Flag{Name: "score-threshold", Value: 0.0},
	},
}
