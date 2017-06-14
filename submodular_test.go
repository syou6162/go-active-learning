package main

import (
	"fmt"
	"testing"
)

func TestGetDF(t *testing.T) {
	e := NewExample("http://b.hatena.ne.jp", POSITIVE)
	e.Body = "こんにちは、日本"
	dfMap := GetDF(*e)

	japan := "BODY:日本"
	if _, ok := dfMap[japan]; !ok {
		t.Error(fmt.Printf("Example must contain %s", japan))
	}
}

func TestGetIDF(t *testing.T) {
	e1 := NewExample("http://b.hatena.ne.jp", POSITIVE)
	e1.Body = "こんにちは、日本"
	idfMap := GetIDF(Examples{e1})

	japan := "BODY:日本"
	if _, ok := idfMap[japan]; !ok {
		t.Error(fmt.Printf("Example must contain %s", japan))
	}
}
