package featureweight_test

import (
	"testing"

	"github.com/codegangsta/cli"
	"github.com/syou6162/go-active-learning/lib/cache"
	"github.com/syou6162/go-active-learning/lib/command"
	"github.com/syou6162/go-active-learning/lib/repository"
	"github.com/syou6162/go-active-learning/lib/service"
	"github.com/syou6162/go-active-learning/lib/util/file"
)

func TestDoListFeatureWeight(t *testing.T) {
	inputFilename := "../../../tech_input_example.txt"
	train, err := file.ReadExamples(inputFilename)
	if err != nil {
		t.Error(err)
	}

	repo, err := repository.New()
	if err != nil {
		t.Error(err)
	}

	cache_, err := cache.New()
	if err != nil {
		t.Error(err)
	}
	defer cache_.Close()

	a := service.NewApp(repo, cache_)
	defer a.Close()

	if err = a.DeleteAllExamples(); err != nil {
		t.Error(err)
	}

	for _, example := range train {
		if err = a.InsertOrUpdateExample(example); err != nil {
			t.Error(err)
		}
	}

	app := cli.NewApp()
	app.Commands = command.Commands
	args := []string{
		"go-active-learning",
		"diagnose",
		"feature-weight",
		"--filter-status-code-ok",
	}

	if err := app.Run(args); err != nil {
		t.Error(err)
	}
}
