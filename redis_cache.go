package main

import (
	"encoding/json"
	"github.com/go-redis/redis"
)

type RedisCache struct {
	Client *redis.Client
}

var redisPrefix = "url"

func NewRedisCache() *RedisCache {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	return &RedisCache{Client: client}
}

// ToDo: return (Example, error)
func (c *RedisCache) Get(example Example) (Example, bool) {
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
func (c *RedisCache) Add(example Example) {
	key := redisPrefix + ":" + example.Url
	json, _ := json.Marshal(example)
	c.Client.Set(key, json, 0).Err()
}
