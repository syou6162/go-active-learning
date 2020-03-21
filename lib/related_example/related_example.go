package related_example

import (
	"bufio"
	"fmt"
	"strconv"
	"strings"

	"os"

	"github.com/syou6162/go-active-learning/lib/model"
	"github.com/syou6162/go-active-learning/lib/service"
	"github.com/urfave/cli"
)

func parseLine(line string) (int, int, error) {
	tokens := strings.Split(line, "\t")
	if len(tokens) == 2 {
		exampleId, _ := strconv.ParseInt(tokens[0], 10, 0)
		relatedExampleId, _ := strconv.ParseInt(tokens[1], 10, 0)
		return int(exampleId), int(relatedExampleId), nil
	}
	return 0, 0, fmt.Errorf("Invalid line: %s", line)
}

func readRelatedExamples(filename string) ([]*model.RelatedExamples, error) {
	fp, err := os.Open(filename)
	defer fp.Close()
	if err != nil {
		return nil, err
	}

	exampleId2RelatedExampleIds := make(map[int][]int)
	scanner := bufio.NewScanner(fp)
	for scanner.Scan() {
		line := scanner.Text()
		exampleId, relatedExampleId, err := parseLine(line)
		if err != nil {
			return nil, err
		}
		if _, ok := exampleId2RelatedExampleIds[exampleId]; ok {
			exampleId2RelatedExampleIds[exampleId] = append(exampleId2RelatedExampleIds[exampleId], relatedExampleId)
		} else {
			exampleId2RelatedExampleIds[exampleId] = []int{relatedExampleId}
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	result := make([]*model.RelatedExamples, 0)
	for exampleId, relatedExampleIds := range exampleId2RelatedExampleIds {
		result = append(result, &model.RelatedExamples{ExampleId: exampleId, RelatedExampleIds: relatedExampleIds})
	}
	return result, nil
}

func doAddRelatedExamples(c *cli.Context) error {
	inputFilename := c.String("input-filename")

	if inputFilename == "" {
		_ = cli.ShowCommandHelp(c, "add-related-examples")
		return cli.NewExitError("`input-filename` is a required field.", 1)
	}

	app, err := service.NewDefaultApp()
	if err != nil {
		return err
	}
	defer app.Close()

	relatedExamplesList, err := readRelatedExamples(inputFilename)
	if err != nil {
		return err
	}
	for _, relatedExamples := range relatedExamplesList {
		for _, related := range relatedExamples.RelatedExampleIds {
			fmt.Print(relatedExamples.ExampleId)
			fmt.Print("\t")
			fmt.Println(related)
		}
		err := app.UpdateRelatedExamples(*relatedExamples)
		if err != nil {
			return err
		}
	}
	return nil
}

var CommandAddRelatedExamples = cli.Command{
	Name:  "add-related-examples",
	Usage: "add related examples",
	Description: `
Add related examples.
`,
	Action: doAddRelatedExamples,
	Flags: []cli.Flag{
		cli.StringFlag{Name: "input-filename"},
	},
}
