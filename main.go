package main

import (
	"fmt"
	"os"
	"runtime"
	"strconv"
	"sync"

	"github.com/codegangsta/cli"
	"github.com/mattn/go-tty"
	"github.com/pkg/browser"
)

type ActionType int

const (
	LABEL_AS_POSITIVE ActionType = iota
	LABEL_AS_NEGATIVE
	SAVE
	HELP
	SKIP
	EXIT
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
			return SKIP, err
		}
	}
	switch r {
	case 'p':
		return LABEL_AS_POSITIVE, nil
	case 'n':
		return LABEL_AS_NEGATIVE, nil
	case 's':
		return SAVE, nil
	case 'h':
		return HELP, nil
	case 'e':
		return EXIT, nil
	default:
		return SKIP, nil
	}
}

func doAnnotate(c *cli.Context) error {
	inputFilename := c.String("input-filename")
	outputFilename := c.String("output-filename")

	if inputFilename == "" {
		_ = cli.ShowCommandHelp(c, "annotate")
		return cli.NewExitError("`input-filename` is a required field.", 1)
	}

	if outputFilename == "" {
		_ = cli.ShowCommandHelp(c, "annotate")
		return cli.NewExitError("`output-filename` is a required field.", 1)
	}

	cacheFilename := "cache.bin"

	cache, _ := LoadCache(cacheFilename)
	examples, _ := ReadExamples(inputFilename)

	shuffle(examples)

	wg := &sync.WaitGroup{}
	cpus := runtime.NumCPU()
	runtime.GOMAXPROCS(cpus)
	sem := make(chan struct{}, cpus)
	for idx, e := range examples {
		wg.Add(1)
		sem <- struct{}{}
		go func(e *Example, idx int) {
			defer wg.Done()
			if example, ok := cache.Cache[e.Url]; ok {
				e.Title = example.Title
				e.Description = example.Description
				e.Body = example.Body
				e.RawHTML = example.RawHTML
			} else {
				article := GetArticle(e.Url)
				fmt.Println("Fetching(" + strconv.Itoa(idx) + "): " + e.Url)
				e.Title = article.Title
				e.Description = article.Description
				e.Body = article.Body
				e.RawHTML = article.RawHTML
				cache.Add(*e)
			}
			e.Fv = removeDuplicate(ExtractFeatures(*e))
			<-sem
		}(e, idx)
	}
	wg.Wait()

	model := TrainedModel(examples)

annotationLoop:
	for {
		unlabeledExamples := model.SortByScore(examples)
		if len(unlabeledExamples) == 0 {
			break
		}
		e := unlabeledExamples[0]
		if e == nil {
			break
		}
		fmt.Println("Label this example (Score: " + fmt.Sprintf("%0.03f", e.Score) + "): " + e.Url + " (" + e.Title + ")")
		browser.OpenURL(e.Url)
		cache.Add(*e)

		act, err := input2ActionType()
		if err != nil {
			return err
		}
		switch act {
		case LABEL_AS_POSITIVE:
			fmt.Println("Labeled as positive")
			e.Annotate(POSITIVE)
		case LABEL_AS_NEGATIVE:
			fmt.Println("Labeled as negative")
			e.Annotate(NEGATIVE)
		case SKIP:
			fmt.Println("Skiped this example")
			continue
		case SAVE:
			fmt.Println("Saved labeld examples")
			WriteExamples(examples, outputFilename)
		case HELP:
			fmt.Println("ToDo: SHOW HELP")
		case EXIT:
			fmt.Println("EXIT")
			break annotationLoop
		default:
			break annotationLoop
		}
		model = TrainedModel(examples)
	}

	WriteExamples(examples, outputFilename)
	cache.Save(cacheFilename)

	return nil
}

var commandAnnotate = cli.Command{
	Name:  "annotate",
	Usage: "Annotate URLs",
	Description: `
Annotate URLs using active learning.
`,
	Action: doAnnotate,
	Flags: []cli.Flag{
		cli.StringFlag{Name: "input-filename"},
		cli.StringFlag{Name: "output-filename"},
	},
}

var Commands = []cli.Command{
	commandAnnotate,
}

func main() {
	app := cli.NewApp()
	app.Name = "go-active-learning"
	app.Commands = Commands

	app.Run(os.Args)
}
