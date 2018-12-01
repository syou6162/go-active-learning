package cache

import (
	"testing"

	"github.com/syou6162/go-active-learning/lib/example"
	"github.com/syou6162/go-active-learning/lib/model"
)

func TestPing(t *testing.T) {
	cache, err := New()
	if err != nil {
		t.Errorf(err.Error())
	}
	defer cache.Close()

	if err := cache.Ping(); err != nil {
		t.Errorf(err.Error())
	}
}

func TestErrorCount(t *testing.T) {
	cache, err := New()
	if err != nil {
		t.Errorf(err.Error())
	}
	defer cache.Close()

	existingUrl := "https://github.com"
	nonExistingUrl := "http://hoge.fuga"
	urls := []string{existingUrl, nonExistingUrl}
	for _, u := range urls {
		key := errorCountPrefix + u
		cache.client.Del(key)
	}

	for _, u := range urls {
		cnt, err := cache.GetErrorCount(u)
		if err != nil {
			t.Errorf("Cannot get error count: %s", err.Error())
		}
		if cnt != 0 {
			t.Errorf("Error count must be 0 for %s", u)
		}
	}

	for _, u := range urls {
		err := cache.IncErrorCount(u)
		if err != nil {
			t.Errorf("Cannot get error count: %s", err.Error())
		}
	}

	for _, u := range urls {
		cnt, err := cache.GetErrorCount(u)
		if err != nil {
			t.Errorf("Cannot get error count: %s", err.Error())
		}
		if cnt != 1 {
			t.Errorf("Error count must be 1 for %s", u)
		}
	}
}

func TestAddExamplesToListAndGetUrlsFromList(t *testing.T) {
	cache, err := New()
	if err != nil {
		t.Errorf(err.Error())
	}
	defer cache.Close()

	listName := "general"
	cache.client.Del("list:" + listName)
	err = cache.AddExamplesToList(listName, model.Examples{})
	if err == nil {
		t.Error("Error should occur when adding empty list")
	}

	e1 := example.NewExample("http://b.hatena.ne.jp", model.POSITIVE)
	e2 := example.NewExample("https://www.yasuhisay.info", model.NEGATIVE)
	e3 := example.NewExample("https://github.com", model.UNLABELED)
	examples := model.Examples{e1, e2, e3}

	err = cache.AddExamplesToList(listName, examples)
	if err != nil {
		t.Error(err.Error())
	}

	list, err := cache.GetUrlsFromList(listName, 0, 100)
	if err != nil {
		t.Error(err.Error())
	}
	if len(list) != 3 {
		t.Errorf("len(list) == %d, want 3", len(list))
	}
}
