package repository

import (
	"github.com/lib/pq"
	"github.com/syou6162/go-active-learning/lib/model"
)

func (r *repository) UpdateRecommendation(rec model.Recommendation) error {
	if _, err := r.db.Exec(`DELETE FROM recommendation WHERE list_type = $1;`, rec.RecommendationListType); err != nil {
		return err
	}
	if _, err := r.db.Exec(`INSERT INTO recommendation (list_type, example_id) VALUES ($1, unnest(cast($2 AS INT[])));`, rec.RecommendationListType, pq.Array(rec.ExampleIds)); err != nil {
		return err
	}
	return nil
}

func (r *repository) FindRecommendation(t model.RecommendationListType) (*model.Recommendation, error) {
	rec := &model.Recommendation{RecommendationListType: t}
	items := make([]int, 0)
	query := `SELECT example_id FROM recommendation WHERE list_type = $1;`
	err := r.db.Select(&items, query, t)
	if err != nil {
		return nil, err
	}
	rec.ExampleIds = items
	return rec, nil
}
