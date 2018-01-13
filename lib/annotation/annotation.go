package annotation

import (
	"github.com/codegangsta/cli"
	"github.com/syou6162/go-active-learning/lib/classifier"
	"github.com/syou6162/go-active-learning/lib/example"
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

func rune2ActionType(r rune) ActionType {
	switch r {
	case 'p':
		return LABEL_AS_POSITIVE
	case 'n':
		return LABEL_AS_NEGATIVE
	case 's':
		return SAVE
	case 'h':
		return HELP
	case 'e':
		return EXIT
	default:
		return HELP
	}
}

func NextExampleToBeAnnotated(model classifier.BinaryClassifier, examples example.Examples) *example.Example {
	unlabeledExamples := model.SortByScore(examples)
	if len(unlabeledExamples) == 0 {
		return nil
	}
	e := unlabeledExamples[0]
	if e == nil {
		return nil
	}
	return e
}

var ActionHelpDoc = `
p: Label this example as positive.
n: Label this example as negative.
s: Save additionally annotated examples in 'output-filename'.
h: Show this help.
e: Exit.
`

var CommandAnnotate = cli.Command{
	Name:  "annotate",
	Usage: "Annotate URLs",
	Description: `
Annotate URLs using active learning.
`,
	Subcommands: []cli.Command{
		{
			Name:  "cli",
			Usage: "Annotate URLs using cli",
			Description: `
Annotate URLs using active learning using cli.
`,
			Action: doAnnotate,
			Flags: []cli.Flag{
				cli.StringFlag{Name: "input-filename"},
				cli.StringFlag{Name: "output-filename"},
				cli.BoolFlag{Name: "open-url", Usage: "Open url in background"},
				cli.BoolFlag{Name: "filter-status-code-ok", Usage: "Use only examples with status code = 200"},
				cli.BoolFlag{Name: "show-active-features"},
			},
		},
		{
			Name:  "slack",
			Usage: "Annotate URLs using slack",
			Description: `
Annotate URLs using active learning using slack.
`,
			Action: doAnnotateWithSlack,
			Flags: []cli.Flag{
				cli.StringFlag{Name: "input-filename"},
				cli.StringFlag{Name: "output-filename"},
				cli.StringFlag{Name: "channel"},
				cli.BoolFlag{Name: "filter-status-code-ok", Usage: "Use only examples with status code = 200"},
			},
		},
	},
}
