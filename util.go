package main

import (
	"bufio"
	"errors"
	"fmt"
	"math"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/syou6162/go-active-learning/lib/cache"
	"github.com/syou6162/go-active-learning/lib/example"
)

func ParseLine(line string) (*example.Example, error) {
	tokens := strings.Split(line, "\t")
	var url string
	if len(tokens) == 1 {
		url = tokens[0]
		return example.NewExample(url, example.UNLABELED), nil
	} else if len(tokens) == 2 {
		url = tokens[0]
		label, _ := strconv.ParseInt(tokens[1], 10, 0)
		switch example.LabelType(label) {
		case example.POSITIVE, example.NEGATIVE, example.UNLABELED:
			return example.NewExample(url, example.LabelType(label)), nil
		default:
			return nil, errors.New(fmt.Sprintf("Invalid Label type %d in %s", label, line))
		}
	} else {
		return nil, errors.New(fmt.Sprintf("Invalid line: %s", line))
	}
}

func ReadExamples(filename string) ([]*example.Example, error) {
	fp, err := os.Open(filename)
	defer fp.Close()
	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(fp)
	var examples example.Examples
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

func WriteExamples(examples example.Examples, filename string) error {
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

func FilterLabeledExamples(examples example.Examples) example.Examples {
	var result example.Examples
	for _, e := range examples {
		if e.IsLabeled() {
			result = append(result, e)
		}
	}
	return result
}

func FilterUnlabeledExamples(examples example.Examples) example.Examples {
	result := example.Examples{}

	alreadyLabeledByURL := make(map[string]bool)
	alreadyLabeledByTitle := make(map[string]bool)
	for _, e := range FilterLabeledExamples(examples) {
		alreadyLabeledByURL[e.Url] = true
		alreadyLabeledByURL[e.FinalUrl] = true
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
			alreadyLabeledByURL[e.FinalUrl] = true
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

func attachMetaData(cache *cache.Cache, examples example.Examples) {
	oldStdout := os.Stdout
	readFile, writeFile, _ := os.Pipe()
	os.Stdout = writeFile

	defer func() {
		writeFile.Close()
		readFile.Close()
		os.Stdout = oldStdout
	}()

	shuffle(examples)

	wg := &sync.WaitGroup{}
	cpus := runtime.NumCPU()
	runtime.GOMAXPROCS(cpus)
	sem := make(chan struct{}, 4)
	for idx, e := range examples {
		wg.Add(1)
		sem <- struct{}{}
		go func(e *example.Example, idx int) {
			defer wg.Done()
			if tmp, ok := cache.Get(*e); ok {
				e.Title = tmp.Title
				e.FinalUrl = tmp.FinalUrl
				e.Description = tmp.Description
				e.Body = tmp.Body
				e.StatusCode = tmp.StatusCode
				e.Fv = tmp.Fv
			} else {
				fmt.Fprintln(os.Stderr, "Fetching("+strconv.Itoa(idx)+"): "+e.Url)
				article := GetArticle(e.Url)
				e.Title = article.Title
				e.FinalUrl = article.Url
				e.Description = article.Description
				e.Body = article.Body
				e.StatusCode = article.StatusCode
				e.Fv = removeDuplicate(example.ExtractFeatures(*e))
				e.Description = ""
				e.Body = ""
				cache.Add(*e)
			}
			<-sem
		}(e, idx)
	}
	wg.Wait()
}

func AttachMetaData(cache *cache.Cache, examples example.Examples) {
	batchSize := 100
	examplesList := make([]example.Examples, 0)
	n := len(examples)

	for i := 0; i < n; i += batchSize {
		max := int(math.Min(float64(i+batchSize), float64(n)))
		examplesList = append(examplesList, examples[i:max])
	}
	for _, l := range examplesList {
		attachMetaData(cache, l)
	}
}

func FilterStatusCodeOkExamples(examples example.Examples) example.Examples {
	result := example.Examples{}

	for _, e := range examples {
		if e.StatusCode == 200 {
			result = append(result, e)
		}
	}

	return result
}

func NewOutputFilename() string {
	t := time.Now()
	return fmt.Sprintf("%d-%02d-%02d-%02d-%02d.txt", t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute())
}

func splitTrainAndDev(examples example.Examples) (train example.Examples, dev example.Examples) {
	shuffle(examples)
	n := int(0.8 * float64(len(examples)))
	return examples[0:n], examples[n:]
}
