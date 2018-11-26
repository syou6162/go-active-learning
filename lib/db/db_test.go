package db_test

import (
	"log"
	"os"
	"testing"
	"time"

	"github.com/syou6162/go-active-learning/lib/db"
	"github.com/syou6162/go-active-learning/lib/example"
	"github.com/syou6162/go-active-learning/lib/model"
)

func TestMain(m *testing.M) {
	err := db.Init()
	if err != nil {
		log.Fatal(err.Error())
	}
	defer db.Close()

	ret := m.Run()
	os.Exit(ret)
}

func TestPing(t *testing.T) {
	if err := db.Ping(); err != nil {
		t.Errorf(err.Error())
	}
}

func TestInsertExamplesFromReader(t *testing.T) {
	_, err := db.DeleteAllExamples()
	if err != nil {
		t.Error(err)
	}

	fp, err := os.Open("../../tech_input_example.txt")
	defer fp.Close()
	if err != nil {
		t.Error(err)
	}
	db.InsertExamplesFromReader(fp)

	examples, err := db.ReadExamples()
	if err != nil {
		t.Error(err)
	}
	if len(examples) == 0 {
		t.Errorf("len(examples) > 0, but %d", len(examples))
	}
}

func TestInsertOrUpdateExample(t *testing.T) {
	_, err := db.DeleteAllExamples()
	if err != nil {
		t.Error(err)
	}

	_, err = db.InsertOrUpdateExample(example.NewExample("http://hoge.com", model.UNLABELED))
	if err != nil {
		t.Error(err)
	}

	examples, err := db.ReadExamples()
	if err != nil {
		t.Error(err)
	}
	if len(examples) != 1 {
		t.Errorf("len(examples) == %d, want 1", len(examples))
	}
	if examples[0].Label != model.UNLABELED {
		t.Errorf("label == %d, want 1", examples[0].Label)
	}

	// same url
	_, err = db.InsertOrUpdateExample(example.NewExample("http://hoge.com", model.NEGATIVE))
	if err != nil {
		t.Error(err)
	}

	examples, err = db.ReadExamples()
	if err != nil {
		t.Error(err)
	}
	if len(examples) != 1 {
		t.Errorf("len(examples) == %d, want 1", len(examples))
	}
	if examples[0].Label != model.NEGATIVE {
		t.Errorf("label == %d, want 1", examples[0].Label)
	}

	// same url but different label
	_, err = db.InsertOrUpdateExample(example.NewExample("http://hoge.com", model.POSITIVE))
	if err != nil {
		t.Error(err)
	}

	examples, err = db.ReadExamples()
	if err != nil {
		t.Error(err)
	}
	if len(examples) != 1 {
		t.Errorf("len(examples) == %d, want 1", len(examples))
	}
	if examples[0].Label != model.POSITIVE {
		t.Errorf("label == %d, want 1", examples[0].Label)
	}

	// cannot update to unlabeled
	_, err = db.InsertOrUpdateExample(example.NewExample("http://hoge.com", model.UNLABELED))
	if err != nil {
		t.Error(err)
	}

	examples, err = db.ReadExamples()
	if err != nil {
		t.Error(err)
	}
	if len(examples) != 1 {
		t.Errorf("len(examples) == %d, want 1", len(examples))
	}
	if examples[0].Label != model.POSITIVE {
		t.Errorf("label == %d, want 1", examples[0].Label)
	}

	// different url
	_, err = db.InsertOrUpdateExample(example.NewExample("http://another.com", model.NEGATIVE))
	if err != nil {
		t.Error(err)
	}

	examples, err = db.ReadExamples()
	if err != nil {
		t.Error(err)
	}
	if len(examples) != 2 {
		t.Errorf("len(examples) == %d, want 2", len(examples))
	}
}

func TestReadLabeledExamples(t *testing.T) {
	_, err := db.DeleteAllExamples()
	if err != nil {
		t.Error(err)
	}

	_, err = db.InsertOrUpdateExample(example.NewExample("http://hoge1.com", model.POSITIVE))
	if err != nil {
		t.Error(err)
	}
	_, err = db.InsertOrUpdateExample(example.NewExample("http://hoge2.com", model.NEGATIVE))
	if err != nil {
		t.Error(err)
	}
	_, err = db.InsertOrUpdateExample(example.NewExample("http://hoge3.com", model.UNLABELED))
	if err != nil {
		t.Error(err)
	}

	examples, err := db.ReadLabeledExamples(10)
	if err != nil {
		t.Error(err)
	}
	if len(examples) != 2 {
		t.Errorf("len(examples) == %d, want 2", len(examples))
	}
}

func TestReadRecentExamples(t *testing.T) {
	_, err := db.DeleteAllExamples()
	if err != nil {
		t.Error(err)
	}

	_, err = db.InsertOrUpdateExample(example.NewExample("http://hoge1.com", model.POSITIVE))
	if err != nil {
		t.Error(err)
	}
	_, err = db.InsertOrUpdateExample(example.NewExample("http://hoge2.com", model.NEGATIVE))
	if err != nil {
		t.Error(err)
	}
	_, err = db.InsertOrUpdateExample(example.NewExample("http://hoge3.com", model.UNLABELED))
	if err != nil {
		t.Error(err)
	}

	examples, err := db.ReadRecentExamples(time.Now().Add(time.Duration(-10) * time.Minute))
	if err != nil {
		t.Error(err)
	}
	if len(examples) != 3 {
		t.Errorf("len(examples) == %d, want 3", len(examples))
	}
}

func TestSearchExamplesByUlr(t *testing.T) {
	_, err := db.DeleteAllExamples()
	if err != nil {
		t.Error(err)
	}

	_, err = db.InsertOrUpdateExample(example.NewExample("http://hoge1.com", model.NEGATIVE))
	if err != nil {
		t.Error(err)
	}
	_, err = db.InsertOrUpdateExample(example.NewExample("http://hoge2.com", model.NEGATIVE))
	if err != nil {
		t.Error(err)
	}
	_, err = db.InsertOrUpdateExample(example.NewExample("http://hoge3.com", model.UNLABELED))
	if err != nil {
		t.Error(err)
	}

	example, err := db.SearchExamplesByUlr("http://hoge1.com")
	if err != nil {
		t.Error(err)
	}
	if example.Url == "" {
		t.Errorf("example.Url == %s, want http://hoge1.com", example.Url)
	}

	example, err = db.SearchExamplesByUlr("http://hoge4.com")
	if err == nil {
		t.Errorf("search result must be nil")
	}
}

func TestSearchExamplesByUlrs(t *testing.T) {
	_, err := db.DeleteAllExamples()
	if err != nil {
		t.Error(err)
	}

	_, err = db.InsertOrUpdateExample(example.NewExample("http://hoge1.com", model.NEGATIVE))
	if err != nil {
		t.Error(err)
	}
	_, err = db.InsertOrUpdateExample(example.NewExample("http://hoge2.com", model.NEGATIVE))
	if err != nil {
		t.Error(err)
	}
	_, err = db.InsertOrUpdateExample(example.NewExample("http://hoge3.com", model.UNLABELED))
	if err != nil {
		t.Error(err)
	}

	examples, err := db.SearchExamplesByUlrs([]string{"http://hoge1.com", "http://hoge2.com"})
	if err != nil {
		t.Error(err)
	}
	if len(examples) != 2 {
		t.Errorf("len(examples) == %d, want 2", len(examples))
	}
}

func TestSearchExamplesByLabels(t *testing.T) {
	_, err := db.DeleteAllExamples()
	if err != nil {
		t.Error(err)
	}

	_, err = db.InsertOrUpdateExample(example.NewExample("http://hoge1.com", model.POSITIVE))
	if err != nil {
		t.Error(err)
	}
	_, err = db.InsertOrUpdateExample(example.NewExample("http://hoge2.com", model.NEGATIVE))
	if err != nil {
		t.Error(err)
	}
	_, err = db.InsertOrUpdateExample(example.NewExample("http://hoge3.com", model.UNLABELED))
	if err != nil {
		t.Error(err)
	}

	examples, err := db.ReadPositiveExamples(10)
	if err != nil {
		t.Error(err)
	}
	if len(examples) != 1 {
		t.Errorf("len(examples) == %d, want 1", len(examples))
	}

	examples, err = db.ReadNegativeExamples(10)
	if err != nil {
		t.Error(err)
	}
	if len(examples) != 1 {
		t.Errorf("len(examples) == %d, want 1", len(examples))
	}

	examples, err = db.ReadUnlabeledExamples(10)
	if err != nil {
		t.Error(err)
	}
	if len(examples) != 1 {
		t.Errorf("len(examples) == %d, want 1", len(examples))
	}
}
