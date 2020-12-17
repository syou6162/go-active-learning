package repository_test

import (
	"testing"
	"time"

	"github.com/syou6162/go-active-learning/lib/example"
	"github.com/syou6162/go-active-learning/lib/model"
	"github.com/syou6162/go-active-learning/lib/repository"
)

func TestUpdateReferringTweets(t *testing.T) {
	repo, err := repository.New()
	if err != nil {
		t.Errorf(err.Error())
	}
	defer repo.Close()

	if err = repo.DeleteAllExamples(); err != nil {
		t.Error(err)
	}

	e := example.NewExample("http://hoge.com", model.UNLABELED)
	err = repo.UpdateOrCreateExample(e)
	if err != nil {
		t.Error(err)
	}
	now := time.Now()
	idStr := "1111111"
	t1 := model.Tweet{
		CreatedAt:       now,
		IdStr:           idStr,
		FullText:        "hello world!!!",
		FavoriteCount:   10,
		RetweetCount:    10,
		Lang:            "en",
		ScreenName:      "syou6162",
		Name:            "syou6162",
		ProfileImageUrl: "http://hogehoge.com/profile.png",
		Score:           1.0,
	}

	tweets := model.ReferringTweets{}
	tweets.Tweets = append(tweets.Tweets, &t1)
	tweets.Count = len(tweets.Tweets)
	e.ReferringTweets = &tweets
	if err = repo.UpdateOrCreateReferringTweets(e); err != nil {
		t.Error(err)
	}

	{
		result, err := repo.SearchReferringTweetsList(model.Examples{e}, 10)
		if err != nil {
			t.Error(err)
		}
		if len(result) == 0 {
			t.Error("result must not be empty")
		}
		if len(result[e.Id].Tweets) == 0 {
			t.Error("result must not be empty")
		}
		if result[e.Id].Count == 0 {
			t.Error("result must not be zero")
		}
		if result[e.Id].Tweets[0].Name != "syou6162" {
			t.Error("Name must be syou6162")
		}
	}

	{
		result, err := repo.FindReferringTweets(e, 10)
		if err != nil {
			t.Error(err)
		}
		if len(result.Tweets) == 0 {
			t.Error("result must not be empty")
		}
		if result.Count == 0 {
			t.Error("result must not be empty")
		}
		if result.Tweets[0].Name != "syou6162" {
			t.Error("Name must be syou6162")
		}
	}

	{
		result, err := repo.FindReferringTweets(e, 0)
		if err != nil {
			t.Error(err)
		}
		if len(result.Tweets) != 0 {
			t.Error("result must be empty")
		}
		if result.Count == 0 {
			t.Error("result must not be empty")
		}
	}

	{
		if err := repo.UpdateTweetLabel(e.Id, idStr, model.NEGATIVE); err != nil {
			t.Error(err)
		}
		result, err := repo.FindReferringTweets(e, 10)
		if err != nil {
			t.Error(err)
		}
		if len(result.Tweets) != 0 {
			t.Error("result must be empty")
		}
		if result.Count != 1 {
			t.Error("result must be 1")
		}
	}
}

func TestSearchReferringTweetsByLabel(t *testing.T) {
	repo, err := repository.New()
	if err != nil {
		t.Errorf(err.Error())
	}
	defer repo.Close()

	if err = repo.DeleteAllExamples(); err != nil {
		t.Error(err)
	}

	e := example.NewExample("http://hoge.com", model.UNLABELED)
	err = repo.UpdateOrCreateExample(e)
	if err != nil {
		t.Error(err)
	}
	now := time.Now()
	idStr := "1111111"
	t1 := model.Tweet{
		CreatedAt:       now,
		IdStr:           idStr,
		FullText:        "hello world!!!",
		FavoriteCount:   10,
		RetweetCount:    10,
		Lang:            "en",
		ScreenName:      "syou6162",
		Name:            "syou6162",
		ProfileImageUrl: "http://hogehoge.com/profile.png",
		Label:           model.POSITIVE,
	}

	tweets := model.ReferringTweets{}
	tweets.Tweets = append(tweets.Tweets, &t1)
	tweets.Count = len(tweets.Tweets)
	e.ReferringTweets = &tweets
	if err = repo.UpdateOrCreateReferringTweets(e); err != nil {
		t.Error(err)
	}

	limit := 10
	{
		result, err := repo.SearchPositiveReferringTweets(3, -1.0, limit)
		if err != nil {
			t.Error(err)
		}
		if len(result.Tweets) != 1 {
			t.Error("len(result) must be 1")
		}
		if result.Count != 1 {
			t.Error("Count must be 1")
		}
	}
	{
		result, err := repo.SearchNegativeReferringTweets(3, -1.0, limit)
		if err != nil {
			t.Error(err)
		}
		if len(result.Tweets) != 0 {
			t.Error("len(result) must be empty")
		}
		if result.Count != 0 {
			t.Error("Count must be zero")
		}
	}
}

func TestSearchRecentReferringTweetsWithHighScore(t *testing.T) {
	repo, err := repository.New()
	if err != nil {
		t.Errorf(err.Error())
	}
	defer repo.Close()

	if err = repo.DeleteAllExamples(); err != nil {
		t.Error(err)
	}

	e := example.NewExample("http://hoge.com", model.UNLABELED)
	err = repo.UpdateOrCreateExample(e)
	if err != nil {
		t.Error(err)
	}
	now := time.Now()
	t1 := model.Tweet{
		CreatedAt:       now,
		IdStr:           "1111111",
		FullText:        "hello world!!!",
		FavoriteCount:   10,
		RetweetCount:    10,
		Lang:            "en",
		ScreenName:      "syou6162",
		Name:            "syou6162",
		ProfileImageUrl: "http://hogehoge.com/profile.png",
		Label:           model.POSITIVE,
		Score:           10.0,
	}
	t2 := model.Tweet{
		CreatedAt:       now,
		IdStr:           "22222222",
		FullText:        "hello world!!!",
		FavoriteCount:   10,
		RetweetCount:    10,
		Lang:            "en",
		ScreenName:      "syou6162",
		Name:            "syou6162",
		ProfileImageUrl: "http://hogehoge.com/profile.png",
		Label:           model.POSITIVE,
		Score:           10.0,
	}
	t3 := model.Tweet{
		CreatedAt:       now,
		IdStr:           "3333333333",
		FullText:        "hello world!!!",
		FavoriteCount:   10,
		RetweetCount:    10,
		Lang:            "en",
		ScreenName:      "syou6162",
		Name:            "syou6162",
		ProfileImageUrl: "http://hogehoge.com/profile.png",
		Label:           model.POSITIVE,
		Score:           -10.0,
	}

	tweets := model.ReferringTweets{}
	tweets.Tweets = append(tweets.Tweets, &t1, &t2, &t3)
	tweets.Count = len(tweets.Tweets)
	e.ReferringTweets = &tweets
	if err = repo.UpdateOrCreateReferringTweets(e); err != nil {
		t.Error(err)
	}

	limit := 10
	{
		result, err := repo.SearchRecentReferringTweetsWithHighScore(now.Add(time.Duration(-10*24)*time.Hour), 0.0, limit)
		if err != nil {
			t.Error(err)
		}
		if len(result.Tweets) != 2 {
			t.Error("len(result) must be 2")
		}
		if result.Count != 2 {
			t.Error("Count must be 2")
		}
	}
}
