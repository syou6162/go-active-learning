package featureweight

import (
	"fmt"
	"sort"

	"github.com/codegangsta/cli"
	"github.com/syou6162/go-active-learning/lib/classifier"
	"github.com/syou6162/go-active-learning/lib/service"
	"github.com/syou6162/go-active-learning/lib/util"
	"github.com/syou6162/go-active-learning/lib/util/converter"
)

type Feature struct {
	Key    string
	Weight float64
}

type FeatureList []Feature

func (p FeatureList) Len() int           { return len(p) }
func (p FeatureList) Less(i, j int) bool { return p[i].Weight < p[j].Weight }
func (p FeatureList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

func DoListFeatureWeight(c *cli.Context) error {
	filterStatusCodeOk := c.Bool("filter-status-code-ok")

	app, err := service.NewDefaultApp()
	if err != nil {
		return err
	}
	defer app.Close()

	examples, err := app.SearchExamples()
	if err != nil {
		return err
	}
	app.Fetch(examples)
	for _, e := range examples {
		app.UpdateFeatureVector(e)
	}
	training := util.FilterLabeledExamples(examples)

	if filterStatusCodeOk {
		training = util.FilterStatusCodeOkExamples(training)
	}

	model, err := classifier.NewMIRAClassifierByCrossValidation(converter.ConvertExamplesToLearningInstances(examples))
	if err != nil {
		return err
	}

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
