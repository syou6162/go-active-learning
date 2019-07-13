package diagnosis

import (
	"github.com/codegangsta/cli"
	featureweight "github.com/syou6162/go-active-learning/lib/diagnosis/feature_weight"
	labelconflict "github.com/syou6162/go-active-learning/lib/diagnosis/label_conflict"
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
			Action: labelconflict.DoLabelConflict,
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
			Action: featureweight.DoListFeatureWeight,
			Flags: []cli.Flag{
				cli.BoolFlag{Name: "filter-status-code-ok", Usage: "Use only examples with status code = 200"},
			},
		},
	},
}
