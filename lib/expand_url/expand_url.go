package expand_url

import (
	"fmt"

	"github.com/codegangsta/cli"
	"github.com/syou6162/go-active-learning/lib/cache"
	"github.com/syou6162/go-active-learning/lib/db"
	"github.com/syou6162/go-active-learning/lib/util/file"
)

func doExpandURL(c *cli.Context) error {
	inputFilename := c.String("input-filename")

	if inputFilename == "" {
		_ = cli.ShowCommandHelp(c, "expand-url")
		return cli.NewExitError("`input-filename` is a required field.", 1)
	}

	conn, err := db.CreateDBConnection()
	if err != nil {
		return err
	}
	defer conn.Close()

	cache, err := cache.NewCache()
	if err != nil {
		return err
	}
	defer cache.Close()

	examples, err := file.ReadExamples(inputFilename)
	if err != nil {
		return err
	}

	cache.AttachMetaData(examples)

	for _, e := range examples {
		_, err = db.InsertOrUpdateExample(conn, e)
		if err != nil {
			return err
		}

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
