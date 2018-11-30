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
		cnt, err := cache.getErrorCount(u)
		if err != nil {
			t.Errorf("Cannot get error count: %s", err.Error())
		}
		if cnt != 0 {
			t.Errorf("Error count must be 0 for %s", u)
		}
	}

	for _, u := range urls {
		err := cache.incErrorCount(u)
		if err != nil {
			t.Errorf("Cannot get error count: %s", err.Error())
		}
	}

	for _, u := range urls {
		cnt, err := cache.getErrorCount(u)
		if err != nil {
			t.Errorf("Cannot get error count: %s", err.Error())
		}
		if cnt != 1 {
			t.Errorf("Error count must be 1 for %s", u)
		}
	}
}

func TestAttachMetaData(t *testing.T) {
	cache, err := New()
	if err != nil {
		t.Errorf(err.Error())
	}
	defer cache.Close()

	e1 := example.NewExample("http://b.hatena.ne.jp", model.POSITIVE)
	e2 := example.NewExample("https://www.yasuhisay.info", model.NEGATIVE)
	e3 := example.NewExample("https://github.com", model.UNLABELED)
	examples := model.Examples{e1, e2, e3}
	for _, e := range examples {
		key := "url:" + e.Url
		cache.client.Del(key)
		cache.client.Del(errorCountPrefix + e.Url)
	}

	if examples[0].Title != "" {
		t.Errorf("Title must be empty for %s", examples[0].Url)
	}
	if len(examples[0].Fv) != 0 {
		t.Errorf("Feature vector must be empty for %s", examples[0].Url)
	}
	cache.AttachMetadata(examples)

	if examples[0].Title != "" {
		t.Errorf("Title must be empty for %s", examples[0].Url)
	}
	if len(examples[0].Fv) != 0 {
		t.Errorf("Feature vector must be empty for %s", examples[0].Url)
	}

	if examples[1].OgType != "" {
		t.Errorf("OgType must be empty for %s", examples[1].Url)
	}

	cache.Fetch(examples)
	cache.UpdateExamplesMetadata(examples)
	if examples[0].Title == "" {
		t.Errorf("Title must not be empty for %s", examples[0].Url)
	}
	if len(examples[0].Fv) == 0 {
		t.Errorf("Feature vector must not be empty for %s", examples[0].Url)
	}

	if examples[1].OgType != "blog" {
		t.Errorf("OgType must be blog for %s", examples[1].Url)
	}

	e4 := example.NewExample("http://b.hatena.ne.jp", model.POSITIVE)
	e5 := example.NewExample("https://www.yasuhisay.info", model.NEGATIVE)
	e6 := example.NewExample("https://github.com", model.UNLABELED)
	examples = model.Examples{e4, e5, e6}
	cache.AttachMetadata(examples)

	if examples[0].Title == "" {
		t.Errorf("Title must be empty for %s", examples[0].Url)
	}
	if len(examples[0].Fv) == 0 {
		t.Errorf("Feature vector must not be empty for %s", examples[0].Url)
	}

	if examples[1].OgType != "blog" {
		t.Errorf("OgType must be blog for %s", examples[1].Url)
	}

	cnt, err := cache.getErrorCount("https://github.com")
	if err != nil {
		t.Errorf("Cannot get error count: %s", err.Error())
	}
	if cnt != 0 {
		t.Errorf("count should be 0, but %d", cnt)
	}
}

func TestAttachMetaDataNonExistingUrls(t *testing.T) {
	cache, err := New()
	if err != nil {
		t.Errorf(err.Error())
	}
	defer cache.Close()

	nonExistingUrl := "http://hoge.fuga"
	e := example.NewExample(nonExistingUrl, model.UNLABELED)
	examples := model.Examples{e}
	for _, e := range examples {
		key := "url:" + e.Url
		cache.client.Del(key)
		cache.client.Del(errorCountPrefix + e.Url)
	}

	for i := 1; i <= 3; i++ {
		cache.Fetch(examples)
		cache.AttachMetadata(examples)
		if examples[0].Title != "" {
			t.Errorf("Title must not be empty for %s", examples[0].Url)
		}
		cnt, err := cache.getErrorCount(nonExistingUrl)
		if err != nil {
			t.Errorf("Cannot get error count: %s", err.Error())
		}
		if i != cnt {
			t.Errorf("Count should be %d, but %d", i, cnt)
		}
	}
}

func TestAttachLightMetaData(t *testing.T) {
	cache, err := New()
	if err != nil {
		t.Errorf(err.Error())
	}
	defer cache.Close()

	e1 := example.NewExample("http://b.hatena.ne.jp", model.POSITIVE)
	e2 := example.NewExample("https://www.yasuhisay.info", model.NEGATIVE)
	e3 := example.NewExample("https://github.com", model.UNLABELED)
	examples := model.Examples{e1, e2, e3}
	for _, e := range examples {
		key := "url:" + e.Url
		cache.client.Del(key)
	}

	if examples[0].Title != "" {
		t.Errorf("Title must be empty for %s", examples[0].Url)
	}
	if len(examples[0].Fv) != 0 {
		t.Errorf("Feature vector must be empty for %s", examples[0].Url)
	}
	cache.AttachMetadata(examples)

	if examples[0].Title != "" {
		t.Errorf("Title must be empty for %s", examples[0].Url)
	}
	if len(examples[0].Fv) != 0 {
		t.Errorf("Feature vector must be empty for %s", examples[0].Url)
	}

	if examples[1].OgType != "" {
		t.Errorf("OgType must be empty for %s", examples[1].Url)
	}

	cache.Fetch(examples)
	if err := cache.UpdateExamplesMetadata(examples); err != nil {
		t.Error(err.Error())
	}

	e1 = example.NewExample("http://b.hatena.ne.jp", model.POSITIVE)
	e2 = example.NewExample("https://www.yasuhisay.info", model.NEGATIVE)
	e3 = example.NewExample("https://github.com", model.UNLABELED)
	examples = model.Examples{e1, e2, e3}

	cache.AttachLightMetadata(examples)

	if examples[0].Title == "" {
		println(examples[0].Title)
		t.Errorf("Title must not be empty for %s", examples[0].Url)
	}
	if len(examples[0].Fv) != 0 {
		t.Errorf("Feature vector must be empty for %s", examples[0].Url)
	}

	if examples[1].OgType != "blog" {
		println(examples[1].OgType)
		t.Errorf("OgType must be blog for %s", examples[1].Url)
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
	for _, e := range examples {
		key := "url:" + e.Url
		cache.client.Del(key)
	}
	cache.Fetch(examples)
	cache.AttachMetadata(examples)

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
