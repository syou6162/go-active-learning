package labelconflict

import (
	"fmt"
	"os"
	"sort"
	"strconv"

	"encoding/csv"

	"github.com/syou6162/go-active-learning/lib/classifier"
	"github.com/syou6162/go-active-learning/lib/model"
	"github.com/syou6162/go-active-learning/lib/service"
	"github.com/syou6162/go-active-learning/lib/util"
	"github.com/syou6162/go-active-learning/lib/util/converter"
	"github.com/urfave/cli"
)

func DoLabelConflict(c *cli.Context) error {
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

	m, err := classifier.NewMIRAClassifierByCrossValidation(classifier.EXAMPLE, converter.ConvertExamplesToLearningInstances(training))
	if err != nil {
		return err
	}

	wrongExamples := model.Examples{}
	correctExamples := model.Examples{}

	for _, e := range training {
		e.Score = m.PredictScore(e.Fv)
		if float64(e.Label)*e.Score < 0 {
			wrongExamples = append(wrongExamples, e)
		} else {
			correctExamples = append(correctExamples, e)
		}
	}

	sort.Sort(sort.Reverse(wrongExamples))
	sort.Sort(correctExamples)
	printResult(*m, correctExamples, wrongExamples)

	return nil
}

func printResult(m classifier.MIRAClassifier, correctExamples model.Examples, wrongExamples model.Examples) error {
	fmt.Println("Index\tLabel\tScore\tURL\tTitle")
	result := append(wrongExamples, correctExamples...)

	w := csv.NewWriter(os.Stdout)
	w.Comma = '\t'

	for idx, e := range result {
		record := []string{
			strconv.Itoa(idx),
			strconv.Itoa(int(e.Label)),
			fmt.Sprintf("%0.03f", m.PredictScore(e.Fv)),
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
