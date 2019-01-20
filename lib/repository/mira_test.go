package repository_test

import (
	"testing"

	"github.com/syou6162/go-active-learning/lib/classifier"
	"github.com/syou6162/go-active-learning/lib/repository"
)

func TestInsertMIRAModel(t *testing.T) {
	repo, err := repository.New()
	if err != nil {
		t.Errorf(err.Error())
	}
	defer repo.Close()

	weight := make(map[string]float64)
	weight["hoge"] = 1.0
	weight["fuga"] = 1.0
	clf := classifier.MIRAClassifier{classifier.EXAMPLE, weight, 10.0, 0.0, 0.0, 0.0, 0.0}
	err = repo.InsertMIRAModel(clf)
	if err != nil {
		t.Error(err)
	}

	{
		clf, err := repo.FindLatestMIRAModel()
		if err != nil {
			t.Error(err)
		}
		if len(clf.Weight) == 0 {
			t.Error("weight must not be empty")
		}
		if clf.C != 10.0 {
			t.Error("C must be 10.0")
		}
	}
}
