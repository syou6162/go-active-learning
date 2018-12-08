package repository

import (
	"encoding/json"

	"github.com/syou6162/go-active-learning/lib/classifier"
)

func (r *repository) InsertMIRAModel(m classifier.MIRAClassifier) error {
	bytes, err := json.Marshal(m)
	if err != nil {
		return err
	}
	if _, err := r.db.Exec(`INSERT INTO model (model) VALUES ($1);`, string(bytes)); err != nil {
		return err
	}
	return nil
}

func (r *repository) FindLatestMIRAModel() (*classifier.MIRAClassifier, error) {
	type Classifier struct {
		Model string
	}
	tmp := Classifier{}

	query := `SELECT model FROM model ORDER BY created_at DESC LIMIT 1;`
	err := r.db.Get(&tmp, query)
	if err != nil {
		return nil, err
	}

	clf := classifier.MIRAClassifier{}
	if err := json.Unmarshal(([]byte)(tmp.Model), &clf); err != nil {
		return nil, err
	}
	return &clf, nil
}
