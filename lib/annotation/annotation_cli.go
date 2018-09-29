package annotation

import (
	"fmt"
	"os"

	"math"
	"sort"

	"github.com/codegangsta/cli"
	"github.com/mattn/go-tty"
	"github.com/pkg/browser"
	"github.com/syou6162/go-active-learning/lib/cache"
	"github.com/syou6162/go-active-learning/lib/classifier"
	"github.com/syou6162/go-active-learning/lib/db"
	"github.com/syou6162/go-active-learning/lib/example"
	"github.com/syou6162/go-active-learning/lib/util"
)

func input2ActionType() (ActionType, error) {
	t, err := tty.Open()
	defer t.Close()
	if err != nil {
		return EXIT, err
	}
	var r rune
	for r == 0 {
		r, err = t.ReadRune()
		if err != nil {
			return HELP, err
		}
	}
	return rune2ActionType(r), nil
}

func doAnnotate(c *cli.Context) error {
	openUrl := c.Bool("open-url")
	filterStatusCodeOk := c.Bool("filter-status-code-ok")
	showActiveFeatures := c.Bool("show-active-features")

	cache, err := cache.NewCache()
	if err != nil {
		return err
	}
	defer cache.Close()

	conn, err := db.CreateDBConnection()
	if err != nil {
		return err
	}
	defer conn.Close()

	examples, err := db.ReadExamples(conn)
	if err != nil {
		return err
	}

	stat := example.GetStat(examples)
	fmt.Fprintln(os.Stderr, fmt.Sprintf("Positive:%d, Negative:%d, Unlabeled:%d", stat["positive"], stat["negative"], stat["unlabeled"]))

	cache.AttachMetaData(examples, true)
	if filterStatusCodeOk {
		examples = util.FilterStatusCodeOkExamples(examples)
	}
	model := classifier.NewBinaryClassifier(examples)

annotationLoop:
	for {
		e := NextExampleToBeAnnotated(model, examples)
		if e == nil {
			fmt.Println("No example")
			break annotationLoop
		}
		fmt.Println("Label this example (Score: " + fmt.Sprintf("%+0.03f", e.Score) + "): " + e.Url + " (" + e.Title + ")")

		if openUrl {
			browser.OpenURL(e.Url)
		}
		if showActiveFeatures {
			ShowActiveFeatures(model, *e, 5)
		}

		act, err := input2ActionType()
		if err != nil {
			return err
		}
		switch act {
		case LABEL_AS_POSITIVE:
			fmt.Println("Labeled as positive")
			e.Annotate(example.POSITIVE)
			db.InsertOrUpdateExample(conn, e)
		case LABEL_AS_NEGATIVE:
			fmt.Println("Labeled as negative")
			e.Annotate(example.NEGATIVE)
			db.InsertOrUpdateExample(conn, e)
		case SKIP:
			fmt.Println("Skiped this example")
			examples = util.RemoveExample(examples, *e)
			continue
		case HELP:
			fmt.Println(ActionHelpDoc)
		case EXIT:
			fmt.Println("EXIT")
			break annotationLoop
		default:
			break annotationLoop
		}
		model = classifier.NewBinaryClassifier(examples)
	}

	return nil
}

type FeatureWeightPair struct {
	Feature string
	Weight  float64
}

type FeatureWeightPairs []FeatureWeightPair

func SortedActiveFeatures(model classifier.BinaryClassifier, example example.Example, n int) FeatureWeightPairs {
	pairs := FeatureWeightPairs{}
	for _, f := range example.Fv {
		pairs = append(pairs, FeatureWeightPair{f, model.GetWeight(f)})
	}
	sort.Sort(sort.Reverse(pairs))

	result := FeatureWeightPairs{}
	cnt := 0
	for _, pair := range pairs {
		if cnt >= n {
			break
		}
		if (example.Score > 0.0 && pair.Weight > 0.0) || (example.Score < 0.0 && pair.Weight < 0.0) {
			result = append(result, pair)
			cnt++
		}
	}
	return result
}

func ShowActiveFeatures(model classifier.BinaryClassifier, example example.Example, n int) {
	for _, pair := range SortedActiveFeatures(model, example, n) {
		fmt.Println(fmt.Sprintf("%+0.1f %s", pair.Weight, pair.Feature))
	}
}

func (slice FeatureWeightPairs) Len() int {
	return len(slice)
}

func (slice FeatureWeightPairs) Less(i, j int) bool {
	return math.Abs(slice[i].Weight) < math.Abs(slice[j].Weight)
}

func (slice FeatureWeightPairs) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}
