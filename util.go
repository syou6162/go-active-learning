package main

import (
	"bufio"
	"errors"
	"os"
	"strconv"
	"strings"
)

func ParseLine(line string) (*Example, error) {
	tokens := strings.Split(line, "\t")
	var url string
	if len(tokens) == 1 {
		url = tokens[0]
		return NewExample(url, UNLABELED), nil
	} else if len(tokens) == 2 {
		url = tokens[0]
		label, _ := strconv.ParseInt(tokens[1], 10, 0)
		switch LabelType(label) {
		case POSITIVE, NEGATIVE, UNLABELED:
			return NewExample(url, LabelType(label)), nil
		default:
			return nil, errors.New("Invalid Label type")
		}
	} else {
		return nil, errors.New("Invalid line")
	}
}

func ReadExamples(filename string) ([]*Example, error) {
	fp, err := os.Open(filename)
	defer fp.Close()
	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(fp)
	var examples Examples
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

func WriteExamples(examples Examples, filename string) error {
	fp, err := os.Create(filename)
	defer fp.Close()
	if err != nil {
		return err
	}

	writer := bufio.NewWriter(fp)
	for _, e := range examples {
		if e.IsNew && e.IsLabeled() {
			_, err := writer.WriteString(e.Url + "\t" + strconv.Itoa(int(e.Label)) + "\n")
			if err != nil {
				return err
			}
		}
	}

	writer.Flush()
	return nil
}

func FilterLabeledExamples(examples Examples) Examples {
	var result Examples
	for _, e := range examples {
		if e.IsLabeled() {
			result = append(result, e)
		}
	}
	return result
}
