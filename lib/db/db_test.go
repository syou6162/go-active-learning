package db_test

import (
	"log"
	"os"
	"testing"
	"time"

	"github.com/syou6162/go-active-learning/lib/db"
	"github.com/syou6162/go-active-learning/lib/example"
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

	_, err = db.InsertOrUpdateExample(&example.Example{Url: "http://hoge.com", Label: example.NEGATIVE})
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

	// same url
	_, err = db.InsertOrUpdateExample(&example.Example{Url: "http://hoge.com", Label: example.NEGATIVE})
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

	// same url but different label
	_, err = db.InsertOrUpdateExample(&example.Example{Url: "http://hoge.com", Label: example.POSITIVE})
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

	// different url
	_, err = db.InsertOrUpdateExample(&example.Example{Url: "http://another.com", Label: example.NEGATIVE})
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

	_, err = db.InsertOrUpdateExample(&example.Example{Url: "http://hoge1.com", Label: example.NEGATIVE})
	if err != nil {
		t.Error(err)
	}
	_, err = db.InsertOrUpdateExample(&example.Example{Url: "http://hoge2.com", Label: example.NEGATIVE})
	if err != nil {
		t.Error(err)
	}
	_, err = db.InsertOrUpdateExample(&example.Example{Url: "http://hoge3.com", Label: example.UNLABELED})
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

	_, err = db.InsertOrUpdateExample(&example.Example{Url: "http://hoge1.com", Label: example.NEGATIVE})
	if err != nil {
		t.Error(err)
	}
	_, err = db.InsertOrUpdateExample(&example.Example{Url: "http://hoge2.com", Label: example.NEGATIVE})
	if err != nil {
		t.Error(err)
	}
	_, err = db.InsertOrUpdateExample(&example.Example{Url: "http://hoge3.com", Label: example.UNLABELED})
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

func TestSearchExamplesByUlrs(t *testing.T) {
	_, err := db.DeleteAllExamples()
	if err != nil {
		t.Error(err)
	}

	_, err = db.InsertOrUpdateExample(&example.Example{Url: "http://hoge1.com", Label: example.NEGATIVE})
	if err != nil {
		t.Error(err)
	}
	_, err = db.InsertOrUpdateExample(&example.Example{Url: "http://hoge2.com", Label: example.NEGATIVE})
	if err != nil {
		t.Error(err)
	}
	_, err = db.InsertOrUpdateExample(&example.Example{Url: "http://hoge3.com", Label: example.UNLABELED})
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
