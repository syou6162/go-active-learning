package classifier

import (
	"testing"

	"github.com/syou6162/go-active-learning/lib/example"
	"github.com/syou6162/go-active-learning/lib/model"
)

func TestPredictScore(t *testing.T) {
	e1 := example.NewExample("http://b.hatena.ne.jp", model.POSITIVE)
	e1.Title = "bookmark"
	e1.Fv = []string{"hoge", "fuga"}
	e2 := example.NewExample("http://google.com", model.NEGATIVE)
	e2.Title = "google"
	e2.Fv = []string{"piyo", "aaa"}
	e3 := example.NewExample("http://hatena.ne.jp", model.POSITIVE)
	e3.Title = "hatena"
	e3.Fv = []string{"hoge", "fuga"}
	e4 := example.NewExample("http://hogehoge.com", model.UNLABELED)
	e4.Title = "hogehoge"
	e4.Fv = []string{"piyo", "hoge"}

	examples := LearningInstances{e1, e2, e3, e4}
	c := NewMIRAClassifier(examples, 1.0)

	if c.PredictScore(e4.Fv) < 0.0 {
		t.Errorf("c.PredictScore(e4.Fv) == %f, want >= 0", c.PredictScore(e4.Fv))
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

	train, dev := splitTrainAndDev(LearningInstances{e1, e2, e3, e4, e5, e6, e7, e8, e9, e10})
	if len(train) != 8 {
		t.Error("Number of training examples should be 8")
	}
	if len(dev) != 2 {
		t.Error("Number of dev examples should be 2")
	}
}

func TestGetWeight(t *testing.T) {
	e1 := example.NewExample("http://b.hatena.ne.jp", model.POSITIVE)
	e1.Title = "bookmark"
	e1.Fv = []string{"hoge", "fuga"}
	e2 := example.NewExample("http://google.com", model.NEGATIVE)
	e2.Title = "google"
	e2.Fv = []string{"piyo", "aaa"}
	e3 := example.NewExample("http://hatena.ne.jp", model.POSITIVE)
	e3.Title = "hatena"
	e3.Fv = []string{"hoge", "fuga"}
	e4 := example.NewExample("http://hogehoge.com", model.UNLABELED)
	e4.Title = "hogehoge"
	e4.Fv = []string{"piyo", "hoge"}

	examples := LearningInstances{e1, e2, e3, e4}
	c := NewMIRAClassifier(examples, 1.0)

	if c.GetWeight("hoge") <= 0.0 {
		t.Errorf("c.GetWeight('hoge') == %f, want > 0", c.GetWeight("hoge"))
	}
}

func TestGetActiveFeatures(t *testing.T) {
	e1 := example.NewExample("http://b.hatena.ne.jp", model.POSITIVE)
	e1.Title = "bookmark"
	e1.Fv = []string{"hoge", "fuga"}
	e2 := example.NewExample("http://google.com", model.NEGATIVE)
	e2.Title = "google"
	e2.Fv = []string{"piyo", "aaa"}
	e3 := example.NewExample("http://hatena.ne.jp", model.POSITIVE)
	e3.Title = "hatena"
	e3.Fv = []string{"hoge", "fuga"}
	e4 := example.NewExample("http://hogehoge.com", model.UNLABELED)
	e4.Title = "hogehoge"
	e4.Fv = []string{"piyo", "hoge"}

	examples := LearningInstances{e1, e2, e3, e4}
	c := NewMIRAClassifier(examples, 1.0)

	if len(c.GetActiveFeatures()) <= 0 {
		t.Errorf("len(c.GetActiveFeatures()) <= %d, want > 0", len(c.GetActiveFeatures()))
	}
}
