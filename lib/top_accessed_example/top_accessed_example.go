package top_accessed_example

import (
	"bufio"
	"fmt"
	"strconv"

	"os"

	"github.com/syou6162/go-active-learning/lib/service"
	"github.com/urfave/cli"
)

func parseLine(line string) (int, error) {
	exampleId, err := strconv.ParseInt(line, 10, 0)
	if err != nil {
		return 0, fmt.Errorf("Invalid line: %s", line)
	}
	return int(exampleId), nil
}

func readTopAccessedExampleIds(filename string) ([]int, error) {
	fp, err := os.Open(filename)
	defer fp.Close()
	if err != nil {
		return nil, err
	}

	exampleIds := make([]int, 0)
	scanner := bufio.NewScanner(fp)
	for scanner.Scan() {
		line := scanner.Text()
		exampleId, err := parseLine(line)
		if err != nil {
			return nil, err
		}
		exampleIds = append(exampleIds, exampleId)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return exampleIds, nil
}

func doAddTopAccessedExamples(c *cli.Context) error {
	inputFilename := c.String("input-filename")

	if inputFilename == "" {
		_ = cli.ShowCommandHelp(c, "add-top-accessed-examples")
		return cli.NewExitError("`input-filename` is a required field.", 1)
	}

	app, err := service.NewDefaultApp()
	if err != nil {
		return err
	}
	defer app.Close()

	exampleIds, err := readTopAccessedExampleIds(inputFilename)
	if err != nil {
		return err
	}
	err = app.UpdateTopAccessedExampleIds(exampleIds)
	if err != nil {
		return err
	}
	return nil
}

var CommandAddTopAccessedExamples = cli.Command{
	Name:  "add-top-accessed-examples",
	Usage: "add top accessed examples",
	Description: `
Add top accessed examples.
`,
	Action: doAddTopAccessedExamples,
	Flags: []cli.Flag{
		cli.StringFlag{Name: "input-filename"},
	},
}
