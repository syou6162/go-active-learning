package main

import (
	"github.com/syou6162/go-active-learning/lib/example"
	"testing"
)

func TestPredictScore(t *testing.T) {
	e1 := example.NewExample("http://b.hatena.ne.jp", example.POSITIVE)
	e1.Title = "bookmark"
	e1.Fv = []string{"hoge", "fuga"}
	e2 := example.NewExample("http://google.com", example.NEGATIVE)
	e2.Title = "google"
	e2.Fv = []string{"piyo", "aaa"}
	e3 := example.NewExample("http://hatena.ne.jp", example.POSITIVE)
	e3.Title = "hatena"
	e3.Fv = []string{"hoge", "fuga"}
	e4 := example.NewExample("http://hogehoge.com", example.UNLABELED)
	e4.Title = "hogehoge"
	e4.Fv = []string{"piyo", "hoge"}

	examples := example.Examples{e1, e2, e3, e4}
	c := NewBinaryClassifier(examples)

	if c.PredictScore(e4.Fv) < 0.0 {
		t.Errorf("c.PredictScore(e4.Fv) == %f, want >= 0", c.PredictScore(e4.Fv))
	}
}

func TestGetWeight(t *testing.T) {
	e1 := example.NewExample("http://b.hatena.ne.jp", example.POSITIVE)
	e1.Title = "bookmark"
	e1.Fv = []string{"hoge", "fuga"}
	e2 := example.NewExample("http://google.com", example.NEGATIVE)
	e2.Title = "google"
	e2.Fv = []string{"piyo", "aaa"}
	e3 := example.NewExample("http://hatena.ne.jp", example.POSITIVE)
	e3.Title = "hatena"
	e3.Fv = []string{"hoge", "fuga"}
	e4 := example.NewExample("http://hogehoge.com", example.UNLABELED)
	e4.Title = "hogehoge"
	e4.Fv = []string{"piyo", "hoge"}

	examples := example.Examples{e1, e2, e3, e4}
	c := NewBinaryClassifier(examples)

	if c.GetWeight("hoge") <= 0.0 {
		t.Errorf("c.GetWeight('hoge') == %f, want > 0", c.GetWeight("hoge"))
	}
}

func TestGetActiveFeatures(t *testing.T) {
	e1 := example.NewExample("http://b.hatena.ne.jp", example.POSITIVE)
	e1.Title = "bookmark"
	e1.Fv = []string{"hoge", "fuga"}
	e2 := example.NewExample("http://google.com", example.NEGATIVE)
	e2.Title = "google"
	e2.Fv = []string{"piyo", "aaa"}
	e3 := example.NewExample("http://hatena.ne.jp", example.POSITIVE)
	e3.Title = "hatena"
	e3.Fv = []string{"hoge", "fuga"}
	e4 := example.NewExample("http://hogehoge.com", example.UNLABELED)
	e4.Title = "hogehoge"
	e4.Fv = []string{"piyo", "hoge"}

	examples := example.Examples{e1, e2, e3, e4}
	c := NewBinaryClassifier(examples)

	if len(c.GetActiveFeatures()) <= 0 {
		t.Errorf("len(c.GetActiveFeatures()) <= %d, want > 0", len(c.GetActiveFeatures()))
	}
}
