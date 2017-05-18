package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sync"
)

type Cache struct {
	Cache map[string]Example `json:"Cache"`
	sync.RWMutex
}

var CacheFilename = "cache.bin"

func NewCache() *Cache {
	return &Cache{Cache: make(map[string]Example)}
}

func (c *Cache) Get(example Example) (Example, bool) {
	c.RLock()
	defer c.RUnlock()

	e, ok := c.Cache[example.Url]
	return e, ok
}

func (c *Cache) Add(example Example) {
	c.Lock()
	defer c.Unlock()

	c.Cache[example.Url] = example
}

func (c *Cache) Save(filename string) error {
	fmt.Fprintln(os.Stderr, "Saving cache...")
	file, err := os.Create(filename)
	defer file.Close()
	if err != nil {
		return err
	}
	json, err := json.Marshal(c.Cache)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(filename, json, 0644)
	if err != nil {
		return err
	}
	return nil
}

func LoadCache(filename string) (*Cache, error) {
	fmt.Fprintln(os.Stderr, "Loading cache...")
	cache := NewCache()
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return cache, err
	}

	if err := json.Unmarshal(bytes, &cache.Cache); err != nil {
		return cache, err
	}
	return cache, nil
}
