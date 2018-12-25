package repository

import (
	"bufio"
	"database/sql"
	"time"

	"io"

	"fmt"

	"github.com/lib/pq"
	_ "github.com/lib/pq"
	"github.com/syou6162/go-active-learning/lib/feature"
	"github.com/syou6162/go-active-learning/lib/model"
	"github.com/syou6162/go-active-learning/lib/util/file"
)

var exampleNotFoundError = model.NotFoundError("example")

// データが存在しなければ追加
// データが存在する場合は、以下の場合にのみ更新する
// - ラベルが正例か負例に変更された
// - クロール対象のサイトが一時的に200以外のステータスで前回データが取得できなかった
func (r *repository) UpdateOrCreateExample(e *model.Example) error {
	now := time.Now()
	e.UpdatedAt = now
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
WHERE
((EXCLUDED.label != 0) AND (example.label != EXCLUDED.label)) OR
((example.status_code != 200) AND (EXCLUDED.status_code = 200))
;`, e)
	if err != nil {
		return err
	}
	tmp, err := r.FindExampleByUlr(e.Url)
	if err != nil {
		return err
	}
	e.Id = tmp.Id
	return nil
}

func (r *repository) UpdateScore(e *model.Example) error {
	if _, err := r.FindExampleByUlr(e.Url); err != nil {
		return err
	}
	if _, err := r.db.Exec(`UPDATE example SET score = $1, updated_at = $2 WHERE url = $3;`, e.Score, time.Now(), e.Url); err != nil {
		return err
	}
	return nil
}

func (r *repository) IncErrorCount(e *model.Example) error {
	errorCount, err := r.GetErrorCount(e)
	if err != nil {
		return err
	}
	if _, err := r.db.Exec(`UPDATE example SET error_count = $1, updated_at = $2 WHERE url = $3;`, errorCount+1, time.Now(), e.Url); err != nil {
		return err
	}
	return nil
}

func (r *repository) GetErrorCount(e *model.Example) (int, error) {
	example, err := r.FindExampleByUlr(e.Url)
	if err != nil {
		return 0, err
	}
	return example.ErrorCount, nil
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
	err = r.UpdateOrCreateExample(e)
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

func (r *repository) searchExamples(query string, args ...interface{}) (model.Examples, error) {
	examples := model.Examples{}
	err := r.db.Select(&examples, query, args...)
	if err != nil {
		return nil, err
	}
	return examples, nil
}

func (r *repository) findExample(query string, args ...interface{}) (*model.Example, error) {
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

func (r *repository) SearchExamples() (model.Examples, error) {
	query := `SELECT * FROM example;`
	return r.searchExamples(query)
}

func (r *repository) SearchRecentExamples(from time.Time, limit int) (model.Examples, error) {
	query := `SELECT * FROM example WHERE created_at > $1 ORDER BY updated_at DESC LIMIT $2;`
	return r.searchExamples(query, from, limit)
}

func (r *repository) SearchRecentExamplesByHost(host string, from time.Time, limit int) (model.Examples, error) {
	query := `SELECT * FROM example WHERE final_url like $1 || '%' AND created_at > $2 ORDER BY updated_at DESC LIMIT $3;`
	return r.searchExamples(query, host, from, limit)
}

func (r *repository) SearchExamplesByLabel(label model.LabelType, limit int) (model.Examples, error) {
	query := `SELECT * FROM example WHERE label = $1 ORDER BY updated_at DESC LIMIT $2;`
	return r.searchExamples(query, label, limit)
}

func (r *repository) SearchLabeledExamples(limit int) (model.Examples, error) {
	query := `SELECT * FROM example WHERE label != 0 ORDER BY updated_at DESC LIMIT $1;`
	return r.searchExamples(query, limit)
}

func (r *repository) SearchPositiveExamples(limit int) (model.Examples, error) {
	return r.SearchExamplesByLabel(model.POSITIVE, limit)
}

func (r *repository) SearchNegativeExamples(limit int) (model.Examples, error) {
	return r.SearchExamplesByLabel(model.NEGATIVE, limit)
}

func (r *repository) SearchUnlabeledExamples(limit int) (model.Examples, error) {
	return r.SearchExamplesByLabel(model.UNLABELED, limit)
}

func (r *repository) SearchPositiveScoredExamples(limit int) (model.Examples, error) {
	query := `SELECT * FROM example WHERE score > 0 ORDER BY updated_at DESC LIMIT $1;`
	return r.searchExamples(query, limit)
}

func (r *repository) FindExampleByUlr(url string) (*model.Example, error) {
	query := `SELECT * FROM example WHERE url = $1;`
	return r.findExample(query, url)
}

func (r *repository) SearchExamplesByUlrs(urls []string) (model.Examples, error) {
	// ref: https://godoc.org/github.com/lib/pq#Array
	query := `SELECT * FROM example WHERE url = ANY($1);`
	return r.searchExamples(query, pq.Array(urls))
}

func (r *repository) SearchExamplesByIds(ids []int) (model.Examples, error) {
	query := `SELECT * FROM example WHERE id = ANY($1);`
	return r.searchExamples(query, pq.Array(ids))
}

func (r *repository) SearchExamplesByKeywords(keywords []string, aggregator string, limit int) (model.Examples, error) {
	if len(keywords) == 0 {
		return model.Examples{}, nil
	}
	regexList := make([]string, 0)
	for _, w := range keywords {
		regexList = append(regexList, fmt.Sprintf(`.*%s.*`, w))
	}
	query := fmt.Sprintf(`SELECT * FROM example WHERE title ~* %s($1) AND label != -1 ORDER BY (label, score) DESC LIMIT $2;`, aggregator)
	return r.searchExamples(query, pq.Array(regexList), limit)
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

func (r *repository) SearchFeatureVector(examples model.Examples) (map[int]feature.FeatureVector, error) {
	type Pair struct {
		ExampleId int    `db:"example_id"`
		Feature   string `db:"feature"`
	}

	fvById := make(map[int]feature.FeatureVector)
	urls := make([]string, 0)
	for _, e := range examples {
		urls = append(urls, e.Url)
	}

	tmp, err := r.SearchExamplesByUlrs(urls)
	if err != nil {
		return fvById, err
	}
	ids := make([]int, 0)
	for _, e := range tmp {
		ids = append(ids, e.Id)
	}

	query := `SELECT example_id, feature FROM feature WHERE example_id = ANY($1);`
	pairs := make([]Pair, 0)
	err = r.db.Select(&pairs, query, pq.Array(ids))
	if err != nil {
		return fvById, err
	}

	for _, pair := range pairs {
		fvById[pair.ExampleId] = append(fvById[pair.ExampleId], pair.Feature)
	}
	return fvById, nil
}

func (r *repository) DeleteAllExamples() error {
	_, err := r.db.Exec(`DELETE FROM example`)
	return err
}
