package repository

func (r *repository) UpdateTopAccessedExampleIds(exampleIds []int) error {
	if _, err := r.db.Exec(`DELETE FROM top_accessed_example;`); err != nil {
		return err
	}
	if _, err := r.db.Exec(`INSERT INTO top_accessed_example (example_id) VALUES (unnest(cast($1 AS INT[])));`, exampleIds); err != nil {
		return err
	}
	return nil
}

func (r *repository) SearchTopAccessedExampleIds() ([]int, error) {
	exampleIds := make([]int, 0)
	query := `SELECT example_id FROM top_accessed_example;`
	err := r.db.Select(&exampleIds, query)
	if err != nil {
		return nil, err
	}
	return exampleIds, nil
}
