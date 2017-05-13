package main

import (
	"encoding/gob"
	"os"
)

type Cache struct {
	cache map[string]string
}

func NewCache() *Cache {
	return &Cache{make(map[string]string)}
}

func (c *Cache) Add(example Example) {
	c.cache[example.url] = example.title
}

func (c *Cache) Save(filename string) error {
	file, err := os.Create(filename)
	defer file.Close()
	if err != nil {
		return err
	}

	enc := gob.NewEncoder(file)
	enc.Encode(&c.cache)
	return nil
}

func LoadCache(filename string) (*Cache, error) {
	cache := NewCache()
	file, err := os.Open(filename)
	defer file.Close()
	if err != nil {
		return cache, err
	}

	decoder := gob.NewDecoder(file)
	c := make(map[string]string)
	decoder.Decode(&c)
	cache.cache = c
	return cache, nil
}
