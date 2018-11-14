package cache

import (
	"log"
	"os"
	"testing"

	"github.com/syou6162/go-active-learning/lib/example"
)

func TestMain(m *testing.M) {
	err := Init()
	if err != nil {
		log.Fatal(err.Error())
	}
	defer Close()

	ret := m.Run()
	os.Exit(ret)
}

func TestPing(t *testing.T) {
	if err := Ping(); err != nil {
		t.Errorf(err.Error())
	}
}

func TestErrorCount(t *testing.T) {
	existingUrl := "https://github.com"
	nonExistingUrl := "http://hoge.fuga"
	urls := []string{existingUrl, nonExistingUrl}
	for _, u := range urls {
		key := errorCountPrefix + u
		client.Del(key)
	}

	for _, u := range urls {
		cnt, err := getErrorCount(u)
		if err != nil {
			t.Errorf("Cannot get error count: %s", err.Error())
		}
		if cnt != 0 {
			t.Errorf("Error count must be 0 for %s", u)
		}
	}

	for _, u := range urls {
		err := incErrorCount(u)
		if err != nil {
			t.Errorf("Cannot get error count: %s", err.Error())
		}
	}

	for _, u := range urls {
		cnt, err := getErrorCount(u)
		if err != nil {
			t.Errorf("Cannot get error count: %s", err.Error())
		}
		if cnt != 1 {
			t.Errorf("Error count must be 1 for %s", u)
		}
	}
}

func TestAttachMetaData(t *testing.T) {
	e1 := example.NewExample("http://b.hatena.ne.jp", example.POSITIVE)
	e2 := example.NewExample("https://www.yasuhisay.info", example.NEGATIVE)
	e3 := example.NewExample("https://github.com", example.UNLABELED)
	examples := example.Examples{e1, e2, e3}
	for _, e := range examples {
		key := "url:" + e.Url
		client.Del(key)
		client.Del(errorCountPrefix + e.Url)
	}

	if examples[0].Title != "" {
		t.Errorf("Title must be empty for %s", examples[0].Url)
	}
	if len(examples[0].Fv) != 0 {
		t.Errorf("Feature vector must be empty for %s", examples[0].Url)
	}
	AttachMetadata(examples, false, false)

	if examples[0].Title != "" {
		t.Errorf("Title must be empty for %s", examples[0].Url)
	}
	if len(examples[0].Fv) != 0 {
		t.Errorf("Feature vector must be empty for %s", examples[0].Url)
	}

	if examples[1].OgType != "" {
		t.Errorf("OgType must be empty for %s", examples[1].Url)
	}

	AttachMetadata(examples, true, false)
	if examples[0].Title == "" {
		t.Errorf("Title must not be empty for %s", examples[0].Url)
	}
	if len(examples[0].Fv) == 0 {
		t.Errorf("Feature vector must not be empty for %s", examples[0].Url)
	}

	if examples[1].OgType != "blog" {
		t.Errorf("OgType must be blog for %s", examples[1].Url)
	}

	e4 := example.NewExample("http://b.hatena.ne.jp", example.POSITIVE)
	e5 := example.NewExample("https://www.yasuhisay.info", example.NEGATIVE)
	e6 := example.NewExample("https://github.com", example.UNLABELED)
	examples = example.Examples{e4, e5, e6}
	AttachMetadata(examples, false, false)

	if examples[0].Title == "" {
		t.Errorf("Title must be empty for %s", examples[0].Url)
	}
	if len(examples[0].Fv) == 0 {
		t.Errorf("Feature vector must not be empty for %s", examples[0].Url)
	}

	if examples[1].OgType != "blog" {
		t.Errorf("OgType must be blog for %s", examples[1].Url)
	}

	cnt, err := getErrorCount("https://github.com")
	if err != nil {
		t.Errorf("Cannot get error count: %s", err.Error())
	}
	if cnt != 0 {
		t.Errorf("count should be 0, but %d", cnt)
	}
}

func TestAttachMetaDataNonExistingUrls(t *testing.T) {
	nonExistingUrl := "http://hoge.fuga"
	e := example.NewExample(nonExistingUrl, example.UNLABELED)
	examples := example.Examples{e}
	for _, e := range examples {
		key := "url:" + e.Url
		client.Del(key)
		client.Del(errorCountPrefix + e.Url)
	}

	for i := 1; i <= 3; i++ {
		AttachMetadata(examples, true, false)
		if examples[0].Title != "" {
			t.Errorf("Title must not be empty for %s", examples[0].Url)
		}
		cnt, err := getErrorCount(nonExistingUrl)
		if err != nil {
			t.Errorf("Cannot get error count: %s", err.Error())
		}
		if i != cnt {
			t.Errorf("Count should be %d, but %d", i, cnt)
		}
	}
}

func TestAttachLightMetaData(t *testing.T) {
	e1 := example.NewExample("http://b.hatena.ne.jp", example.POSITIVE)
	e2 := example.NewExample("https://www.yasuhisay.info", example.NEGATIVE)
	e3 := example.NewExample("https://github.com", example.UNLABELED)
	examples := example.Examples{e1, e2, e3}
	for _, e := range examples {
		key := "url:" + e.Url
		client.Del(key)
	}

	if examples[0].Title != "" {
		t.Errorf("Title must be empty for %s", examples[0].Url)
	}
	if len(examples[0].Fv) != 0 {
		t.Errorf("Feature vector must be empty for %s", examples[0].Url)
	}
	AttachMetadata(examples, false, false)

	if examples[0].Title != "" {
		t.Errorf("Title must be empty for %s", examples[0].Url)
	}
	if len(examples[0].Fv) != 0 {
		t.Errorf("Feature vector must be empty for %s", examples[0].Url)
	}

	if examples[1].OgType != "" {
		t.Errorf("OgType must be empty for %s", examples[1].Url)
	}

	AttachMetadata(examples, true, true)

	e1 = example.NewExample("http://b.hatena.ne.jp", example.POSITIVE)
	e2 = example.NewExample("https://www.yasuhisay.info", example.NEGATIVE)
	e3 = example.NewExample("https://github.com", example.UNLABELED)
	examples = example.Examples{e1, e2, e3}

	AttachMetadata(examples, false, true)

	if examples[0].Title == "" {
		t.Errorf("Title must not be empty for %s", examples[0].Url)
	}
	if len(examples[0].Fv) != 0 {
		t.Errorf("Feature vector must be empty for %s", examples[0].Url)
	}

	if examples[1].OgType != "blog" {
		t.Errorf("OgType must be blog for %s", examples[1].Url)
	}
}

func TestReferringTweets(t *testing.T) {
	e1 := example.NewExample("http://b.hatena.ne.jp", example.POSITIVE)
	examples := example.Examples{e1}
	for _, e := range examples {
		key := "url:" + e.Url
		client.Del(key)
	}

	AttachMetadata(examples, true, true)
	e1.ReferringTweets = example.ReferringTweets{"https:/twitter.com/1"}
	SetExample(*e1)
	e2 := example.NewExample("http://b.hatena.ne.jp", example.POSITIVE)
	examples = example.Examples{e2}
	AttachMetadata(examples, false, true)

	if len(examples[0].ReferringTweets) != 1 {
		t.Errorf("len(examples[0].ReferringTweets) should be 1, but %d", len(examples[0].ReferringTweets))
	}
}

func TestAddExamplesToListAndGetUrlsFromList(t *testing.T) {
	listName := "general"
	client.Del("list:" + listName)
	err := AddExamplesToList(listName, example.Examples{})
	if err == nil {
		t.Error("Error should occur when adding empty list")
	}

	e1 := example.NewExample("http://b.hatena.ne.jp", example.POSITIVE)
	e2 := example.NewExample("https://www.yasuhisay.info", example.NEGATIVE)
	e3 := example.NewExample("https://github.com", example.UNLABELED)
	examples := example.Examples{e1, e2, e3}
	for _, e := range examples {
		key := "url:" + e.Url
		client.Del(key)
	}
	AttachMetadata(examples, true, false)

	err = AddExamplesToList(listName, examples)
	if err != nil {
		t.Error(err.Error())
	}

	list, err := GetUrlsFromList(listName, 0, 100)
	if err != nil {
		t.Error(err.Error())
	}
	if len(list) != 3 {
		t.Errorf("len(list) == %d, want 3", len(list))
	}
}
