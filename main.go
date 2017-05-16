package main

import (
	"fmt"
	"os"
	"sync"
	"github.com/mattn/go-tty"
)

type ActionType int

const (
	LABEL_AS_POSITIVE ActionType = iota
	LABEL_AS_NEGATIVE
	SAVE
	HELP
	SKIP
	EXIT
)

func input2ActionType() (ActionType, error) {
	t, err := tty.Open()
	defer t.Close()
	if err != nil {
		return EXIT, err
	}
	var r rune
	for r == 0 {
		r, err = t.ReadRune()
		if err != nil {
			return SKIP, err
		}
	}
	switch r {
	case 'p':
		return LABEL_AS_POSITIVE, nil
	case 'n':
		return LABEL_AS_NEGATIVE, nil
	case 's':
		return SAVE, nil
	case 'h':
		return HELP, nil
	case 'e':
		return EXIT, nil
	default:
		return SKIP, nil
	}
}

func main() {
	cacheFilename := "cache.bin"

	cache, _ := LoadCache(cacheFilename)
	examples, _ := ReadExamples(os.Args[1])

	outputFilename := os.Args[2]
	shuffle(examples)

	wg := &sync.WaitGroup{}
	cpus := 20
	semaphore := make(chan int, cpus)
	for _, e := range examples {
		wg.Add(1)
		go func(example *Example) {
			defer wg.Done()
			semaphore <- 1
			if e, ok := cache.Cache[example.Url]; ok {
				example.Title = e.Title
				example.Description = e.Description
				example.Body = e.Body
			} else {
				article := GetArticle(example.Url)
				fmt.Println("Fetching: " + example.Url)
				example.Title = article.Title
				example.Description = article.Description
				example.Body = article.Body
				cache.Add(*example)
			}
			example.Fv = ExtractFeatures(*example)
			<-semaphore
		}(e)
	}
	wg.Wait()

	model := TrainedModel(examples)

annotationLoop:
	for {
		unlabeledExamples := model.SortByScore(examples)
		if len(unlabeledExamples) == 0 {
			break
		}
		e := unlabeledExamples[0]
		if e == nil {
			break
		}
		fmt.Println("Label this example (Score: " + fmt.Sprintf("%0.03f", e.Score) + "): " + e.Url + " (" + e.Title + ")")
		cache.Add(*e)

		act, err := input2ActionType()
		if err != nil {
			return
		}
		switch act {
		case LABEL_AS_POSITIVE:
			fmt.Println("Labeled as positive")
			e.Annotate(POSITIVE)
		case LABEL_AS_NEGATIVE:
			fmt.Println("Labeled as negative")
			e.Annotate(NEGATIVE)
		case SKIP:
			fmt.Println("Skiped this example")
			continue
		case SAVE:
			fmt.Println("Saved labeld examples")
			WriteExamples(examples, outputFilename)
		case HELP:
			fmt.Println("ToDo: SHOW HELP")
		case EXIT:
			fmt.Println("EXIT")
			break annotationLoop
		default:
			break annotationLoop
		}
		model = TrainedModel(examples)
	}

	WriteExamples(examples, outputFilename)
	cache.Save(cacheFilename)
}
