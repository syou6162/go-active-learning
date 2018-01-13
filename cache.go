package main

import (
	"encoding/json"
	"github.com/go-redis/redis"
)

type Cache struct {
	Client *redis.Client
}

var redisPrefix = "url"

func NewCache() (*Cache, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	_, err := client.Ping().Result()
	if err != nil {
		return nil, err
	}
	return &Cache{Client: client}, nil
}

// ToDo: return (Example, error)
func (c *Cache) Get(example Example) (Example, bool) {
	key := redisPrefix + ":" + example.Url
	exampleStr, err := c.Client.Get(key).Result()
	e := Example{}
	if err != nil {
		return e, false
	}
	if err := json.Unmarshal([]byte(exampleStr), &e); err != nil {
		return e, false
	}
	return e, true
}

// ToDo: return error...
func (c *Cache) Add(example Example) {
	key := redisPrefix + ":" + example.Url
	json, _ := json.Marshal(example)
	c.Client.Set(key, json, 0).Err()
}
