package cache

import (
	"encoding/json"
	"fmt"

	"math"
	"os"
	"runtime"
	"strconv"
	"sync"
	"time"

	"github.com/go-redis/redis"
	"github.com/syou6162/go-active-learning/lib/example"
	"github.com/syou6162/go-active-learning/lib/fetcher"
	"github.com/syou6162/go-active-learning/lib/util"
)

type Cache struct {
	Client *redis.Client
}

var redisPrefix = "url"

func NewCache() (*Cache, error) {
	host := util.GetEnv("REDIS_HOST", "localhost")
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:6379", host),
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	_, err := client.Ping().Result()
	if err != nil {
		return nil, err
	}
	return &Cache{Client: client}, nil
}

func (c *Cache) Close() error {
	return c.Client.Close()
}

// ToDo: return (Example, error)
func (c *Cache) GetExample(exa example.Example) (example.Example, bool) {
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

func (c *Cache) AddExample(example example.Example) error {
	key := redisPrefix + ":" + example.Url
	json, _ := json.Marshal(example)
	if err := c.Client.Set(key, json, 0).Err(); err != nil {
		return err
	}
	if err := c.Client.Expire(key, time.Hour*240).Err(); err != nil {
		return err
	}
	return nil
}

var listPrefix = "list:"

func (c *Cache) AddExamplesToList(listName string, examples example.Examples) error {
	if err := c.Client.Del(listPrefix + listName).Err(); err != nil {
		return err
	}

	result := make([]redis.Z, 0)
	for _, e := range examples {
		url := e.Url
		if e.FinalUrl != "" {
			url = e.FinalUrl
		}
		result = append(result, redis.Z{Score: e.Score, Member: url})
	}
	err := c.Client.ZAdd(listPrefix+listName, result...).Err()
	if err != nil {
		return err
	}
	return nil
}

func (c *Cache) GetUrlsFromList(listName string, from int64, to int64) ([]string, error) {
	sliceCmd := c.Client.ZRevRange(listPrefix+listName, from, to)
	if sliceCmd.Err() != nil {
		return nil, sliceCmd.Err()
	}
	return sliceCmd.Val(), nil
}

func (cache *Cache) attachMetaData(examples example.Examples) {
	oldStdout := os.Stdout
	readFile, writeFile, _ := os.Pipe()
	os.Stdout = writeFile

	defer func() {
		writeFile.Close()
		readFile.Close()
		os.Stdout = oldStdout
	}()

	wg := &sync.WaitGroup{}
	cpus := runtime.NumCPU()
	runtime.GOMAXPROCS(cpus)
	sem := make(chan struct{}, 4)
	for idx, e := range examples {
		wg.Add(1)
		sem <- struct{}{}
		go func(e *example.Example, idx int) {
			defer wg.Done()
			if tmp, ok := cache.GetExample(*e); ok {
				e.Title = tmp.Title
				e.FinalUrl = tmp.FinalUrl
				e.Description = tmp.Description
				e.OgDescription = tmp.OgDescription
				e.Body = tmp.Body
				e.CleanedText = tmp.CleanedText
				e.StatusCode = tmp.StatusCode
				e.Fv = tmp.Fv
				e.Score = tmp.Score
			} else {
				fmt.Fprintln(os.Stderr, "Fetching("+strconv.Itoa(idx)+"): "+e.Url)
				article := fetcher.GetArticle(e.Url)
				e.Title = article.Title
				e.FinalUrl = article.Url
				e.Description = article.Description
				e.OgDescription = article.OgDescription
				e.Body = article.Body
				e.CleanedText = article.CleanedText
				e.StatusCode = article.StatusCode
				e.Fv = util.RemoveDuplicate(example.ExtractFeatures(*e))
				e.Body = ""
				cache.AddExample(*e)
			}
			<-sem
		}(e, idx)
	}
	wg.Wait()
}

func (cache *Cache) AttachMetaData(examples example.Examples) {
	batchSize := 100
	examplesList := make([]example.Examples, 0)
	n := len(examples)

	for i := 0; i < n; i += batchSize {
		max := int(math.Min(float64(i+batchSize), float64(n)))
		examplesList = append(examplesList, examples[i:max])
	}
	for _, l := range examplesList {
		cache.attachMetaData(l)
	}
}
