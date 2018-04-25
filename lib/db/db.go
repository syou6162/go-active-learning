package db

import (
	"bufio"
	"database/sql"
	"fmt"
	"time"

	"io"

	_ "github.com/lib/pq"
	"github.com/syou6162/go-active-learning/lib/example"
	"github.com/syou6162/go-active-learning/lib/util"
	"github.com/syou6162/go-active-learning/lib/util/file"
)

func CreateDBConnection() (*sql.DB, error) {
	host := util.GetEnv("POSTGRES_HOST", "localhost")
	dbUser := util.GetEnv("DB_USER", "nobody")
	dbPassword := util.GetEnv("DB_PASSWORD", "nobody")
	dbName := util.GetEnv("DB_NAME", "go-active-learning")
	return sql.Open("postgres", fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable", host, dbUser, dbPassword, dbName))
}

func InsertOrUpdateExample(db *sql.DB, e *example.Example) (sql.Result, error) {
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

func InsertExampleFromScanner(db *sql.DB, scanner *bufio.Scanner) (*example.Example, error) {
	line := scanner.Text()
	e, err := file.ParseLine(line)
	if err != nil {
		return nil, err
	}
	_, err = InsertOrUpdateExample(db, e)
	if err != nil {
		return nil, err
	}
	return e, nil
}

func InsertExamplesFromReader(r io.Reader) error {
	scanner := bufio.NewScanner(r)

	conn, err := CreateDBConnection()
	if err != nil {
		return err
	}
	defer conn.Close()

	for scanner.Scan() {
		_, err := InsertExampleFromScanner(conn, scanner)
		if err != nil {
			return err
		}
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	return nil
}

func readExamples(db *sql.DB, query string) ([]*example.Example, error) {
	rows, err := db.Query(query)
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

func ReadExamples(db *sql.DB) ([]*example.Example, error) {
	query := `SELECT url, label FROM example;`
	return readExamples(db, query)
}

func ReadLabeledExamples(db *sql.DB) ([]*example.Example, error) {
	query := `SELECT url, label FROM example WHERE label != 0;`
	return readExamples(db, query)
}

func DeleteAllExamples(db *sql.DB) (sql.Result, error) {
	return db.Exec(`DELETE FROM example`)
}
