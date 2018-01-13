package main

import (
	"fmt"
	"testing"
)

func TestCacheGet(t *testing.T) {
	c, err := NewCache()
	if err != nil {
		t.Error("Cannot connect to redis")
	}

	example := NewExample("http://a.hatena.ne.jp", POSITIVE)
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
