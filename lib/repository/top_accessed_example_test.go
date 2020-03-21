package repository_test

import (
	"testing"

	"github.com/syou6162/go-active-learning/lib/example"
	"github.com/syou6162/go-active-learning/lib/model"
	"github.com/syou6162/go-active-learning/lib/repository"
)

func TestUpdateTopAccessedExamples(t *testing.T) {
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
	err = repo.UpdateTopAccessedExamples(examples)
	if err != nil {
		t.Error(err)
	}

	{
		top, err := repo.SearchTopAccessedExampleIds()
		if err != nil {
			t.Error(err)
		}
		if len(top) != 3 {
			t.Error("len(top) must be 3")
		}
	}
}
