package cache

import (
	"testing"

	"github.com/syou6162/go-active-learning/lib/example"
)

func TestAttachMetaData(t *testing.T) {
	e1 := example.NewExample("http://b.hatena.ne.jp", example.POSITIVE)
	e2 := example.NewExample("http://www.yasuhisay.info", example.NEGATIVE)
	e3 := example.NewExample("https://github.com", example.UNLABELED)
	examples := example.Examples{e1, e2, e3}
	err := Init()
	if err != nil {
		t.Error("Cannot connect to redis")
	}
	defer Close()
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

	AttachMetadata(examples, true, false)
	if examples[0].Title == "" {
		t.Errorf("Title must not be empty for %s", examples[0].Url)
	}
	if len(examples[0].Fv) == 0 {
		t.Errorf("Feature vector must not be empty for %s", examples[0].Url)
	}

	if examples[1].OgType == "" {
		t.Errorf("OgType must not be empty for %s", examples[1].Url)
	}
}
