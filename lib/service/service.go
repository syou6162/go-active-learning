package service

import (
	"bufio"
	"io"
	"time"

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
	FindExampleById(id int) (*model.Example, error)
	SearchExamplesByUlrs(urls []string) (model.Examples, error)
	SearchExamplesByIds(ids []int) (model.Examples, error)
	SearchExamplesByKeywords(keywords []string, aggregator string, limit int) (model.Examples, error)
	DeleteAllExamples() error
	CountPositiveExamples() (int, error)
	CountNegativeExamples() (int, error)
	CountUnlabeledExamples() (int, error)

	InsertMIRAModel(m classifier.MIRAClassifier) error
	FindLatestMIRAModel(modelType classifier.ModelType) (*classifier.MIRAClassifier, error)

	UpdateFeatureVector(e *model.Example) error
	UpdateHatenaBookmark(e *model.Example) error
	UpdateOrCreateReferringTweets(e *model.Example) error
	UpdateTweetLabel(exampleId int, idStr string, label model.LabelType) error
	SearchReferringTweets(limit int) (model.ReferringTweets, error)
	SearchPositiveReferringTweets(scoreThreshold float64, tweetsLimitInSameExample int, limit int) (model.ReferringTweets, error)
	SearchNegativeReferringTweets(scoreThreshold float64, tweetsLimitInSameExample int, limit int) (model.ReferringTweets, error)
	SearchUnlabeledReferringTweets(scoreThreshold float64, tweetsLimitInSameExample int, limit int) (model.ReferringTweets, error)
	SearchRecentReferringTweetsWithHighScore(from time.Time, scoreThreshold float64, limit int) (model.ReferringTweets, error)
	Fetch(examples model.Examples)

	AttachMetadataIncludingFeatureVector(examples model.Examples, bookmarkLimit int, tweetLimit int) error
	AttachMetadata(examples model.Examples, bookmarkLimit, tweetLimit int) error

	UpdateRecommendation(listName string, examples model.Examples) error
	GetRecommendation(listName string) (model.Examples, error)

	UpdateRelatedExamples(related model.RelatedExamples) error
	SearchRelatedExamples(e *model.Example) (model.Examples, error)

	UpdateTopAccessedExampleIds(exampleIds []int) error
	SearchTopAccessedExamples() (model.Examples, error)

	Ping() error
	Close() error
}

func NewApp(repo repository.Repository) GoActiveLearningApp {
	return &goActiveLearningApp{repo: repo}
}

func NewDefaultApp() (GoActiveLearningApp, error) {
	repo, err := repository.New()
	if err != nil {
		return nil, err
	}
	return &goActiveLearningApp{repo: repo}, nil
}

type goActiveLearningApp struct {
	repo repository.Repository
}

func (app *goActiveLearningApp) InsertMIRAModel(m classifier.MIRAClassifier) error {
	return app.repo.InsertMIRAModel(m)
}

func (app *goActiveLearningApp) FindLatestMIRAModel(modelType classifier.ModelType) (*classifier.MIRAClassifier, error) {
	return app.repo.FindLatestMIRAModel(modelType)
}

func (app *goActiveLearningApp) Ping() error {
	if err := app.repo.Ping(); err != nil {
		return err
	}
	return nil
}

func (app *goActiveLearningApp) Close() error {
	if err := app.repo.Close(); err != nil {
		return err
	}
	return nil
}
