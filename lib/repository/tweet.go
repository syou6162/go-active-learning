package repository

import (
	"github.com/syou6162/go-active-learning/lib/model"
	"github.com/lib/pq"
)

func (r *repository) UpdateReferringTweets(e *model.Example) error {
	if len(*e.ReferringTweets) == 0 {
		return nil
	}

	tmp, err := r.FindExampleByUlr(e.Url)
	if err != nil {
		return err
	}
	id := tmp.Id
	if _, err = r.db.Exec(`DELETE FROM tweet WHERE example_id = $1;`, id); err != nil {
		return err
	}

	for _, t := range *e.ReferringTweets {
		t.ExampleId = id
		if _, err = r.db.NamedExec(`
INSERT INTO tweet
( example_id,  created_at,  id_str,  full_text,  favorite_count,  retweet_count,  lang,  screen_name,  name,  profile_image_url)
VALUES
(:example_id, :created_at, :id_str, :full_text, :favorite_count, :retweet_count, :lang, :screen_name, :name, :profile_image_url)
;`, t); err != nil {
			return err
		}
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

	query := `SELECT * FROM tweet WHERE example_id = ANY($1);`
	err := r.db.Select(&referringTweets, query, pq.Array(exampleIds))
	if err != nil {
		return referringTweetsByExampleId, err
	}

	for _, t := range referringTweets {
		tmp := referringTweetsByExampleId[t.ExampleId]
		tmp = append(tmp, t)
	}
	return referringTweetsByExampleId, nil
}
