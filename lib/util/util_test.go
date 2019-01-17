package util

import (
	"testing"

	"github.com/syou6162/go-active-learning/lib/example"
	"github.com/syou6162/go-active-learning/lib/model"
)

func TestFilterLabeledExamples(t *testing.T) {
	e1 := example.NewExample("http://b.hatena.ne.jp", model.POSITIVE)
	e2 := example.NewExample("https://www.yasuhisay.info", model.NEGATIVE)
	e3 := example.NewExample("http://google.com", model.UNLABELED)

	examples := FilterLabeledExamples(model.Examples{e1, e2, e3})
	if len(examples) != 2 {
		t.Error("Number of labeled examples should be 2")
	}
}

func TestFilterUnlabeledExamples(t *testing.T) {
	e1 := example.NewExample("http://b.hatena.ne.jp", model.POSITIVE)
	e2 := example.NewExample("https://www.yasuhisay.info", model.NEGATIVE)
	e3 := example.NewExample("http://google.com", model.UNLABELED)
	e3.Title = "Google"

	examples := FilterUnlabeledExamples(model.Examples{e1, e2, e3})
	if len(examples) != 1 {
		t.Error("Number of unlabeled examples should be 1")
	}
}

func TestFilterStatusCodeOkExamples(t *testing.T) {
	e1 := example.NewExample("http://b.hatena.ne.jp", model.POSITIVE)
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

func TestSplitTrainAndDev(t *testing.T) {
	e1 := example.NewExample("http://a.hatena.ne.jp", model.POSITIVE)
	e2 := example.NewExample("https://www.yasuhisay.info", model.NEGATIVE)
	e3 := example.NewExample("http://google.com", model.UNLABELED)
	e4 := example.NewExample("http://a.hatena.ne.jp", model.POSITIVE)
	e5 := example.NewExample("https://www.yasuhisay.info", model.NEGATIVE)
	e6 := example.NewExample("http://a.hatena.ne.jp", model.POSITIVE)
	e7 := example.NewExample("https://www.yasuhisay.info", model.NEGATIVE)
	e8 := example.NewExample("http://google.com", model.UNLABELED)
	e9 := example.NewExample("http://a.hatena.ne.jp", model.POSITIVE)
	e10 := example.NewExample("https://www.yasuhisay.info", model.NEGATIVE)

	train, dev := SplitTrainAndDev(model.Examples{e1, e2, e3, e4, e5, e6, e7, e8, e9, e10})
	if len(train) != 8 {
		t.Error("Number of training examples should be 8")
	}
	if len(dev) != 2 {
		t.Error("Number of dev examples should be 2")
	}
}

func TestSortByCommentUsefulness(t *testing.T) {
	e1 := example.NewExample("https://ai.googleblog.com/2019/01/looking-back-at-googles-research.html", model.POSITIVE)
	e1.Title = "Google AI Blog: Looking Back at Google’s Research Efforts in 2018"

	tweets := model.ReferringTweets{
		&model.Tweet{ScreenName: "aaa", FullText: "Google AI Blog: Looking Back at Google’s Research Efforts in 2018 https://t.co/YhzgxeTAft"},
		&model.Tweet{ScreenName: "bbb", FullText: "Googleすごい https://t.co/YhzgxeTAft"},
	}

	result := SortByCommentUsefulness(*e1, tweets)
	if len(result) != 2 {
		t.Error("Number of dev examples should be 2")
	}
	if result[0].ScreenName != "bbb" {
		t.Error("result must not be sorted")
	}
}
