package main

import (
	"testing"
)

func TestPredictScore(t *testing.T) {
	e1 := NewExample("http://b.hatena.ne.jp", POSITIVE)
	e1.Title = "bookmark"
	e1.Fv = []string{"hoge", "fuga"}
	e2 := NewExample("http://google.com", NEGATIVE)
	e2.Title = "google"
	e2.Fv = []string{"piyo", "aaa"}
	e3 := NewExample("http://hatena.ne.jp", POSITIVE)
	e3.Title = "hatena"
	e3.Fv = []string{"hoge", "fuga"}
	e4 := NewExample("http://hogehoge.com", UNLABELED)
	e4.Title = "hogehoge"
	e4.Fv = []string{"piyo", "hoge"}

	examples := Examples{e1, e2, e3, e4}
	c := NewBinaryClassifier(examples)

	if c.PredictScore(e4.Fv) <= 0.0 {
		t.Errorf("c.PredictScore(e4.Fv) == %f, want > 0", c.PredictScore(e4.Fv))
	}
}

func TestGetWeight(t *testing.T) {
	e1 := NewExample("http://b.hatena.ne.jp", POSITIVE)
	e1.Title = "bookmark"
	e1.Fv = []string{"hoge", "fuga"}
	e2 := NewExample("http://google.com", NEGATIVE)
	e2.Title = "google"
	e2.Fv = []string{"piyo", "aaa"}
	e3 := NewExample("http://hatena.ne.jp", POSITIVE)
	e3.Title = "hatena"
	e3.Fv = []string{"hoge", "fuga"}
	e4 := NewExample("http://hogehoge.com", UNLABELED)
	e4.Title = "hogehoge"
	e4.Fv = []string{"piyo", "hoge"}

	examples := Examples{e1, e2, e3, e4}
	c := NewBinaryClassifier(examples)

	if c.GetWeight("hoge") <= 0.0 {
		t.Errorf("c.GetWeight('hoge') == %f, want > 0", c.GetWeight("hoge"))
	}
}

func TestGetActiveFeatures(t *testing.T) {
	e1 := NewExample("http://b.hatena.ne.jp", POSITIVE)
	e1.Title = "bookmark"
	e1.Fv = []string{"hoge", "fuga"}
	e2 := NewExample("http://google.com", NEGATIVE)
	e2.Title = "google"
	e2.Fv = []string{"piyo", "aaa"}
	e3 := NewExample("http://hatena.ne.jp", POSITIVE)
	e3.Title = "hatena"
	e3.Fv = []string{"hoge", "fuga"}
	e4 := NewExample("http://hogehoge.com", UNLABELED)
	e4.Title = "hogehoge"
	e4.Fv = []string{"piyo", "hoge"}

	examples := Examples{e1, e2, e3, e4}
	c := NewBinaryClassifier(examples)

	if len(c.GetActiveFeatures()) <= 0 {
		t.Errorf("len(c.GetActiveFeatures()) <= %f, want > 0", len(c.GetActiveFeatures()))
	}
}
