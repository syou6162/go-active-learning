package repository

import (
	"github.com/lib/pq"
	"github.com/syou6162/go-active-learning/lib/model"
)

func (r *repository) UpdateOrCreateReferringTweets(e *model.Example) error {
	if e.ReferringTweets == nil || len(*e.ReferringTweets) == 0 {
		return nil
	}

	tmp, err := r.FindExampleByUlr(e.Url)
	if err != nil {
		return err
	}
	id := tmp.Id

	for _, t := range *e.ReferringTweets {
		t.ExampleId = id
		if _, err = r.db.NamedExec(`
INSERT INTO tweet
( example_id,  created_at,  id_str,  full_text,  favorite_count,  retweet_count,  lang,  screen_name,  name,  profile_image_url,  label,  score)
VALUES
(:example_id, :created_at, :id_str, :full_text, :favorite_count, :retweet_count, :lang, :screen_name, :name, :profile_image_url, :label, :score)
ON CONFLICT (example_id, id_str)
DO UPDATE SET
favorite_count = :favorite_count,  retweet_count = :retweet_count, label = :label
WHERE
EXCLUDED.label != 0 AND tweet.label != EXCLUDED.label
;`, t); err != nil {
			return err
		}
	}
	return nil
}

func (r *repository) UpdateTweetLabel(exampleId int, idStr string, label model.LabelType) error {
	if _, err := r.db.Exec(`UPDATE tweet SET label = $1 WHERE example_id = $2 AND id_str = $3;`, label, exampleId, idStr); err != nil {
		return err
	}
	return nil
}

func (r *repository) SearchReferringTweetsList(examples model.Examples) (map[int]model.ReferringTweets, error) {
	referringTweetsByExampleId := make(map[int]model.ReferringTweets)

	referringTweets := model.ReferringTweets{}
	exampleIds := make([]int, 0)
	for _, e := range examples {
		exampleIds = append(exampleIds, e.Id)
	}

	query := `SELECT * FROM tweet WHERE example_id = ANY($1) AND label != -1 AND score > -1.0 AND (lang = 'en' OR lang = 'ja') ORDER BY favorite_count DESC;`
	err := r.db.Select(&referringTweets, query, pq.Array(exampleIds))
	if err != nil {
		return referringTweetsByExampleId, err
	}

	for _, t := range referringTweets {
		referringTweetsByExampleId[t.ExampleId] = append(referringTweetsByExampleId[t.ExampleId], t)
	}
	return referringTweetsByExampleId, nil
}

func (r *repository) SearchReferringTweets(limit int) (model.ReferringTweets, error) {
	referringTweets := model.ReferringTweets{}
	query := `SELECT * FROM tweet WHERE lang = 'en' OR lang = 'ja' ORDER BY created_at DESC LIMIT $1;`
	err := r.db.Select(&referringTweets, query, limit)
	if err != nil {
		return referringTweets, err
	}
	return referringTweets, nil
}

func (r *repository) searchReferringTweetsByLabel(label model.LabelType, limit int) (model.ReferringTweets, error) {
	referringTweets := model.ReferringTweets{}
	query := `
SELECT * FROM tweet WHERE id IN
  (SELECT id FROM
    (SELECT tweet.id, ROW_NUMBER() OVER(partition BY example_id ORDER BY favorite_count DESC) AS rank
    FROM tweet
    INNER JOIN example ON tweet.example_id = example.id
    WHERE tweet.label = $1 AND (lang = 'en' OR lang = 'ja') AND (example.label = 1 OR example.label = 0)
  ) AS t WHERE rank < 4)
ORDER BY created_at DESC LIMIT $2
;`
	err := r.db.Select(&referringTweets, query, label, limit)
	if err != nil {
		return referringTweets, err
	}
	return referringTweets, nil
}

func (r *repository) SearchPositiveReferringTweets(limit int) (model.ReferringTweets, error) {
	return r.searchReferringTweetsByLabel(model.POSITIVE, limit)
}

func (r *repository) SearchNegativeReferringTweets(limit int) (model.ReferringTweets, error) {
	return r.searchReferringTweetsByLabel(model.NEGATIVE, limit)
}

func (r *repository) SearchUnlabeledReferringTweets(limit int) (model.ReferringTweets, error) {
	return r.searchReferringTweetsByLabel(model.UNLABELED, limit)
}

func (r *repository) FindReferringTweets(e *model.Example) (model.ReferringTweets, error) {
	referringTweets := model.ReferringTweets{}

	query := `SELECT * FROM tweet WHERE example_id = $1 AND label != -1 AND score > 0.0 AND (lang = 'en' OR lang = 'ja') ORDER BY favorite_count DESC;`
	err := r.db.Select(&referringTweets, query, e.Id)
	if err != nil {
		return referringTweets, err
	}
	return referringTweets, nil
}
