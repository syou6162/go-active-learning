package service

import (
	"bufio"
	"io"
	"time"

	"github.com/syou6162/go-active-learning/lib/model"
	"github.com/syou6162/go-active-learning/lib/repository"
)

type GoActiveLearningApp interface {
	InsertOrUpdateExample(e *model.Example) error
	InsertExampleFromScanner(scanner *bufio.Scanner) (*model.Example, error)
	InsertExamplesFromReader(reader io.Reader) error
	ReadExamples() (model.Examples, error)
	ReadRecentExamples(from time.Time) (model.Examples, error)
	ReadExamplesByLabel(label model.LabelType, limit int) (model.Examples, error)
	ReadLabeledExamples(limit int) (model.Examples, error)
	ReadPositiveExamples(limit int) (model.Examples, error)
	ReadNegativeExamples(limit int) (model.Examples, error)
	ReadUnlabeledExamples(limit int) (model.Examples, error)
	SearchExamplesByUlr(url string) (*model.Example, error)
	SearchExamplesByUlrs(urls []string) (model.Examples, error)
	DeleteAllExamples() error
	Ping() error
	Close() error
}

func NewApp(repo repository.Repository) GoActiveLearningApp {
	return &goActiveLearningApp{repo}
}

type goActiveLearningApp struct {
	repo repository.Repository
}

func (app *goActiveLearningApp) Ping() error {
	return app.repo.Ping()
}

func (app *goActiveLearningApp) Close() error {
	return app.repo.Close()
}
