package repository_test

import (
	"log"
	"os"
	"testing"
	"time"

	"github.com/syou6162/go-active-learning/lib/example"
	"github.com/syou6162/go-active-learning/lib/feature"
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

	examples, err := repo.SearchExamples()
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

	examples, err := repo.SearchExamples()
	if err != nil {
		t.Error(err)
	}
	if len(examples) != 1 {
		t.Errorf("len(examples) == %d, want 1", len(examples))
	}
	if examples[0].Label != model.UNLABELED {
		t.Errorf("label == %d, want 1", examples[0].Label)
	}
	if examples[0].Id == 0 {
		t.Error("id must not be 0")
	}

	// same url
	err = repo.InsertOrUpdateExample(example.NewExample("http://hoge.com", model.NEGATIVE))
	if err != nil {
		t.Error(err)
	}

	examples, err = repo.SearchExamples()
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

	examples, err = repo.SearchExamples()
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

	examples, err = repo.SearchExamples()
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

	examples, err = repo.SearchExamples()
	if err != nil {
		t.Error(err)
	}
	if len(examples) != 2 {
		t.Errorf("len(examples) == %d, want 2", len(examples))
	}
}

func TestUpdateScore(t *testing.T) {
	repo, err := repository.New()
	if err != nil {
		t.Errorf(err.Error())
	}
	defer repo.Close()

	if err = repo.DeleteAllExamples(); err != nil {
		t.Error(err)
	}

	url := "http://hoge.com"
	e := example.NewExample(url, model.UNLABELED)
	e.Score = 1.0
	err = repo.InsertOrUpdateExample(e)
	if err != nil {
		t.Error(err)
	}

	e, err = repo.FindExampleByUlr(url)
	if err != nil {
		t.Error(err)
	}
	if e.Score != 1.0 {
		t.Errorf("e.Score == %f, want 1.0", e.Score)
	}

	e.Score = 100.0
	err = repo.UpdateScore(e)
	if err != nil {
		t.Error(err)
	}

	e, err = repo.FindExampleByUlr(url)
	if err != nil {
		t.Error(err)
	}
	if e.Score != 100.0 {
		t.Errorf("e.Score == %f, want 100.0", e.Score)
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

	examples, err := repo.SearchLabeledExamples(10)
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

	examples, err := repo.SearchRecentExamples(time.Now().Add(time.Duration(-10) * time.Minute), 10)
	if err != nil {
		t.Error(err)
	}
	if len(examples) != 3 {
		t.Errorf("len(examples) == %d, want 3", len(examples))
	}
}

func TestReadRecentExamplesByHost(t *testing.T) {
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

	examples, err := repo.SearchRecentExamplesByHost("http://hoge1.com", time.Now().Add(time.Duration(-10)*time.Minute), 10)
	if err != nil {
		t.Error(err)
	}
	if len(examples) != 1 {
		t.Errorf("len(examples) == %d, want 1", len(examples))
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

	example, err := repo.FindExampleByUlr("http://hoge1.com")
	if err != nil {
		t.Error(err)
	}
	if example.Url == "" {
		t.Errorf("example.Url == %s, want http://hoge1.com", example.Url)
	}

	example, err = repo.FindExampleByUlr("http://hoge4.com")
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

	examples, err := repo.SearchPositiveExamples(10)
	if err != nil {
		t.Error(err)
	}
	if len(examples) != 1 {
		t.Errorf("len(examples) == %d, want 1", len(examples))
	}

	examples, err = repo.SearchNegativeExamples(10)
	if err != nil {
		t.Error(err)
	}
	if len(examples) != 1 {
		t.Errorf("len(examples) == %d, want 1", len(examples))
	}

	examples, err = repo.SearchUnlabeledExamples(10)
	if err != nil {
		t.Error(err)
	}
	if len(examples) != 1 {
		t.Errorf("len(examples) == %d, want 1", len(examples))
	}
}

func TestFeatureVectorReadWrite(t *testing.T) {
	repo, err := repository.New()
	if err != nil {
		t.Errorf(err.Error())
	}
	defer repo.Close()

	if err = repo.DeleteAllExamples(); err != nil {
		t.Error(err)
	}

	e1 := example.NewExample("http://hoge.com", model.UNLABELED)
	err = repo.InsertOrUpdateExample(e1)
	if err != nil {
		t.Error(err)
	}
	e1.Fv = feature.FeatureVector{"BIAS"}

	if err = repo.UpdateFeatureVector(e1); err != nil {
		t.Error(err)
	}

	fv, err := repo.FindFeatureVector(e1)
	if err != nil {
		t.Error(err)
	}
	if len(fv) != 1 {
		t.Errorf("len(fv) == %d, want 1", len(fv))
	}

	e2 := example.NewExample("http://fuga.com", model.UNLABELED)
	err = repo.InsertOrUpdateExample(e2)
	if err != nil {
		t.Error(err)
	}
	e2.Fv = feature.FeatureVector{"hoge"}
	if err = repo.UpdateFeatureVector(e2); err != nil {
		t.Error(err)
	}
	fvList, err := repo.SearchFeatureVector(model.Examples{e1, e2})
	if err != nil {
		t.Error(err)
	}
	if len(fvList) != 2 {
		t.Errorf("len(fvList) == %d, want 2", len(fvList))
	}
	if fvList[e2.Id][0] != "hoge" {
		t.Errorf("fvList[e2.Id][0] == %s, want hoge", fvList[e2.Id][0])
	}
}

func TestSearchExamplesByWords(t *testing.T) {
	repo, err := repository.New()
	if err != nil {
		t.Errorf(err.Error())
	}
	defer repo.Close()

	if err = repo.DeleteAllExamples(); err != nil {
		t.Error(err)
	}

	e1 := example.NewExample("http://hoge.com", model.UNLABELED)
	e1.Title = "日本語"
	err = repo.InsertOrUpdateExample(e1)
	if err != nil {
		t.Error(err)
	}

	e2 := example.NewExample("http://fuga.com", model.UNLABELED)
	e2.Title = "英語"
	err = repo.InsertOrUpdateExample(e2)
	if err != nil {
		t.Error(err)
	}

	examples, err := repo.SearchExamplesByKeywords([]string{"日本語"}, "ALL", 100)
	if len(examples) != 1 {
		t.Errorf("len(examples) == %d, want 1", len(examples))
	}
	examples, err = repo.SearchExamplesByKeywords([]string{"語"}, "ALL", 100)
	if len(examples) != 2 {
		t.Errorf("len(examples) == %d, want 2", len(examples))
	}
	examples, err = repo.SearchExamplesByKeywords([]string{"日本語", "英語"}, "ALL", 100)
	if len(examples) != 0 {
		t.Errorf("len(examples) == %d, want 0", len(examples))
	}
	examples, err = repo.SearchExamplesByKeywords([]string{"日本語", "英語"}, "ANY", 100)
	if len(examples) != 2 {
		t.Errorf("len(examples) == %d, want 2", len(examples))
	}
}
