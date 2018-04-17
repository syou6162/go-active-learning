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

func getEnv(key, fallback string) string {
	value, ok := os.LookupEnv(key)
	if !ok {
		value = fallback
	}
	return value
}

func CreateDBConnection() (*sql.DB, error) {
	dbUser := getEnv("DB_USER", "nobody")
	dbPassword := getEnv("DB_PASSWORD", "nobody")
	dbName := getEnv("DB_NAME", "go-active-learning-test")
	return sql.Open("postgres", fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", dbUser, dbPassword, dbName))
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
