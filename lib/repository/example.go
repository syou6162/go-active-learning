package repository

import (
	"bufio"
	"database/sql"
	"time"

	"io"

	"github.com/lib/pq"
	_ "github.com/lib/pq"
	"github.com/syou6162/go-active-learning/lib/model"
	"github.com/syou6162/go-active-learning/lib/util/file"
)

func (r *repository) InsertOrUpdateExample(e *model.Example) (sql.Result, error) {
	var label model.LabelType

	url := e.FinalUrl
	if url == "" {
		url = e.Url
	}

	err := r.db.QueryRow(`SELECT label FROM example WHERE url = $1`, url).Scan(&label)
	switch {
	case err == sql.ErrNoRows:
		return r.db.Exec(`INSERT INTO example (url, label, created_at, updated_at) VALUES ($1, $2, $3, $4)`, url, e.Label, e.CreatedAt, e.UpdatedAt)
	case err != nil:
		return nil, err
	default:
		if label != e.Label && // ラベルが変更される
			e.Label != model.UNLABELED { // 変更されるラベルはPOSITIVEかNEGATIVEのみ
			return r.db.Exec(`UPDATE example SET label = $2, updated_at = $3 WHERE url = $1 `, url, e.Label, e.UpdatedAt)
		}
		return nil, nil
	}
}

func (r *repository) InsertExampleFromScanner(scanner *bufio.Scanner) (*model.Example, error) {
	line := scanner.Text()
	e, err := file.ParseLine(line)
	if err != nil {
		return nil, err
	}
	_, err = r.InsertOrUpdateExample(e)
	if err != nil {
		return nil, err
	}
	return e, nil
}

func (r *repository) InsertExamplesFromReader(reader io.Reader) error {
	scanner := bufio.NewScanner(reader)

	for scanner.Scan() {
		_, err := r.InsertExampleFromScanner(scanner)
		if err != nil {
			return err
		}
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	return nil
}

func (r *repository) readExamples(query string, args ...interface{}) (model.Examples, error) {
	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var examples model.Examples

	for rows.Next() {
		var label model.LabelType
		var url string
		var createdAt time.Time
		var updatedAt time.Time
		if err := rows.Scan(&url, &label, &createdAt, &updatedAt); err != nil {
			return nil, err
		}
		e := model.Example{Url: url, Label: label, CreatedAt: createdAt, UpdatedAt: updatedAt}
		examples = append(examples, &e)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return examples, nil
}

func (r *repository) readExample(query string, args ...interface{}) (*model.Example, error) {
	var label model.LabelType
	var url string
	var createdAt time.Time
	var updatedAt time.Time

	row := r.db.QueryRow(query, args...)
	if err := row.Scan(&url, &label, &createdAt, &updatedAt); err != nil {
		return nil, err
	}
	e := model.Example{Url: url, Label: label, CreatedAt: createdAt, UpdatedAt: updatedAt}
	return &e, nil
}

func (r *repository) ReadExamples() (model.Examples, error) {
	query := `SELECT url, label, created_at, updated_at FROM example;`
	return r.readExamples(query)
}

func (r *repository) ReadRecentExamples(from time.Time) (model.Examples, error) {
	query := `SELECT url, label, created_at, updated_at FROM example WHERE created_at > $1 ORDER BY updated_at DESC;`
	return r.readExamples(query, from)
}

func (r *repository) ReadExamplesByLabel(label model.LabelType, limit int) (model.Examples, error) {
	query := `SELECT url, label, created_at, updated_at FROM example WHERE label = $1 ORDER BY updated_at DESC LIMIT $2;`
	return r.readExamples(query, label, limit)
}

func (r *repository) ReadLabeledExamples(limit int) (model.Examples, error) {
	query := `SELECT url, label, created_at, updated_at FROM example WHERE label != 0 ORDER BY updated_at DESC LIMIT $1;`
	return r.readExamples(query, limit)
}

func (r *repository) ReadPositiveExamples(limit int) (model.Examples, error) {
	return r.ReadExamplesByLabel(model.POSITIVE, limit)
}

func (r *repository) ReadNegativeExamples(limit int) (model.Examples, error) {
	return r.ReadExamplesByLabel(model.NEGATIVE, limit)
}

func (r *repository) ReadUnlabeledExamples(limit int) (model.Examples, error) {
	return r.ReadExamplesByLabel(model.UNLABELED, limit)
}

func (r *repository) SearchExamplesByUlr(url string) (*model.Example, error) {
	query := `SELECT url, label, created_at, updated_at FROM example WHERE url = $1;`
	return r.readExample(query, url)
}

func (r *repository) SearchExamplesByUlrs(urls []string) (model.Examples, error) {
	// ref: https://godoc.org/github.com/lib/pq#Array
	query := `SELECT url, label, created_at, updated_at FROM example WHERE url = ANY($1);`
	return r.readExamples(query, pq.Array(urls))
}

func (r *repository) DeleteAllExamples() (sql.Result, error) {
	return r.db.Exec(`DELETE FROM example`)
}
