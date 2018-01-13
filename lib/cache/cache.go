package cache

import (
	"encoding/json"
	"github.com/go-redis/redis"
	"github.com/syou6162/go-active-learning/lib/example"
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
func (c *Cache) Get(exa example.Example) (example.Example, bool) {
	key := redisPrefix + ":" + exa.Url
	exampleStr, err := c.Client.Get(key).Result()
	e := example.Example{}
	if err != nil {
		return e, false
	}
	if err := json.Unmarshal([]byte(exampleStr), &e); err != nil {
		return e, false
	}
	return e, true
}

// ToDo: return error...
func (c *Cache) Add(example example.Example) {
	key := redisPrefix + ":" + example.Url
	json, _ := json.Marshal(example)
	c.Client.Set(key, json, 0).Err()
}
