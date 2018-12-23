package service

import (
	"bufio"
	"io"
	"time"

	"github.com/syou6162/go-active-learning/lib/cache"
	"github.com/syou6162/go-active-learning/lib/classifier"
	"github.com/syou6162/go-active-learning/lib/model"
	"github.com/syou6162/go-active-learning/lib/repository"
)

type GoActiveLearningApp interface {
	UpdateOrCreateExample(e *model.Example) error
	UpdateScore(e *model.Example) error
	InsertExampleFromScanner(scanner *bufio.Scanner) (*model.Example, error)
	InsertExamplesFromReader(reader io.Reader) error
	SearchExamples() (model.Examples, error)
	SearchRecentExamples(from time.Time, limit int) (model.Examples, error)
	SearchRecentExamplesByHost(host string, from time.Time, limit int) (model.Examples, error)
	SearchExamplesByLabel(label model.LabelType, limit int) (model.Examples, error)
	SearchLabeledExamples(limit int) (model.Examples, error)
	SearchPositiveExamples(limit int) (model.Examples, error)
	SearchNegativeExamples(limit int) (model.Examples, error)
	SearchUnlabeledExamples(limit int) (model.Examples, error)
	SearchPositiveScoredExamples(limit int) (model.Examples, error)
	FindExampleByUlr(url string) (*model.Example, error)
	SearchExamplesByUlrs(urls []string) (model.Examples, error)
	SearchExamplesByKeywords(keywords []string, aggregator string, limit int) (model.Examples, error)
	DeleteAllExamples() error

	InsertMIRAModel(m classifier.MIRAClassifier) error
	FindLatestMIRAModel() (*classifier.MIRAClassifier, error)

	UpdateExampleMetadata(e model.Example) error
	UpdateExamplesMetadata(examples model.Examples) error
	UpdateFeatureVector(e *model.Example) error
	UpdateHatenaBookmark(e *model.Example) error
	UpdateReferringTweets(e *model.Example) error
	Fetch(examples model.Examples)

	AttachMetadata(examples model.Examples) error
	AttachLightMetadata(examples model.Examples) error

	AddExamplesToList(listName string, examples model.Examples) error
	GetUrlsFromList(listName string, from int64, to int64) ([]string, error)

	Ping() error
	Close() error
}

func NewApp(repo repository.Repository, c cache.Cache) GoActiveLearningApp {
	return &goActiveLearningApp{repo: repo, cache: c}
}

func NewDefaultApp() (GoActiveLearningApp, error) {
	repo, err := repository.New()
	if err != nil {
		return nil, err
	}

	c, err := cache.New()
	if err != nil {
		return nil, err
	}

	return &goActiveLearningApp{repo: repo, cache: c}, nil
}

type goActiveLearningApp struct {
	repo  repository.Repository
	cache cache.Cache
}

func (app *goActiveLearningApp) InsertMIRAModel(m classifier.MIRAClassifier) error {
	return app.repo.InsertMIRAModel(m)
}

func (app *goActiveLearningApp) FindLatestMIRAModel() (*classifier.MIRAClassifier, error) {
	return app.repo.FindLatestMIRAModel()
}

func (app *goActiveLearningApp) Ping() error {
	if err := app.repo.Ping(); err != nil {
		return err
	}
	if err := app.cache.Ping(); err != nil {
		return err
	}
	return nil
}

func (app *goActiveLearningApp) Close() error {
	if err := app.repo.Close(); err != nil {
		return err
	}
	if err := app.cache.Close(); err != nil {
		return err
	}
	return nil
}
