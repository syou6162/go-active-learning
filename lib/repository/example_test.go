package repository_test

import (
	"log"
	"os"
	"testing"
	"time"

	"github.com/syou6162/go-active-learning/lib/example"
	"github.com/syou6162/go-active-learning/lib/model"
	"github.com/syou6162/go-active-learning/lib/repository"
)

func TestMain(m *testing.M) {
	repo, err := repository.New()
	if err != nil {
		log.Fatal(err.Error())
	}
	defer repo.Close()

	ret := m.Run()
	os.Exit(ret)
}

func TestPing(t *testing.T) {
	repo, err := repository.New()
	if err != nil {
		t.Errorf(err.Error())
	}
	defer repo.Close()

	if err := repo.Ping(); err != nil {
		t.Errorf(err.Error())
	}
}

func TestInsertExamplesFromReader(t *testing.T) {
	repo, err := repository.New()
	if err != nil {
		t.Errorf(err.Error())
	}
	defer repo.Close()

	if err = repo.DeleteAllExamples(); err != nil {
		t.Error(err)
	}

	fp, err := os.Open("../../tech_input_example.txt")
	defer fp.Close()
	if err != nil {
		t.Error(err)
	}
	repo.InsertExamplesFromReader(fp)

	examples, err := repo.ReadExamples()
	if err != nil {
		t.Error(err)
	}
	if len(examples) == 0 {
		t.Errorf("len(examples) > 0, but %d", len(examples))
	}
}

func TestInsertOrUpdateExample(t *testing.T) {
	repo, err := repository.New()
	if err != nil {
		t.Errorf(err.Error())
	}
	defer repo.Close()

	if err = repo.DeleteAllExamples(); err != nil {
		t.Error(err)
	}

	err = repo.InsertOrUpdateExample(example.NewExample("http://hoge.com", model.UNLABELED))
	if err != nil {
		t.Error(err)
	}

	examples, err := repo.ReadExamples()
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
	err = repo.InsertOrUpdateExample(example.NewExample("http://hoge.com", model.NEGATIVE))
	if err != nil {
		t.Error(err)
	}

	examples, err = repo.ReadExamples()
	if err != nil {
		t.Error(err)
	}
	if len(examples) != 1 {
		t.Errorf("len(examples) == %d, want 1", len(examples))
	}
	if examples[0].Label != model.NEGATIVE {
		t.Errorf("label == %d, want -1", examples[0].Label)
	}

	// same url but different label
	err = repo.InsertOrUpdateExample(example.NewExample("http://hoge.com", model.POSITIVE))
	if err != nil {
		t.Error(err)
	}

	examples, err = repo.ReadExamples()
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
	err = repo.InsertOrUpdateExample(example.NewExample("http://hoge.com", model.UNLABELED))
	if err != nil {
		t.Error(err)
	}

	examples, err = repo.ReadExamples()
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
	err = repo.InsertOrUpdateExample(example.NewExample("http://another.com", model.NEGATIVE))
	if err != nil {
		t.Error(err)
	}

	examples, err = repo.ReadExamples()
	if err != nil {
		t.Error(err)
	}
	if len(examples) != 2 {
		t.Errorf("len(examples) == %d, want 2", len(examples))
	}
}

func TestReadLabeledExamples(t *testing.T) {
	repo, err := repository.New()
	if err != nil {
		t.Errorf(err.Error())
	}
	defer repo.Close()

	if err = repo.DeleteAllExamples(); err != nil {
		t.Error(err)
	}

	err = repo.InsertOrUpdateExample(example.NewExample("http://hoge1.com", model.POSITIVE))
	if err != nil {
		t.Error(err)
	}
	err = repo.InsertOrUpdateExample(example.NewExample("http://hoge2.com", model.NEGATIVE))
	if err != nil {
		t.Error(err)
	}
	err = repo.InsertOrUpdateExample(example.NewExample("http://hoge3.com", model.UNLABELED))
	if err != nil {
		t.Error(err)
	}

	examples, err := repo.ReadLabeledExamples(10)
	if err != nil {
		t.Error(err)
	}
	if len(examples) != 2 {
		t.Errorf("len(examples) == %d, want 2", len(examples))
	}
}

func TestReadRecentExamples(t *testing.T) {
	repo, err := repository.New()
	if err != nil {
		t.Errorf(err.Error())
	}
	defer repo.Close()

	if err = repo.DeleteAllExamples(); err != nil {
		t.Error(err)
	}

	err = repo.InsertOrUpdateExample(example.NewExample("http://hoge1.com", model.POSITIVE))
	if err != nil {
		t.Error(err)
	}
	err = repo.InsertOrUpdateExample(example.NewExample("http://hoge2.com", model.NEGATIVE))
	if err != nil {
		t.Error(err)
	}
	err = repo.InsertOrUpdateExample(example.NewExample("http://hoge3.com", model.UNLABELED))
	if err != nil {
		t.Error(err)
	}

	examples, err := repo.ReadRecentExamples(time.Now().Add(time.Duration(-10) * time.Minute))
	if err != nil {
		t.Error(err)
	}
	if len(examples) != 3 {
		t.Errorf("len(examples) == %d, want 3", len(examples))
	}
}

func TestSearchExamplesByUlr(t *testing.T) {
	repo, err := repository.New()
	if err != nil {
		t.Errorf(err.Error())
	}
	defer repo.Close()

	if err = repo.DeleteAllExamples(); err != nil {
		t.Error(err)
	}

	err = repo.InsertOrUpdateExample(example.NewExample("http://hoge1.com", model.NEGATIVE))
	if err != nil {
		t.Error(err)
	}
	err = repo.InsertOrUpdateExample(example.NewExample("http://hoge2.com", model.NEGATIVE))
	if err != nil {
		t.Error(err)
	}
	err = repo.InsertOrUpdateExample(example.NewExample("http://hoge3.com", model.UNLABELED))
	if err != nil {
		t.Error(err)
	}

	example, err := repo.SearchExamplesByUlr("http://hoge1.com")
	if err != nil {
		t.Error(err)
	}
	if example.Url == "" {
		t.Errorf("example.Url == %s, want http://hoge1.com", example.Url)
	}

	example, err = repo.SearchExamplesByUlr("http://hoge4.com")
	if err == nil {
		t.Errorf("search result must be nil")
	}
}

func TestSearchExamplesByUlrs(t *testing.T) {
	repo, err := repository.New()
	if err != nil {
		t.Errorf(err.Error())
	}
	defer repo.Close()

	if err = repo.DeleteAllExamples(); err != nil {
		t.Error(err)
	}

	err = repo.InsertOrUpdateExample(example.NewExample("http://hoge1.com", model.NEGATIVE))
	if err != nil {
		t.Error(err)
	}
	err = repo.InsertOrUpdateExample(example.NewExample("http://hoge2.com", model.NEGATIVE))
	if err != nil {
		t.Error(err)
	}
	err = repo.InsertOrUpdateExample(example.NewExample("http://hoge3.com", model.UNLABELED))
	if err != nil {
		t.Error(err)
	}

	examples, err := repo.SearchExamplesByUlrs([]string{"http://hoge1.com", "http://hoge2.com"})
	if err != nil {
		t.Error(err)
	}
	if len(examples) != 2 {
		t.Errorf("len(examples) == %d, want 2", len(examples))
	}
}

func TestSearchExamplesByLabels(t *testing.T) {
	repo, err := repository.New()
	if err != nil {
		t.Errorf(err.Error())
	}
	defer repo.Close()

	if err = repo.DeleteAllExamples(); err != nil {
		t.Error(err)
	}

	err = repo.InsertOrUpdateExample(example.NewExample("http://hoge1.com", model.POSITIVE))
	if err != nil {
		t.Error(err)
	}
	err = repo.InsertOrUpdateExample(example.NewExample("http://hoge2.com", model.NEGATIVE))
	if err != nil {
		t.Error(err)
	}
	err = repo.InsertOrUpdateExample(example.NewExample("http://hoge3.com", model.UNLABELED))
	if err != nil {
		t.Error(err)
	}

	examples, err := repo.ReadPositiveExamples(10)
	if err != nil {
		t.Error(err)
	}
	if len(examples) != 1 {
		t.Errorf("len(examples) == %d, want 1", len(examples))
	}

	examples, err = repo.ReadNegativeExamples(10)
	if err != nil {
		t.Error(err)
	}
	if len(examples) != 1 {
		t.Errorf("len(examples) == %d, want 1", len(examples))
	}

	examples, err = repo.ReadUnlabeledExamples(10)
	if err != nil {
		t.Error(err)
	}
	if len(examples) != 1 {
		t.Errorf("len(examples) == %d, want 1", len(examples))
	}
}
