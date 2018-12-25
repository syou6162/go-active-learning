package cache

import (
	"fmt"

	"time"

	"strconv"

	"github.com/go-redis/redis"
	"github.com/syou6162/go-active-learning/lib/util"
)

type Cache interface {
	Ping() error
	Close() error

	IncErrorCount(url string) error
	GetErrorCount(url string) (int, error)
}

type cache struct {
	client *redis.Client
}

func New() (*cache, error) {
	host := util.GetEnv("REDIS_HOST", "localhost")
	client := redis.NewClient(&redis.Options{
		Addr:        fmt.Sprintf("%s:6379", host),
		Password:    "", // no password set
		DB:          0,  // use default DB
		PoolSize:    100,
		PoolTimeout: time.Duration(5) * time.Second,
		IdleTimeout: time.Duration(10) * time.Second,
	})
	if _, err := client.Ping().Result(); err != nil {
		return nil, err
	}
	return &cache{client: client}, nil
}

func (c *cache) Ping() error {
	return c.client.Ping().Err()
}

func (c *cache) Close() error {
	if c.client != nil {
		return c.client.Close()
	} else {
		return nil
	}
}

var errorCountPrefix = "errorCountPrefix:"

func (c *cache) IncErrorCount(url string) error {
	key := errorCountPrefix + url
	exist, err := c.client.Exists(key).Result()
	if err != nil {
		return err
	}
	if exist == 0 {
		hour := 24 * 10
		c.client.Set(key, 1, time.Hour*time.Duration(hour))
		return nil
	} else {
		if _, err = c.client.Incr(key).Result(); err != nil {
			return err
		}
	}
	return nil
}

func (c *cache) GetErrorCount(url string) (int, error) {
	key := errorCountPrefix + url
	ok, err := c.client.Exists(key).Result()
	if err != nil {
		return 0, err
	}
	if ok == 0 {
		return 0, nil
	}

	cntStr, err := c.client.Get(key).Result()
	if err != nil {
		return 0, err
	}
	cnt, err := strconv.Atoi(cntStr)
	if err != nil {
		return 0, err
	}
	return cnt, nil
}
