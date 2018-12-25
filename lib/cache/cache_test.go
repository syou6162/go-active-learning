package cache

import (
	"testing"
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
