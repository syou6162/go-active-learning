package db

import (
	"bufio"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"io"

	"github.com/lib/pq"
	_ "github.com/lib/pq"
	"github.com/syou6162/go-active-learning/lib/example"
	"github.com/syou6162/go-active-learning/lib/util"
	"github.com/syou6162/go-active-learning/lib/util/file"
)

var db *sql.DB

func Init() error {
	if db != nil {
		return errors.New("Init method is called more than twice")
	}
	host := util.GetEnv("POSTGRES_HOST", "localhost")
	dbUser := util.GetEnv("DB_USER", "nobody")
	dbPassword := util.GetEnv("DB_PASSWORD", "nobody")
	dbName := util.GetEnv("DB_NAME", "go-active-learning")
	var err error
	db, err = sql.Open("postgres", fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable", host, dbUser, dbPassword, dbName))
	if err != nil {
		return err
	}
	db.SetMaxOpenConns(50)
	return nil
}

func Close() error {
	if db != nil {
		return db.Close()
	} else {
		return nil
	}
}

func InsertOrUpdateExample(e *example.Example) (sql.Result, error) {
	var label example.LabelType
	now := time.Now()

	url := e.FinalUrl
	if url == "" {
		url = e.Url
	}

	err := db.QueryRow(`SELECT label FROM example WHERE url = $1`, url).Scan(&label)
	switch {
	case err == sql.ErrNoRows:
		return db.Exec(`INSERT INTO example (url, label, created_at, updated_at) VALUES ($1, $2, $3, $4)`, url, e.Label, now, now)
	case err != nil:
		return nil, err
	default:
		if label != e.Label {
			return db.Exec(`UPDATE example SET label = $2, updated_at = $3 WHERE url = $1 `, url, e.Label, now)
		}
		return nil, nil
	}
}

func InsertExampleFromScanner(scanner *bufio.Scanner) (*example.Example, error) {
	line := scanner.Text()
	e, err := file.ParseLine(line)
	if err != nil {
		return nil, err
	}
	_, err = InsertOrUpdateExample(e)
	if err != nil {
		return nil, err
	}
	return e, nil
}

func InsertExamplesFromReader(r io.Reader) error {
	scanner := bufio.NewScanner(r)

	for scanner.Scan() {
		_, err := InsertExampleFromScanner(scanner)
		if err != nil {
			return err
		}
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	return nil
}

func readExamples(query string, args ...interface{}) ([]*example.Example, error) {
	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var examples example.Examples

	for rows.Next() {
		var label example.LabelType
		var url string
		if err := rows.Scan(&url, &label); err != nil {
			return nil, err
		}
		e := example.Example{Url: url, Label: label}
		examples = append(examples, &e)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return examples, nil
}

func ReadExamples() ([]*example.Example, error) {
	query := `SELECT url, label FROM example;`
	return readExamples(query)
}

func ReadRecentExamples(from time.Time) ([]*example.Example, error) {
	query := `SELECT url, label FROM example WHERE created_at > $1 ORDER BY updated_at DESC;`
	return readExamples(query, from)
}

func ReadExamplesByLabel(label example.LabelType, limit int) ([]*example.Example, error) {
	query := `SELECT url, label FROM example WHERE label = $1 ORDER BY updated_at DESC LIMIT $2;`
	return readExamples(query, label, limit)
}

func ReadLabeledExamples(limit int) ([]*example.Example, error) {
	query := `SELECT url, label FROM example WHERE label != 0 ORDER BY updated_at DESC LIMIT $1;`
	return readExamples(query, limit)
}

func ReadPositiveExamples(limit int) ([]*example.Example, error) {
	return ReadExamplesByLabel(example.POSITIVE, limit)
}

func ReadNegativeExamples(limit int) ([]*example.Example, error) {
	return ReadExamplesByLabel(example.NEGATIVE, limit)
}

func ReadUnlabeledExamples(limit int) ([]*example.Example, error) {
	return ReadExamplesByLabel(example.UNLABELED, limit)
}

func SearchExamplesByUlrs(urls []string) (example.Examples, error) {
	// ref: https://godoc.org/github.com/lib/pq#Array
	query := `SELECT url, label FROM example WHERE url = ANY($1);`
	return readExamples(query, pq.Array(urls))
}

func DeleteAllExamples() (sql.Result, error) {
	return db.Exec(`DELETE FROM example`)
}
