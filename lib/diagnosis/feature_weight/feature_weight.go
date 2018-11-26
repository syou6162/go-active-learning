package featureweight

import (
	"fmt"
	"sort"

	"github.com/codegangsta/cli"
	"github.com/syou6162/go-active-learning/lib/cache"
	"github.com/syou6162/go-active-learning/lib/classifier"
	"github.com/syou6162/go-active-learning/lib/repository"
	"github.com/syou6162/go-active-learning/lib/util"
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

	err := cache.Init()
	if err != nil {
		return err
	}
	defer cache.Close()

	repo, err := repository.New()
	if err != nil {
		return err
	}
	defer repo.Close()

	examples, err := repo.ReadExamples()
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
