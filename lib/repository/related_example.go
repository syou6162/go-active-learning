package repository

import (
	"github.com/lib/pq"
	"github.com/syou6162/go-active-learning/lib/model"
)

func (r *repository) UpdateRelatedExamples(related model.RelatedExamples) error {
	if _, err := r.db.Exec(`DELETE FROM related_example WHERE example_id = $1;`, related.ExampleId); err != nil {
		return err
	}
	if _, err := r.db.Exec(`INSERT INTO related_example (example_id, related_example_id) VALUES ($1, unnest(cast($2 AS INT[])));`, related.ExampleId, pq.Array(related.RelatedExampleIds)); err != nil {
		return err
	}
	return nil
}

func (r *repository) FindRelatedExamples(e *model.Example) (*model.RelatedExamples, error) {
	related := &model.RelatedExamples{ExampleId: e.Id}
	items := make([]int, 0)
	query := `SELECT related_example_id FROM related_example WHERE example_id = $1;`
	err := r.db.Select(&items, query, e.Id)
	if err != nil {
		return nil, err
	}
	related.RelatedExampleIds = items
	return related, nil
}
