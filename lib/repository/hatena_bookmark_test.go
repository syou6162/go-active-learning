package repository_test

import (
	"testing"
	"time"

	"github.com/syou6162/go-active-learning/lib/example"
	"github.com/syou6162/go-active-learning/lib/model"
	"github.com/syou6162/go-active-learning/lib/repository"
)

func TestUpdateHatenaBookmark(t *testing.T) {
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
	b1 := model.Bookmark{
		User:      "syou6162",
		Comment:   "面白いサイトですね",
		Timestamp: model.HatenaBookmarkTime{Time: &now},
		Tags:      model.Tags{"hack"},
	}
	hb := model.HatenaBookmark{
		ExampleId: e.Id,
		Title:     "hoge",
		Count:     10,
		Bookmarks: []*model.Bookmark{&b1},
	}
	e.HatenaBookmark = &hb
	if err = repo.UpdateHatenaBookmark(e); err != nil {
		t.Error(err)
	}

	{
		result, err := repo.SearchHatenaBookmarks(model.Examples{e})
		if err != nil {
			t.Error(err)
		}

		for _, tmp := range result {
			if tmp.Title == "" {
				t.Error("Title must not be empty")
			}
			for _, b := range tmp.Bookmarks {
				if b.User == "" {
					t.Error("User must not be empty")
				}
				if len(b.Tags) == 0 {
					t.Error("Tags must not be empty")
				}
			}
		}
	}

	{
		result, err := repo.FindHatenaBookmark(e)
		if err != nil {
			t.Error(err)
		}

		if result.Title == "" {
			t.Error("Title must not be empty")
		}
		for _, b := range result.Bookmarks {
			if b.User == "" {
				t.Error("User must not be empty")
			}
			if len(b.Tags) == 0 {
				t.Error("Tags must not be empty")
			}
		}
	}
}
