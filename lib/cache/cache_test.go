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

	example := example.NewExample("http://a.hatena.ne.jp", example.POSITIVE)
	c.Client.Del("url:http://a.hatena.ne.jp")
	e, ok := c.Get(*example)
	if ok {
		t.Error(fmt.Printf("Cache must not contain %s", example.Url))
	}

	c.Add(*example)
	e, ok = c.Get(*example)
	if !ok {
		t.Error(fmt.Printf("Cache must return %s", example.Url))
	}
	if example.Url != e.Url {
		t.Error(fmt.Printf("Urls must be same(%s, %s)", example.Url, e.Url))
	}
}
