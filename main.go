package main

import (
	"fmt"
	"os"

	"github.com/mattn/go-tty"
)

func ask(e *Example) (bool, error) {
	fmt.Println("Label this example: " + e.url)
	t, err := tty.Open()
	defer t.Close()
	if err != nil {
		return false, err
	}
	var r rune
	for r == 0 {
		r, err = t.ReadRune()
		if err != nil {
			return false, err
		}
	}
	switch r {
	case 'p':
		e.Annotate(POSITIVE)
		fmt.Println("Labeled as positive")
		return true, nil
	case 'n':
		e.Annotate(NEGATIVE)
		fmt.Println("Labeled as negative")
		return true, nil
	default:
		return false, nil
	}
}

func main() {
	examples, _ := ReadExamples(os.Args[1])

	for {
		e := RandomSelectOneExample(examples)
		if e == nil {
			break
		}
		ask(e)
	}

	for _, e := range examples {
		fmt.Println(e)
	}
}
