package util

import (
	"testing"

	"github.com/syou6162/go-active-learning/lib/example"
	"github.com/syou6162/go-active-learning/lib/model"
)

func TestFilterLabeledExamples(t *testing.T) {
	e1 := example.NewExample("https://b.hatena.ne.jp", model.POSITIVE)
	e2 := example.NewExample("https://www.yasuhisay.info", model.NEGATIVE)
	e3 := example.NewExample("http://google.com", model.UNLABELED)

	examples := FilterLabeledExamples(model.Examples{e1, e2, e3})
	if len(examples) != 2 {
		t.Error("Number of labeled examples should be 2")
	}
}

func TestFilterUnlabeledExamples(t *testing.T) {
	e1 := example.NewExample("https://b.hatena.ne.jp", model.POSITIVE)
	e2 := example.NewExample("https://www.yasuhisay.info", model.NEGATIVE)
	e3 := example.NewExample("http://google.com", model.UNLABELED)
	e3.Title = "Google"

	examples := FilterUnlabeledExamples(model.Examples{e1, e2, e3})
	if len(examples) != 1 {
		t.Error("Number of unlabeled examples should be 1")
	}
}

func TestFilterStatusCodeOkExamples(t *testing.T) {
	e1 := example.NewExample("https://b.hatena.ne.jp", model.POSITIVE)
	e1.StatusCode = 200
	e2 := example.NewExample("https://www.yasuhisay.info", model.NEGATIVE)
	e2.StatusCode = 404
	e3 := example.NewExample("http://google.com", model.UNLABELED)
	e3.StatusCode = 304

	examples := FilterStatusCodeOkExamples(model.Examples{e1, e2, e3})
	if len(examples) != 1 {
		t.Error("Number of examples (status code = 200) should be 1")
	}
}

func TestUniqueByFinalUrl(t *testing.T) {
	e1 := model.Example{FinalUrl: "aaa"}
	e2 := model.Example{FinalUrl: "bbb"}
	e3 := model.Example{FinalUrl: "aaa"}
	examples := model.Examples{&e1, &e2, &e3}
	result := UniqueByFinalUrl(examples)
	if len(result) != 2 {
		t.Errorf("length(result) should be %d, but %d", 2, len(result))
	}
}

func TestRemoveDuplicate(t *testing.T) {
	args := []string{"hoge", "fuga", "piyo", "hoge"}

	result := RemoveDuplicate(args)
	if len(result) != 3 {
		t.Error("Number of unique string in args should be 3")
	}
}
