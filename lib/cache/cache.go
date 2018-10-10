package cache

import (
	"fmt"

	"log"
	"math"
	"os"
	"runtime"
	"strconv"
	"sync"
	"time"

	"github.com/go-redis/redis"
	"github.com/syou6162/go-active-learning/lib/example"
	"github.com/syou6162/go-active-learning/lib/feature"
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
		Addr:        fmt.Sprintf("%s:6379", host),
		Password:    "", // no password set
		DB:          0,  // use default DB
		PoolSize:    100,
		MaxRetries:  4,
		PoolTimeout: time.Duration(10) * time.Second,
		IdleTimeout: time.Duration(60) * time.Second,
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

func (c *Cache) attachMetadata(examples example.Examples) error {
	for _, e := range examples {
		key := redisPrefix + ":" + e.Url
		vals, err := c.Client.HMGet(key,
			"Fv",            // 0
			"FinalUrl",      // 1
			"Title",         // 2
			"Description",   // 3
			"OgDescription", // 4
			"Body",          // 5
			"Score",         // 6
			"IsNew",         // 7
			"StatusCode",    // 8
		).Result()
		if err != nil {
			return err
		}

		// Fv
		if result, ok := vals[0].(feature.FeatureVector); ok {
			e.Fv = result
		}
		// FinalUrl
		if result, ok := vals[1].(string); ok {
			e.FinalUrl = result
		}
		// Title
		if result, ok := vals[2].(string); ok {
			e.Title = result
		}
		// Description
		if result, ok := vals[3].(string); ok {
			e.Description = result
		}
		// OgDescription
		if result, ok := vals[4].(string); ok {
			e.OgDescription = result
		}
		// Body
		if result, ok := vals[5].(string); ok {
			e.Body = result
		}
		// Score
		if result, ok := vals[6].(float64); ok {
			e.Score = result
		}
		// IsNew
		if result, ok := vals[7].(bool); ok {
			e.IsNew = result
		}
		// StatusCode
		if result, ok := vals[8].(int); ok {
			e.StatusCode = result
		}
	}
	return nil
}

func fetchMetaData(e *example.Example) {
	article := fetcher.GetArticle(e.Url)
	e.Title = article.Title
	e.FinalUrl = article.Url
	e.Description = article.Description
	e.OgDescription = article.OgDescription
	e.Body = article.Body
	e.StatusCode = article.StatusCode
	e.Fv = util.RemoveDuplicate(example.ExtractFeatures(*e))
}

func (c *Cache) SetExample(example example.Example) error {
	key := redisPrefix + ":" + example.Url

	vals := make(map[string]interface{})
	vals["Label"] = &example.Label
	vals["Fv"] = &example.Fv
	vals["Url"] = example.Url
	vals["FinalUrl"] = example.FinalUrl
	vals["Title"] = example.Title
	vals["Description"] = example.Description
	vals["OgDescription"] = example.OgDescription
	vals["Body"] = example.Body
	vals["Score"] = example.Score
	vals["IsNew"] = example.IsNew
	vals["StatusCode"] = example.StatusCode

	if err := c.Client.HMSet(key, vals).Err(); err != nil {
		return err
	}
	if err := c.Client.Expire(key, time.Hour*240).Err(); err != nil {
		return err
	}
	return nil
}

func (cache *Cache) AttachMetadata(examples example.Examples, fetchNewExamples bool) {
	batchSize := 100
	examplesList := make([]example.Examples, 0)
	n := len(examples)

	for i := 0; i < n; i += batchSize {
		max := int(math.Min(float64(i+batchSize), float64(n)))
		examplesList = append(examplesList, examples[i:max])
	}
	for _, l := range examplesList {
		if err := cache.attachMetadata(l); err != nil {
			log.Println(err.Error())
		}
		examplesWithEmptyMetaData := example.Examples{}
		for _, e := range l {
			if e.StatusCode != 200 && fetchNewExamples {
				examplesWithEmptyMetaData = append(examplesWithEmptyMetaData, e)
			}
		}
		wg := &sync.WaitGroup{}
		cpus := runtime.NumCPU()
		runtime.GOMAXPROCS(cpus)
		sem := make(chan struct{}, batchSize)
		for idx, e := range examplesWithEmptyMetaData {
			wg.Add(1)
			sem <- struct{}{}
			go func(e *example.Example, idx int) {
				defer wg.Done()
				fmt.Fprintln(os.Stderr, "Fetching("+strconv.Itoa(idx)+"): "+e.Url)
				fetchMetaData(e)
				if err := cache.SetExample(*e); err != nil {
					log.Println(err.Error())
				}
				<-sem
			}(e, idx)
		}
		wg.Wait()
	}
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
	// ToDo: take care the case when result is empty
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
