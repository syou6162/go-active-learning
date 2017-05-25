package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
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
			return nil, errors.New(fmt.Sprintf("Invalid Label type %d in %s", label, line))
		}
	} else {
		return nil, errors.New(fmt.Sprintf("Invalid line: %s", line))
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

func FilterUnlabeledExamples(examples Examples) Examples {
	result := Examples{}

	alreadyLabeledByURL := make(map[string]bool)
	alreadyLabeledByTitle := make(map[string]bool)
	for _, e := range FilterLabeledExamples(examples) {
		alreadyLabeledByURL[e.Url] = true
		alreadyLabeledByTitle[e.Title] = true
	}

	for _, e := range examples {
		if _, ok := alreadyLabeledByURL[e.Url]; ok {
			continue
		}
		if _, ok := alreadyLabeledByTitle[e.Title]; ok {
			continue
		}
		if !e.IsLabeled() {
			alreadyLabeledByURL[e.Url] = true
			alreadyLabeledByTitle[e.Title] = true
			result = append(result, e)
		}
	}
	return result
}

func removeDuplicate(args []string) []string {
	results := make([]string, 0)
	encountered := map[string]bool{}
	for i := 0; i < len(args); i++ {
		if !encountered[args[i]] {
			encountered[args[i]] = true
			results = append(results, args[i])
		}
	}
	return results
}

func AttachMetaData(cache *Cache, examples Examples) {
	shuffle(examples)

	wg := &sync.WaitGroup{}
	cpus := runtime.NumCPU()
	runtime.GOMAXPROCS(cpus)
	sem := make(chan struct{}, 2000)
	for idx, e := range examples {
		wg.Add(1)
		sem <- struct{}{}
		go func(e *Example, idx int) {
			defer wg.Done()
			if example, ok := cache.Get(*e); ok {
				e.Title = example.Title
				e.Description = example.Description
				e.Body = example.Body
				e.StatusCode = example.StatusCode
			} else {
				article := GetArticle(e.Url)
				fmt.Fprintln(os.Stderr, "Fetching("+strconv.Itoa(idx)+"): "+e.Url)
				e.Title = article.Title
				e.Description = article.Description
				e.Body = article.Body
				e.StatusCode = article.StatusCode
				cache.Add(*e)
			}
			e.Fv = removeDuplicate(ExtractFeatures(*e))
			<-sem
		}(e, idx)
	}
	wg.Wait()
}

func FilterStatusCodeOkExamples(examples Examples) Examples {
	result := Examples{}

	for _, e := range examples {
		if e.StatusCode == 200 {
			result = append(result, e)
		}
	}

	return result
}
