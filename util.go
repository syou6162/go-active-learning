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
			return nil, errors.New("Invalid label type")
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
