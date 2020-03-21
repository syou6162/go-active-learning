package repository_test

import (
	"testing"

	"github.com/syou6162/go-active-learning/lib/example"
	"github.com/syou6162/go-active-learning/lib/model"
	"github.com/syou6162/go-active-learning/lib/repository"
)

func TestUpdateRelatedExamples(t *testing.T) {
	repo, err := repository.New()
	if err != nil {
		t.Errorf(err.Error())
	}
	defer repo.Close()

	e1 := example.NewExample("http://hoge1.com", model.POSITIVE)
	e2 := example.NewExample("http://hoge2.com", model.NEGATIVE)
	e3 := example.NewExample("http://hoge3.com", model.UNLABELED)
	examples := model.Examples{e1, e2, e3}
	for _, e := range examples {
		err = repo.UpdateOrCreateExample(e)
		if err != nil {
			t.Error(err)
		}
	}
	related := model.RelatedExamples{ExampleId: e1.Id, RelatedExampleIds: []int{e2.Id, e3.Id}}
	err = repo.UpdateRelatedExamples(related)
	if err != nil {
		t.Error(err)
	}

	{
		related, err := repo.FindRelatedExamples(e1)
		if err != nil {
			t.Error(err)
		}
		if len(related.RelatedExampleIds) != 2 {
			t.Error("len(related.RelatedExampleIds) must be 2")
		}
	}
}

func TestUpdateRelatedExamplesMyOwn(t *testing.T) {
	repo, err := repository.New()
	if err != nil {
		t.Errorf(err.Error())
	}
	defer repo.Close()

	e1 := example.NewExample("http://hoge1.com", model.POSITIVE)
	e2 := example.NewExample("http://hoge2.com", model.NEGATIVE)
	e3 := example.NewExample("http://hoge3.com", model.UNLABELED)
	examples := model.Examples{e1, e2, e3}
	for _, e := range examples {
		err = repo.UpdateOrCreateExample(e)
		if err != nil {
			t.Error(err)
		}
	}
	related := model.RelatedExamples{ExampleId: e1.Id, RelatedExampleIds: []int{e1.Id, e2.Id, e3.Id}}
	err = repo.UpdateRelatedExamples(related)
	if err == nil {
		t.Error("自身と同一のexample_idを持つ事例はrelated_example_idに追加できない")
	}
}
