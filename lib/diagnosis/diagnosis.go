package diagnosis

import (
	"fmt"
	"os"
	"sort"
	"strconv"

	"encoding/csv"

	"github.com/codegangsta/cli"
	"github.com/syou6162/go-active-learning/lib/cache"
	"github.com/syou6162/go-active-learning/lib/classifier"
	"github.com/syou6162/go-active-learning/lib/db"
	"github.com/syou6162/go-active-learning/lib/example"
	"github.com/syou6162/go-active-learning/lib/util"
)

var CommandDiagnose = cli.Command{
	Name:  "diagnose",
	Usage: "Diagnose training data or learned model",
	Description: `
Diagnose training data or learned model. This mode has two subcommand: label-conflict and feature-weight.
`,

	Subcommands: []cli.Command{
		{
			Name:  "label-conflict",
			Usage: "Diagnose label conflicts in training data",
			Description: `
Diagnose label conflicts in training data. 'conflict' means that an annotated label is '-1/1', but a predicted label by model is '1/-1'.
`,
			Action: doDiagnose,
			Flags: []cli.Flag{
				cli.BoolFlag{Name: "filter-status-code-ok", Usage: "Use only examples with status code = 200"},
			},
		},
		{
			Name:  "feature-weight",
			Usage: "List feature weight",
			Description: `
List feature weight.
`,
			Action: doListFeatureWeight,
			Flags: []cli.Flag{
				cli.BoolFlag{Name: "filter-status-code-ok", Usage: "Use only examples with status code = 200"},
			},
		},
	},
}

func doDiagnose(c *cli.Context) error {
	filterStatusCodeOk := c.Bool("filter-status-code-ok")

	err := cache.Init()
	if err != nil {
		return err
	}
	defer cache.Close()

	err = db.Init()
	if err != nil {
		return err
	}
	defer db.Close()

	examples, err := db.ReadExamples()
	if err != nil {
		return err
	}
	cache.AttachMetadata(examples, true, false)
	training := util.FilterLabeledExamples(examples)

	if filterStatusCodeOk {
		training = util.FilterStatusCodeOkExamples(training)
	}

	model := classifier.NewBinaryClassifier(training)

	wrongExamples := example.Examples{}
	correctExamples := example.Examples{}

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

	return nil
}

func printResult(model classifier.BinaryClassifier, correctExamples example.Examples, wrongExamples example.Examples) error {
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

type Feature struct {
	Key    string
	Weight float64
}

type FeatureList []Feature

func (p FeatureList) Len() int           { return len(p) }
func (p FeatureList) Less(i, j int) bool { return p[i].Weight < p[j].Weight }
func (p FeatureList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

func doListFeatureWeight(c *cli.Context) error {
	filterStatusCodeOk := c.Bool("filter-status-code-ok")

	err := cache.Init()
	if err != nil {
		return err
	}
	defer cache.Close()

	err = db.Init()
	if err != nil {
		return err
	}
	defer db.Close()

	examples, err := db.ReadExamples()
	if err != nil {
		return err
	}
	cache.AttachMetadata(examples, true, false)
	training := util.FilterLabeledExamples(examples)

	if filterStatusCodeOk {
		training = util.FilterStatusCodeOkExamples(training)
	}

	model := classifier.NewBinaryClassifier(training)

	tmp := make(FeatureList, 0)
	for _, k := range model.GetActiveFeatures() {
		tmp = append(tmp, Feature{k, model.GetWeight(k)})
	}
	sort.Sort(sort.Reverse(tmp))

	for _, p := range tmp {
		fmt.Println(fmt.Sprintf("%+0.2f\t%s", p.Weight, p.Key))
	}

	return nil
}
