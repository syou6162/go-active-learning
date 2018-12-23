package add

import (
	"fmt"
	"log"

	"github.com/codegangsta/cli"
	"github.com/syou6162/go-active-learning/lib/hatena_bookmark"
	"github.com/syou6162/go-active-learning/lib/service"
	"github.com/syou6162/go-active-learning/lib/util"
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

	if err := app.AttachLightMetadata(examples); err != nil {
		return err
	}

	examples = util.FilterStatusCodeNotOkExamples(examples)
	app.Fetch(examples)
	examples = util.FilterStatusCodeOkExamples(examples)

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
		if err = app.UpdateOrCreateExample(e); err != nil {
			log.Println(fmt.Sprintf("Error occured proccessing %s %s", e.Url, err.Error()))
			continue
		}
		if err = app.UpdateFeatureVector(e); err != nil {
			log.Println(fmt.Sprintf("Error occured proccessing %s feature vector %s", e.Url, err.Error()))
			continue
		}
		if bookmark, err := hatena_bookmark.GetHatenaBookmark(e.FinalUrl); err == nil {
			e.HatenaBookmark = bookmark
			app.UpdateHatenaBookmark(e)
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
