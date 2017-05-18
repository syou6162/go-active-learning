package main

import (
	"fmt"

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

var ActionHelpDoc = `
p: Label this example as positive.
n: Label this example as negative.
s: Save additionally annotated examples in 'output-filename'.
h: Show this help.
e: Exit.
`

func doAnnotate(c *cli.Context) error {
	inputFilename := c.String("input-filename")
	outputFilename := c.String("output-filename")
	openUrl := c.Bool("openurl")
	filterStatusCodeOk := c.Bool("filter-status-code-ok")

	if inputFilename == "" {
		_ = cli.ShowCommandHelp(c, "annotate")
		return cli.NewExitError("`input-filename` is a required field.", 1)
	}

	if outputFilename == "" {
		_ = cli.ShowCommandHelp(c, "annotate")
		return cli.NewExitError("`output-filename` is a required field.", 1)
	}

	cacheFilename := CacheFilename

	cache, _ := LoadCache(cacheFilename)
	examples, _ := ReadExamples(inputFilename)
	AttachMetaData(cache, examples)
	if filterStatusCodeOk {
		examples = FilterStatusCodeOkExamples(examples)
	}
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
		cache.Add(*e)

		if openUrl {
			browser.OpenURL(e.Url)
		}

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
			fmt.Println(ActionHelpDoc)
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
		cli.BoolFlag{Name: "openurl", Usage: "Open url in background"},
		cli.BoolFlag{Name: "filter-status-code-ok", Usage: "Use only examples with status code = 200"},
	},
}
