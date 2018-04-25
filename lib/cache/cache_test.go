package cache

import (
	"fmt"
	"testing"

	"github.com/syou6162/go-active-learning/lib/example"
)

func TestCacheGet(t *testing.T) {
	c, err := NewCache()
	if err != nil {
		t.Error("Cannot connect to redis")
	}
	defer c.Close()

	example := example.NewExample("http://a.hatena.ne.jp", example.POSITIVE)
	c.Client.Del("url:http://a.hatena.ne.jp")
	e, ok := c.GetExample(*example)
	if ok {
		t.Error(fmt.Printf("Cache must not contain %s", example.Url))
	}

	c.AddExample(*example)
	e, ok = c.GetExample(*example)
	if !ok {
		t.Error(fmt.Printf("Cache must return %s", example.Url))
	}
	if example.Url != e.Url {
		t.Error(fmt.Printf("Urls must be same(%s, %s)", example.Url, e.Url))
	}
}

func TestAttachMetaData(t *testing.T) {
	e1 := example.NewExample("http://b.hatena.ne.jp", example.POSITIVE)
	e2 := example.NewExample("http://www.yasuhisay.info", example.NEGATIVE)
	e3 := example.NewExample("https://github.com", example.UNLABELED)
	examples := example.Examples{e1, e2, e3}
	cache, err := NewCache()
	if err != nil {
		t.Error("Cannot connect to redis")
	}
	defer cache.Close()
	cache.AttachMetaData(examples)

	if examples[0].Title == "" {
		t.Errorf("Title must not be empty for %s", examples[0].Url)
	}
	if len(examples[0].Fv) == 0 {
		t.Errorf("Feature vector must not be empty for %s", examples[0].Url)
	}
}
