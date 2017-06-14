package main

import (
	"fmt"
	"testing"
)

func TestCacheGet(t *testing.T) {
	c := NewCache()
	example := NewExample("http://b.hatena.ne.jp", POSITIVE)
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

func TestCacheSave(t *testing.T) {
	c := NewCache()
	example := NewExample("http://b.hatena.ne.jp", POSITIVE)
	c.Add(*example)
	c.Get(*example)
	err := c.Save(CacheFilename)

	if err != nil {
		t.Error(fmt.Printf("Error (%s) occurs when saving cache file", err))
	}
}

func TestLoadCache(t *testing.T) {
	c := NewCache()
	example := NewExample("http://b.hatena.ne.jp", POSITIVE)
	c.Add(*example)
	c.Get(*example)
	c.Save(CacheFilename)

	c, err := LoadCache(CacheFilename)
	if err != nil {
		t.Error(fmt.Printf("Error (%s) occurs when loading cache file", err))
	}

	e, ok := c.Get(*example)
	if !ok {
		t.Error(fmt.Printf("Cache must return %s", example.Url))
	}
	if example.Url != e.Url {
		t.Error(fmt.Printf("Urls must be same(%s, %s)", example.Url, e.Url))
	}
}
