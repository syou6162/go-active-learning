package repository

import (
	"bufio"
	"database/sql"
	"time"

	"io"

	"github.com/lib/pq"
	_ "github.com/lib/pq"
	"github.com/syou6162/go-active-learning/lib/feature"
	"github.com/syou6162/go-active-learning/lib/model"
	"github.com/syou6162/go-active-learning/lib/util/file"
)

var exampleNotFoundError = model.NotFoundError("example")

func (r *repository) InsertOrUpdateExample(e *model.Example) error {
	_, err := r.db.NamedExec(`
INSERT INTO example
( url,  final_url,  title,  description,  og_description,  og_type,  og_image,  body,  score,  is_new,  status_code,  favicon,  label,  created_at,  updated_at)
VALUES
(:url, :final_url, :title, :description, :og_description, :og_type, :og_image, :body, :score, :is_new, :status_code, :favicon, :label, :created_at, :updated_at)
ON CONFLICT (url)
DO UPDATE SET
url = :url, final_url = :final_url, title = :title,
description = :description, og_description = :og_description, og_type = :og_type, og_image = :og_image,
body = :body, score = :score, is_new = :is_new, status_code = :status_code, favicon = :favicon,
label = :label, created_at = :created_at, updated_at = :updated_at
WHERE (:label != 0) AND (example.label != EXCLUDED.label)
;`, e)
	return err
}

func (r *repository) UpdateFeatureVector(e *model.Example) error {
	tmp, err := r.FindExampleByUlr(e.Url)
	if err != nil {
		return err
	}
	id := tmp.Id
	if _, err = r.db.Exec(`DELETE FROM feature WHERE example_id = $1;`, id); err != nil {
		return err
	}
	_, err = r.db.Exec(`INSERT INTO feature (example_id, feature) VALUES ($1, unnest(cast($2 AS TEXT[])));`, id, pq.Array(e.Fv))
	return err
}

func (r *repository) InsertExampleFromScanner(scanner *bufio.Scanner) (*model.Example, error) {
	line := scanner.Text()
	e, err := file.ParseLine(line)
	if err != nil {
		return nil, err
	}
	err = r.InsertOrUpdateExample(e)
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
	examples := model.Examples{}
	err := r.db.Select(&examples, query, args...)
	if err != nil {
		return nil, err
	}
	return examples, nil
}

func (r *repository) readExample(query string, args ...interface{}) (*model.Example, error) {
	e := model.Example{}

	err := r.db.Get(&e, query, args...)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, exampleNotFoundError
		}
		return nil, err
	}
	return &e, nil
}

func (r *repository) ReadExamples() (model.Examples, error) {
	query := `SELECT * FROM example;`
	return r.readExamples(query)
}

func (r *repository) ReadRecentExamples(from time.Time) (model.Examples, error) {
	query := `SELECT * FROM example WHERE created_at > $1 ORDER BY updated_at DESC;`
	return r.readExamples(query, from)
}

func (r *repository) ReadExamplesByLabel(label model.LabelType, limit int) (model.Examples, error) {
	query := `SELECT * FROM example WHERE label = $1 ORDER BY updated_at DESC LIMIT $2;`
	return r.readExamples(query, label, limit)
}

func (r *repository) ReadLabeledExamples(limit int) (model.Examples, error) {
	query := `SELECT * FROM example WHERE label != 0 ORDER BY updated_at DESC LIMIT $1;`
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

func (r *repository) FindExampleByUlr(url string) (*model.Example, error) {
	query := `SELECT * FROM example WHERE url = $1;`
	return r.readExample(query, url)
}

func (r *repository) SearchExamplesByUlrs(urls []string) (model.Examples, error) {
	// ref: https://godoc.org/github.com/lib/pq#Array
	query := `SELECT * FROM example WHERE url = ANY($1);`
	return r.readExamples(query, pq.Array(urls))
}

func (r *repository) FindFeatureVector(e *model.Example) (feature.FeatureVector, error) {
	fv := feature.FeatureVector{}
	tmp, err := r.FindExampleByUlr(e.Url)
	if err != nil {
		return fv, err
	}
	id := tmp.Id
	query := `SELECT feature FROM feature WHERE example_id = $1;`
	err = r.db.Select(&fv, query, id)
	if err != nil {
		return fv, err
	}
	return fv, nil
}

func (r *repository) DeleteAllExamples() error {
	_, err := r.db.Exec(`DELETE FROM example`)
	return err
}
