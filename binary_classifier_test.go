package main

import (
	"testing"
)

func TestSortByScore(t *testing.T) {
	e0 := NewExample("http://www.yasuhisay.info/", POSITIVE)
	e0.Title = "yasuhisa"
	e0.Fv = []string{"hoge", "fuga"}
	e1 := NewExample("http://b.hatena.ne.jp", POSITIVE)
	e1.Title = "bookmark"
	e1.Fv = []string{"hoge", "fuga", "aaa"}
	e2 := NewExample("http://google.com", NEGATIVE)
	e2.Title = "google"
	e2.Fv = []string{"piyo", "aaa"}
	e3 := NewExample("http://hatena.ne.jp", UNLABELED)
	e3.Title = "hatena"
	e3.Fv = []string{"hoge", "fuga"}
	e4 := NewExample("http://hogehoge.com", UNLABELED)
	e4.Title = "hogehoge"
	e4.Fv = []string{"piyo", "hoge"}

	examples := Examples{e0, e1, e2, e3, e4}
	c := NewBinaryClassifier(examples)
	examples = SortByScore(c, examples)

	if len(examples) != 2 {
		t.Errorf("len(example) == %d, want 2", len(examples))
	}

	if examples[0].Title != "hogehoge" {
		t.Errorf("example[0].Title == %s, want 'hogehoge'", examples[0].Title)
	}
}
