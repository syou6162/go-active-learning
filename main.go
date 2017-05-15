package main

import (
	"fmt"
	"os"

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
	for _, e := range examples {
		if title, ok := cache.cache[e.url]; ok {
			e.title = title
		} else {
			title = GetTitle(e.url)
			fmt.Println("Fetching: " + e.url)
			e.title = title
			cache.Add(*e)
		}
		e.fv = ExtractFeatures(e.title)
	}
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
		fmt.Println("Label this example (score: " + fmt.Sprintf("%0.03f", e.score) + "): " + e.url + " (" + e.title + ")")
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

	WriteExamples(examples, os.Args[2])
	cache.Save(cacheFilename)
}
