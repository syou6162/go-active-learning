package file

import (
	"fmt"
	"testing"

	"github.com/syou6162/go-active-learning/lib/example"
	"github.com/syou6162/go-active-learning/lib/model"
)

func TestParseLine(t *testing.T) {
	line1 := "http://model.com\t1"
	e, err := ParseLine(line1)

	if err != nil {
		t.Error("cannot parse line1")
	}
	if e.Label != model.POSITIVE {
		t.Error("Label must be POSITIVE")
	}

	line2 := "http://model.com\t-1"
	e, err = ParseLine(line2)

	if err != nil {
		t.Error("cannot parse line2")
	}
	if e.Label != model.NEGATIVE {
		t.Error("Label must be NEGATIVE")
	}

	line3 := "http://model.com"
	e, err = ParseLine(line3)

	if err != nil {
		t.Error("cannot parse line3")
	}
	if e.Label != model.UNLABELED {
		t.Error("Label must be UNLABELED")
	}

	line4 := "http://model.com\t2"
	e, err = ParseLine(line4)

	if e != nil {
		t.Error("wrong line format")
	}
}

func TestReadExamples(t *testing.T) {
	filename := "../../../tech_input_example.txt"
	examples, err := ReadExamples(filename)

	if err != nil {
		fmt.Println(err.Error())
		t.Error(fmt.Printf("Cannot read examples from %s\n", filename))
	}
	if len(examples) == 0 {
		t.Error(fmt.Printf("%s should contain more than one examples\n", filename))
	}
}

func TestWriteExamples(t *testing.T) {
	filename := ".write_test.txt"
	e1 := example.NewExample("http://b.hatena.ne.jp", model.POSITIVE)
	e2 := example.NewExample("http://www.yasuhisay.info", model.NEGATIVE)

	err := WriteExamples(model.Examples{e1, e2}, filename)
	if err != nil {
		t.Error(fmt.Printf("Cannot write examples to %s", filename))
	}

	examples, err := ReadExamples(filename)
	if err != nil {
		t.Error(fmt.Printf("Cannot read examples from %s", filename))
	}
	if len(examples) == 2 {
		t.Error(fmt.Printf("%s should contain two examples", filename))
	}
}
