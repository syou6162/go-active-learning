package db

import (
	"bufio"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/syou6162/go-active-learning/lib/example"
	"github.com/syou6162/go-active-learning/lib/util"
	"os"
	"time"
)

func CreateDBConnection() (*sql.DB, error) {
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	return sql.Open("postgres", fmt.Sprintf("user=%s password=%s dbname=go-active-learning sslmode=disable", dbUser, dbPassword))
}

func CreateExampleTable(db *sql.DB) (sql.Result, error) {
	schema := `
CREATE TABLE IF NOT EXISTS example (
  "id" SERIAL,
  "url" TEXT NOT NULL,
  "label" INT NOT NULL,
  "created_at" timestamp NOT NULL,
  "updated_at" timestamp NOT NULL
);
CREATE UNIQUE INDEX IF NOT EXISTS "url_idx_example" ON example ("url");
`
	return db.Exec(schema)
}

func InsertExample(db *sql.DB, e *example.Example) (sql.Result, error) {
	now := time.Now()
	return db.Exec(`
INSERT INTO example (url, label, created_at, updated_at) VALUES ($1, $2, $3, $4)
`, e.Url, e.Label, now, now)
}

func InsertExampleFromScanner(db *sql.DB, scanner *bufio.Scanner) (*example.Example, error) {
	line := scanner.Text()
	e, err := util.ParseLine(line)
	if err != nil {
		return nil, err
	}
	_, err = InsertExample(db, e)
	if err != nil {
		return nil, err
	}
	return e, nil
}

func ReadExamples(db *sql.DB) ([]*example.Example, error) {
	rows, err := db.Query(`SELECT url, label FROM example`)
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
