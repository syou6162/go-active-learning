package add

import (
	"fmt"
	"log"

	"github.com/codegangsta/cli"
	"github.com/syou6162/go-active-learning/lib/hatena_bookmark"
	"github.com/syou6162/go-active-learning/lib/service"
	"github.com/syou6162/go-active-learning/lib/util/file"
)

func doAdd(c *cli.Context) error {
	inputFilename := c.String("input-filename")

	if inputFilename == "" {
		_ = cli.ShowCommandHelp(c, "add")
		return cli.NewExitError("`input-filename` is a required field.", 1)
	}

	app, err := service.NewDefaultApp()
	if err != nil {
		return err
	}
	defer app.Close()

	examples, err := file.ReadExamples(inputFilename)
	if err != nil {
		return err
	}

	app.Fetch(examples)
	app.UpdateExamplesMetadata(examples)

	m, err := app.FindLatestMIRAModel()
	skipPredictScore := false
	if err != nil {
		log.Println(fmt.Sprintf("Error to load model %s", err.Error()))
		skipPredictScore = true
	}

	for _, e := range examples {
		if !skipPredictScore {
			e.Score = m.PredictScore(e.Fv)
		}
		if err = app.InsertOrUpdateExample(e); err != nil {
			log.Println(fmt.Sprintf("Error occured proccessing %s %s", e.Url, err.Error()))
		}
		if bookmark, err := hatena_bookmark.GetHatenaBookmark(e.FinalUrl); err == nil {
			e.HatenaBookmark = bookmark
			app.UpdateExampleMetadata(*e)
		}
	}

	return nil
}

var CommandAdd = cli.Command{
	Name:  "add",
	Usage: "add urls",
	Description: `
Add urls.
`,
	Action: doAdd,
	Flags: []cli.Flag{
		cli.StringFlag{Name: "input-filename"},
	},
}
