package main

import (
	"fmt"
	"os"

	"github.com/codegangsta/cli"
)

func doExpandURL(c *cli.Context) error {
	inputFilename := c.String("input-filename")

	if inputFilename == "" {
		_ = cli.ShowCommandHelp(c, "expand-url")
		return cli.NewExitError("`input-filename` is a required field.", 1)
	}

	cacheFilename := CacheFilename

	cache, err := LoadCache(cacheFilename)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
	}

	examples, err := ReadExamples(inputFilename)
	if err != nil {
		return err
	}

	AttachMetaData(cache, examples)

	for _, e := range examples {
		url := e.Url
		if e.FinalUrl != "" {
			url = e.FinalUrl
		}
		fmt.Println(fmt.Sprintf("%s\t%d", url, e.Label))
	}

	cache.Save(cacheFilename)
	return nil
}

var commandExpandURL = cli.Command{
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
