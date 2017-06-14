package main

import (
	"testing"
)

func TestGetArticle(t *testing.T) {
	a := GetArticle("http://b.hatena.ne.jp")

	if a.Title == "" {
		t.Error("Title must not be empty")
	}
	if a.Body == "" {
		t.Error("Body must not be empty")
	}
	if a.StatusCode != 200 {
		t.Error("StatusCode must be 200")
	}
}
