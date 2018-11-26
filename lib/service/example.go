package service

import (
	"github.com/syou6162/go-active-learning/lib/model"
)

func (app *goActiveLearningApp) InsertOrUpdateExample(e *model.Example) error {
	return app.repo.InsertOrUpdateExample(e)
}

func (app *goActiveLearningApp) DeleteAllExamples() error {
	return app.repo.DeleteAllExamples()
}
