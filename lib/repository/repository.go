package repository

import (
	"database/sql"
	"fmt"
	"io"
	"time"

	"bufio"

	_ "github.com/lib/pq"
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
	SearchExamplesByUlr(url string) (*model.Example, error)
	SearchExamplesByUlrs(urls []string) (model.Examples, error)
	DeleteAllExamples() error
	Ping() error
	Close() error
}

type repository struct {
	db *sql.DB
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
	db, err := sql.Open("postgres", GetDataSourceName())
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