package main

import "testing"

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
