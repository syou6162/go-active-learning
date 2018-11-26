package file

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/syou6162/go-active-learning/lib/example"
	"github.com/syou6162/go-active-learning/lib/model"
)

func ParseLine(line string) (*model.Example, error) {
	tokens := strings.Split(line, "\t")
	var url string
	if len(tokens) == 1 {
		url = tokens[0]
		return example.NewExample(url, model.UNLABELED), nil
	} else if len(tokens) == 2 {
		url = tokens[0]
		label, _ := strconv.ParseInt(tokens[1], 10, 0)
		switch model.LabelType(label) {
		case model.POSITIVE, model.NEGATIVE, model.UNLABELED:
			return example.NewExample(url, model.LabelType(label)), nil
		default:
			return nil, errors.New(fmt.Sprintf("Invalid Label type %d in %s", label, line))
		}
	} else {
		return nil, errors.New(fmt.Sprintf("Invalid line: %s", line))
	}
}

func ReadExamples(filename string) ([]*model.Example, error) {
	fp, err := os.Open(filename)
	defer fp.Close()
	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(fp)
	var examples model.Examples
	for scanner.Scan() {
		line := scanner.Text()
		e, err := ParseLine(line)
		if err != nil {
			return nil, err
		}
		examples = append(examples, e)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return examples, nil
}

func WriteExamples(examples model.Examples, filename string) error {
	fp, err := os.Create(filename)
	defer fp.Close()
	if err != nil {
		return err
	}

	writer := bufio.NewWriter(fp)
	for _, e := range examples {
		if e.IsNew && e.IsLabeled() {
			url := e.FinalUrl
			if url == "" {
				url = e.Url
			}
			_, err := writer.WriteString(url + "\t" + strconv.Itoa(int(e.Label)) + "\n")
			if err != nil {
				return err
			}
		}
	}

	writer.Flush()
	return nil
}
