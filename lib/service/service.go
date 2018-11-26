package service

import (
	"github.com/syou6162/go-active-learning/lib/repository"
	"github.com/syou6162/go-active-learning/lib/model"
)

type GoActiveLearningApp interface {
	InsertOrUpdateExample(e *model.Example) error
	DeleteAllExamples() error
	Close() error
}

func NewApp(repo repository.Repository) GoActiveLearningApp {
	return &goActiveLearningApp{repo}
}

type goActiveLearningApp struct {
	repo         repository.Repository
}

func (app *goActiveLearningApp) Close() error {
	return app.repo.Close()
}