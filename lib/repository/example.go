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

var exampleNotFoundError = model.NotFoundError("example")

func (r *repository) InsertOrUpdateExample(e *model.Example) error {
	var label model.LabelType

	url := e.FinalUrl
	if url == "" {
		url = e.Url
	}

	err := r.db.QueryRow(`SELECT label FROM example WHERE url = $1`, url).Scan(&label)
	switch {
	case err == sql.ErrNoRows:
		_, err = r.db.Exec(`INSERT INTO example (url, final_url, title, description, og_description, og_type, og_image, body, score, is_new, status_code, favicon, label, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)`,
			url, e.FinalUrl, e.Title, e.Description, e.OgDescription, e.OgType, e.OgImage, e.Body, e.Score, e.IsNew, e.StatusCode, e.Favicon, e.Label, e.CreatedAt, e.UpdatedAt)
		return err
	case err != nil:
		return err
	default:
		if label != e.Label && // ラベルが変更される
			e.Label != model.UNLABELED { // 変更されるラベルはPOSITIVEかNEGATIVEのみ
			_, err = r.db.Exec(`UPDATE example SET label = $2, updated_at = $3, title = $4, description = $5, og_description = $6, og_type = $7, og_image = $8, body = $9, score = $10, is_new = $11, status_code = $12, favicon = $13 WHERE url = $1 `,
				url, e.Label, e.UpdatedAt, e.Title, e.Description, e.OgDescription, e.OgType, e.OgImage, e.Body, e.Score, e.IsNew, e.StatusCode, e.Favicon)
			return err
		}
		return nil
	}
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

func (r *repository) SearchExamplesByUlr(url string) (*model.Example, error) {
	query := `SELECT * FROM example WHERE url = $1;`
	return r.readExample(query, url)
}

func (r *repository) SearchExamplesByUlrs(urls []string) (model.Examples, error) {
	// ref: https://godoc.org/github.com/lib/pq#Array
	query := `SELECT * FROM example WHERE url = ANY($1);`
	return r.readExamples(query, pq.Array(urls))
}

func (r *repository) DeleteAllExamples() error {
	_, err := r.db.Exec(`DELETE FROM example`)
	return err
}
