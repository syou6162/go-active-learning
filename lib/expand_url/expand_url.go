package expand_url

import (
	"fmt"

	"github.com/codegangsta/cli"
	"github.com/syou6162/go-active-learning/lib/cache"
	"github.com/syou6162/go-active-learning/lib/util"
	"os"
)

func doExpandURL(c *cli.Context) error {
	inputFilename := c.String("input-filename")

	if inputFilename == "" {
		_ = cli.ShowCommandHelp(c, "expand-url")
		return cli.NewExitError("`input-filename` is a required field.", 1)
	}

	cache, err := cache.NewCache()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
	}

	examples, err := util.ReadExamples(inputFilename)
	if err != nil {
		return err
	}

	util.AttachMetaData(cache, examples)

	for _, e := range examples {
		url := e.Url
		if e.FinalUrl != "" {
			url = e.FinalUrl
		}
		fmt.Println(fmt.Sprintf("%s\t%d", url, e.Label))
	}

	return nil
}

var CommandExpandURL = cli.Command{
	Name:  "expand-url",
	Usage: "Expand shortened url",
	Description: `
Expand shortened url.
`,
	Action: doExpandURL,
	Flags: []cli.Flag{
		cli.StringFlag{Name: "input-filename"},
	},
}
