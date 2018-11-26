package labelconflict

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
	"github.com/syou6162/go-active-learning/lib/model"
	"github.com/syou6162/go-active-learning/lib/util"
)

func DoLabelConflict(c *cli.Context) error {
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

	m := classifier.NewBinaryClassifier(training)

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
	printResult(m, correctExamples, wrongExamples)

	return nil
}

func printResult(m classifier.BinaryClassifier, correctExamples model.Examples, wrongExamples model.Examples) error {
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
