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
	query := `INSERT INTO model (model_type, model, c, accuracy, precision, recall, fvalue) VALUES ($1, $2, $3, $4, $5, $6, $7);`
	if _, err := r.db.Exec(query, m.ModelType, string(bytes), m.C, m.Accuracy, m.Precision, m.Recall, m.Fvalue); err != nil {
		return err
	}
	return nil
}

func (r *repository) FindLatestMIRAModel(modelType classifier.ModelType) (*classifier.MIRAClassifier, error) {
	type Classifier struct {
		Model string
	}
	tmp := Classifier{}

	query := `SELECT model FROM model WHERE model_type = $1 ORDER BY created_at DESC LIMIT 1;`
	err := r.db.Get(&tmp, query, modelType)
	if err != nil {
		return nil, err
	}

	clf := classifier.MIRAClassifier{}
	if err := json.Unmarshal(([]byte)(tmp.Model), &clf); err != nil {
		return nil, err
	}
	return &clf, nil
}
