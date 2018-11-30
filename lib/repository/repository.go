package repository

import (
	"fmt"
	"io"
	"time"

	"github.com/jmoiron/sqlx"

	"bufio"

	_ "github.com/lib/pq"
	"github.com/syou6162/go-active-learning/lib/feature"
	"github.com/syou6162/go-active-learning/lib/model"
	"github.com/syou6162/go-active-learning/lib/util"
)

type Repository interface {
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
	FindExampleByUlr(url string) (*model.Example, error)
	SearchExamplesByUlrs(urls []string) (model.Examples, error)
	SearchExamplesByKeywords(keywords []string, aggregator string, limit int) (model.Examples, error)
	DeleteAllExamples() error

	UpdateFeatureVector(e *model.Example) error
	FindFeatureVector(e *model.Example) (feature.FeatureVector, error)
	SearchFeatureVector(examples model.Examples) ([]feature.FeatureVector, error)

	UpdateHatenaBookmark(e *model.Example) error
	SearchHatenaBookmarks(examples model.Examples) ([]*model.HatenaBookmark, error)
	FindHatenaBookmark(e *model.Example) (*model.HatenaBookmark, error)

	UpdateReferringTweets(e *model.Example) error
	SearchReferringTweetsList(examples model.Examples) (map[int]model.ReferringTweets, error)
	FindReferringTweets(e *model.Example) (model.ReferringTweets, error)

	Ping() error
	Close() error
}

type repository struct {
	db *sqlx.DB
}

func GetDataSourceName() string {
	host := util.GetEnv("POSTGRES_HOST", "localhost")
	dbUser := util.GetEnv("DB_USER", "nobody")
	dbPassword := util.GetEnv("DB_PASSWORD", "nobody")
	dbName := util.GetEnv("DB_NAME", "go-active-learning")
	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s sslmode=disable",
		host, dbUser, dbPassword, dbName,
	)
}

func New() (*repository, error) {
	db, err := sqlx.Open("postgres", GetDataSourceName())
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(50)
	return &repository{db: db}, nil
}

func (r *repository) Ping() error {
	return r.db.Ping()
}

func (r *repository) Close() error {
	if r.db != nil {
		return r.db.Close()
	} else {
		return nil
	}
}
