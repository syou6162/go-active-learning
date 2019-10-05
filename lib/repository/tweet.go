package repository

import (
	"time"

	"github.com/lib/pq"
	"github.com/syou6162/go-active-learning/lib/model"
)

func (r *repository) UpdateOrCreateReferringTweets(e *model.Example) error {
	if e.ReferringTweets == nil || len((*e).ReferringTweets.Tweets) == 0 || (*e).ReferringTweets.Count == 0 {
		return nil
	}

	tmp, err := r.FindExampleByUlr(e.Url)
	if err != nil {
		return err
	}
	id := tmp.Id

	for _, t := range (*e).ReferringTweets.Tweets {
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

type exampleIdWithTweetsCount struct {
	ExampleId   int `db:"example_id"`
	TweetsCount int `db:"tweets_count"`
}

func (r *repository) SearchReferringTweetsList(examples model.Examples, limitForEachExample int) (map[int]model.ReferringTweets, error) {
	referringTweetsByExampleId := make(map[int]model.ReferringTweets)

	exampleIds := make([]int, 0)
	for _, e := range examples {
		exampleIds = append(exampleIds, e.Id)
	}

	exampleIdsWithTweetsCount := make([]exampleIdWithTweetsCount, 0)
	tweetsCountByExampleQuery := `SELECT example_id, COUNT(*) AS tweets_count FROM tweet WHERE example_id = ANY($1) GROUP BY example_id ORDER BY tweets_count DESC;`
	err := r.db.Select(&exampleIdsWithTweetsCount, tweetsCountByExampleQuery, pq.Array(exampleIds))
	if err != nil {
		return referringTweetsByExampleId, err
	}
	tweetsCountByExampleId := make(map[int]int)
	for _, e := range exampleIdsWithTweetsCount {
		tweetsCountByExampleId[e.ExampleId] = e.TweetsCount
	}

	if limitForEachExample == 0 {
		for _, exampleId := range exampleIds {
			referringTweets := model.ReferringTweets{Count: 0, Tweets: make([]*model.Tweet, 0)}
			if cnt, ok := tweetsCountByExampleId[exampleId]; ok {
				referringTweets.Count = cnt
			}
			referringTweetsByExampleId[exampleId] = referringTweets
		}
		return referringTweetsByExampleId, nil
	}

	tweets := make([]*model.Tweet, 0)
	query := `SELECT * FROM tweet WHERE example_id = ANY($1) AND label != -1 AND score > -1.0 AND (lang = 'en' OR lang = 'ja') ORDER BY favorite_count DESC LIMIT $2;`
	err = r.db.Select(&tweets, query, pq.Array(exampleIds), limitForEachExample)
	if err != nil {
		return referringTweetsByExampleId, err
	}
	tweetsByExampleId := make(map[int][]*model.Tweet)
	for _, t := range tweets {
		tweetsByExampleId[t.ExampleId] = append(tweetsByExampleId[t.ExampleId], t)
	}

	for _, exampleId := range exampleIds {
		referringTweets := model.ReferringTweets{Count: 0, Tweets: make([]*model.Tweet, 0)}
		if tweets, ok := tweetsByExampleId[exampleId]; ok {
			referringTweets.Tweets = tweets
		}
		if cnt, ok := tweetsCountByExampleId[exampleId]; ok {
			referringTweets.Count = cnt
		}
		referringTweetsByExampleId[exampleId] = referringTweets
	}
	return referringTweetsByExampleId, nil
}

func (r *repository) SearchReferringTweets(limit int) (model.ReferringTweets, error) {
	referringTweets := model.ReferringTweets{Count: 0, Tweets: make([]*model.Tweet, 0)}
	query := `SELECT * FROM tweet WHERE lang = 'en' OR lang = 'ja' ORDER BY created_at DESC LIMIT $1;`
	err := r.db.Select(&referringTweets.Tweets, query, limit)
	if err != nil {
		return referringTweets, err
	}
	referringTweets.Count = len(referringTweets.Tweets)
	return referringTweets, nil
}

func (r *repository) SearchRecentReferringTweetsWithHighScore(from time.Time, scoreThreshold float64, limit int) (model.ReferringTweets, error) {
	referringTweets := model.ReferringTweets{Count: 0, Tweets: make([]*model.Tweet, 0)}
	query := `
SELECT 
	tweet.id,
	tweet.example_id,

	tweet.created_at,
	tweet.id_str,
	tweet.full_text,
	tweet.favorite_count,
	tweet.retweet_count,
	tweet.lang,

	tweet.screen_name,
	tweet.name,
	tweet.profile_image_url,
	tweet.label,
	tweet.score
FROM 
	tweet 
INNER JOIN 
	example ON example.id = example_id 
WHERE
	tweet.created_at > $1 AND 
	tweet.label != -1 AND 
	example.label != -1 AND 
	tweet.score > $2 AND 
	(favorite_count > 0 OR retweet_count > 0) AND
	(lang = 'en' OR lang = 'ja')
ORDER BY tweet.score DESC
LIMIT $3
;
`
	err := r.db.Select(&referringTweets.Tweets, query, from, scoreThreshold, limit)
	if err != nil {
		return referringTweets, err
	}
	referringTweets.Count = len(referringTweets.Tweets)
	return referringTweets, nil
}

func (r *repository) searchReferringTweetsByLabel(label model.LabelType, limit int) (model.ReferringTweets, error) {
	referringTweets := model.ReferringTweets{Count: 0, Tweets: make([]*model.Tweet, 0)}
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
	err := r.db.Select(&referringTweets.Tweets, query, label, limit)
	if err != nil {
		return referringTweets, err
	}
	referringTweets.Count = len(referringTweets.Tweets)
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

type tweetsCount struct {
	Count int `db:"count"`
}

func (r *repository) FindReferringTweets(e *model.Example, limit int) (model.ReferringTweets, error) {
	referringTweets := model.ReferringTweets{Count: 0, Tweets: make([]*model.Tweet, 0)}

	countQuery := `SELECT COUNT(*) AS count FROM tweet WHERE example_id = $1;`
	cnt := tweetsCount{}
	err := r.db.Get(&cnt, countQuery, e.Id)
	if err != nil {
		return referringTweets, err
	}
	referringTweets.Count = cnt.Count
	if limit == 0 {
		return referringTweets, err
	}

	query := `SELECT * FROM tweet WHERE example_id = $1 AND label != -1 AND score > 0.0 AND (lang = 'en' OR lang = 'ja') ORDER BY favorite_count DESC LIMIT $2;`
	err = r.db.Select(&referringTweets.Tweets, query, e.Id, limit)
	if err != nil {
		return referringTweets, err
	}
	return referringTweets, nil
}
