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
	err = repo.InsertOrUpdateExample(e)
	if err != nil {
		t.Error(err)
	}
	now := time.Now()
	t1 := model.Tweet{
		CreatedAt: now,
		IdStr: "1111111",
		FullText: "hello world!!!",
		FavoriteCount: 10,
		RetweetCount: 10,
		Lang: "en",
		ScreenName: "syou6162",
		Name: "syou6162",
		ProfileImageUrl: "http://hogehoge.com/profile.png",
	}

	tweets := model.ReferringTweets{&t1}
	e.ReferringTweets = &tweets
	if err = repo.UpdateReferringTweets(e); err != nil {
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
}
