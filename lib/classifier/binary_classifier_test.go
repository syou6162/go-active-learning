package classifier

import (
	"testing"

	"github.com/syou6162/go-active-learning/lib/example"
	"github.com/syou6162/go-active-learning/lib/model"
)

func TestSortByScore(t *testing.T) {
	e0 := example.NewExample("https://www.yasuhisay.info/", model.POSITIVE)
	e0.Title = "yasuhisa"
	e0.Fv = []string{"hoge", "fuga"}
	e1 := example.NewExample("http://b.hatena.ne.jp", model.POSITIVE)
	e1.Title = "bookmark"
	e1.Fv = []string{"hoge", "fuga", "aaa"}
	e2 := example.NewExample("http://google.com", model.NEGATIVE)
	e2.Title = "google"
	e2.Fv = []string{"piyo", "aaa"}
	e3 := example.NewExample("http://hatena.ne.jp", model.UNLABELED)
	e3.Title = "hatena"
	e3.Fv = []string{"hoge", "fuga"}
	e4 := example.NewExample("http://hogehoge.com", model.UNLABELED)
	e4.Title = "hogehoge"
	e4.Fv = []string{"piyo", "hoge"}

	examples := model.Examples{e0, e1, e2, e3, e4}
	c := NewBinaryClassifier(examples)
	examples = SortByScore(c, examples)

	if len(examples) != 2 {
		t.Errorf("len(example) == %d, want 2", len(examples))
	}

	if examples[0].Title != "hogehoge" {
		t.Errorf("example[0].Title == %s, want 'hogehoge'", examples[0].Title)
	}
}
