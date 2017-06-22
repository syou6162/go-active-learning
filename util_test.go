package main

import (
	"fmt"
	"testing"
)

func TestParseLine(t *testing.T) {
	line1 := "http://example.com\t1"
	e, err := ParseLine(line1)

	if err != nil {
		t.Error("cannot parse line1")
	}
	if e.Label != POSITIVE {
		t.Error("Label must be POSITIVE")
	}

	line2 := "http://example.com\t-1"
	e, err = ParseLine(line2)

	if err != nil {
		t.Error("cannot parse line2")
	}
	if e.Label != NEGATIVE {
		t.Error("Label must be NEGATIVE")
	}

	line3 := "http://example.com"
	e, err = ParseLine(line3)

	if err != nil {
		t.Error("cannot parse line3")
	}
	if e.Label != UNLABELED {
		t.Error("Label must be UNLABELED")
	}

	line4 := "http://example.com\t2"
	e, err = ParseLine(line4)

	if e != nil {
		t.Error("wrong line format")
	}
}

func TestReadExamples(t *testing.T) {
	filename := "tech_input_example.txt"
	examples, err := ReadExamples(filename)

	if err != nil {
		t.Error(fmt.Printf("Cannot read examples from %s", filename))
	}
	if len(examples) == 0 {
		t.Error(fmt.Printf("%s should contain more than one examples", filename))
	}
}

func TestWriteExamples(t *testing.T) {
	filename := ".write_test.txt"
	e1 := NewExample("http://b.hatena.ne.jp", POSITIVE)
	e2 := NewExample("http://www.yasuhisay.info", NEGATIVE)

	err := WriteExamples(Examples{e1, e2}, filename)
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

func TestFilterLabeledExamples(t *testing.T) {
	e1 := NewExample("http://b.hatena.ne.jp", POSITIVE)
	e2 := NewExample("http://www.yasuhisay.info", NEGATIVE)
	e3 := NewExample("http://google.com", UNLABELED)

	examples := FilterLabeledExamples(Examples{e1, e2, e3})
	if len(examples) != 2 {
		t.Error("Number of labeled examples should be 2")
	}
}

func TestFilterUnlabeledExamples(t *testing.T) {
	e1 := NewExample("http://b.hatena.ne.jp", POSITIVE)
	e2 := NewExample("http://www.yasuhisay.info", NEGATIVE)
	e3 := NewExample("http://google.com", UNLABELED)
	e3.Title = "Google"

	examples := FilterUnlabeledExamples(Examples{e1, e2, e3})
	if len(examples) != 1 {
		t.Error("Number of unlabeled examples should be 1")
	}
}

func TestFilterStatusCodeOkExamples(t *testing.T) {
	e1 := NewExample("http://b.hatena.ne.jp", POSITIVE)
	e1.StatusCode = 200
	e2 := NewExample("http://www.yasuhisay.info", NEGATIVE)
	e2.StatusCode = 404
	e3 := NewExample("http://google.com", UNLABELED)
	e3.StatusCode = 304

	examples := FilterStatusCodeOkExamples(Examples{e1, e2, e3})
	if len(examples) != 1 {
		t.Error("Number of examples (status code = 200) should be 1")
	}
}

func TestRemoveDuplicate(t *testing.T) {
	args := []string{"hoge", "fuga", "piyo", "hoge"}

	result := removeDuplicate(args)
	if len(result) != 3 {
		t.Error("Number of unique string in args should be 3")
	}
}

func TestSplitTrainAndDev(t *testing.T) {
	e1 := NewExample("http://b.hatena.ne.jp", POSITIVE)
	e2 := NewExample("http://www.yasuhisay.info", NEGATIVE)
	e3 := NewExample("http://google.com", UNLABELED)
	e4 := NewExample("http://b.hatena.ne.jp", POSITIVE)
	e5 := NewExample("http://www.yasuhisay.info", NEGATIVE)
	e6 := NewExample("http://b.hatena.ne.jp", POSITIVE)
	e7 := NewExample("http://www.yasuhisay.info", NEGATIVE)
	e8 := NewExample("http://google.com", UNLABELED)
	e9 := NewExample("http://b.hatena.ne.jp", POSITIVE)
	e10 := NewExample("http://www.yasuhisay.info", NEGATIVE)

	train, dev := splitTrainAndDev(Examples{e1, e2, e3, e4, e5, e6, e7, e8, e9, e10})
	if len(train) != 8 {
		t.Error("Number of training examples should be 8")
	}
	if len(dev) != 2 {
		t.Error("Number of dev examples should be 2")
	}
}

func TestAttachMetaData(t *testing.T) {
	e1 := NewExample("http://b.hatena.ne.jp", POSITIVE)
	e2 := NewExample("http://www.yasuhisay.info", NEGATIVE)
	e3 := NewExample("http://google.com", UNLABELED)
	examples := Examples{e1, e2, e3}
	AttachMetaData(NewCache(), examples)

	if examples[0].Title == "" {
		t.Errorf("Title must not be empty for %s", examples[0].Url)
	}
	if len(examples[0].Fv) == 0 {
		t.Errorf("Feature vector must not be empty for %s", examples[0].Url)
	}
}
