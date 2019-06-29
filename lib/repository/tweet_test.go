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

	tweets := model.ReferringTweets{&t1}
	e.ReferringTweets = &tweets
	if err = repo.UpdateOrCreateReferringTweets(e); err != nil {
		t.Error(err)
	}

	{
		result, err := repo.SearchReferringTweetsList(model.Examples{e})
		if err != nil {
			t.Error(err)
		}
		if len(result) == 0 {
			t.Error("result must not be empty")
		}
		if len(result[e.Id]) == 0 {
			t.Error("result must not be empty")
		}
		if result[e.Id][0].Name != "syou6162" {
			t.Error("Name must be syou6162")
		}
	}

	{
		result, err := repo.FindReferringTweets(e)
		if err != nil {
			t.Error(err)
		}
		if len(result) == 0 {
			t.Error("result must not be empty")
		}
		if result[0].Name != "syou6162" {
			t.Error("Name must be syou6162")
		}
	}

	{
		if err := repo.UpdateTweetLabel(e.Id, idStr, model.NEGATIVE); err != nil {
			t.Error(err)
		}
		result, err := repo.FindReferringTweets(e)
		if err != nil {
			t.Error(err)
		}
		if len(result) != 0 {
			t.Error("result must be empty")
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

	tweets := model.ReferringTweets{&t1}
	e.ReferringTweets = &tweets
	if err = repo.UpdateOrCreateReferringTweets(e); err != nil {
		t.Error(err)
	}

	limit := 10
	{
		result, err := repo.SearchPositiveReferringTweets(limit)
		if err != nil {
			t.Error(err)
		}
		if len(result) != 1 {
			t.Error("len(result) must be 1")
		}
	}
	{
		result, err := repo.SearchNegativeReferringTweets(limit)
		if err != nil {
			t.Error(err)
		}
		if len(result) != 0 {
			t.Error("len(result) must be empty")
		}
	}
}
