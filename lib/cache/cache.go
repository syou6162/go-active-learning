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

var client *redis.Client
var once sync.Once

type Cache struct {
	Client *redis.Client
}

var redisPrefix = "url"

func Init() error {
	var err error
	once.Do(func() {
		host := util.GetEnv("REDIS_HOST", "localhost")
		client = redis.NewClient(&redis.Options{
			Addr:        fmt.Sprintf("%s:6379", host),
			Password:    "", // no password set
			DB:          0,  // use default DB
			PoolSize:    100,
			PoolTimeout: time.Duration(5) * time.Second,
			IdleTimeout: time.Duration(10) * time.Second,
		})
		_, err = client.Ping().Result()
		if err != nil {
			return
		}
	})
	if err != nil {
		return err
	}
	return nil
}

func Close() error {
	if client != nil {
		return client.Close()
	} else {
		return nil
	}
}

func attachMetadata(examples example.Examples) error {
	for _, e := range examples {
		key := redisPrefix + ":" + e.Url
		vals, err := client.HMGet(key,
			"Fv",            // 0
			"FinalUrl",      // 1
			"Title",         // 2
			"Description",   // 3
			"OgDescription", // 4
			"OgType",        // 5
			"Body",          // 6
			"Score",         // 7
			"IsNew",         // 8
			"StatusCode",    // 9
		).Result()
		if err != nil {
			return err
		}

		// Fv
		if result, ok := vals[0].(string); ok {
			fv := feature.FeatureVector{}
			if err := fv.UnmarshalBinary([]byte(result)); err == nil {
				e.Fv = fv
			}
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
		// OgType
		if result, ok := vals[5].(string); ok {
			e.OgType = result
		}
		// Body
		if result, ok := vals[6].(string); ok {
			e.Body = result
		}
		// Score
		if result, ok := vals[7].(string); ok {
			if score, err := strconv.ParseFloat(result, 64); err == nil {
				e.Score = score
			}
		}
		// IsNew
		if result, ok := vals[8].(string); ok {
			if isNew, err := strconv.ParseBool(result); err == nil {
				e.IsNew = isNew
			}
		}
		// StatusCode
		if result, ok := vals[9].(string); ok {
			if statusCode, err := strconv.Atoi(result); err == nil {
				e.StatusCode = statusCode
			}
		}
	}
	return nil
}

func attachLightMetadata(examples example.Examples) error {
	url2Cmd := make(map[string]*redis.SliceCmd)
	url2Example := make(map[string]*example.Example)
	pipe := client.Pipeline()

	for _, e := range examples {
		key := redisPrefix + ":" + e.Url
		url2Cmd[key] = pipe.HMGet(key,
			"FinalUrl",      // 0
			"Title",         // 1
			"Description",   // 2
			"OgDescription", // 3
			"OgType",        // 4
			"Score",         // 5
			"StatusCode",    // 6
		)
		url2Example[key] = e
	}
	_, err := pipe.Exec()
	if err != nil {
		return err
	}

	for k, cmd := range url2Cmd {
		e := url2Example[k]
		vals, err := cmd.Result()
		if err != nil {
			return err
		}
		// FinalUrl
		if result, ok := vals[0].(string); ok {
			e.FinalUrl = result
		}
		// Title
		if result, ok := vals[1].(string); ok {
			e.Title = result
		}
		// Description
		if result, ok := vals[2].(string); ok {
			e.Description = result
		}
		// OgDescription
		if result, ok := vals[3].(string); ok {
			e.OgDescription = result
		}
		// OgType
		if result, ok := vals[4].(string); ok {
			e.OgDescription = result
		}
		// Score
		if result, ok := vals[5].(string); ok {
			if score, err := strconv.ParseFloat(result, 64); err == nil {
				e.Score = score
			}
		}
		// StatusCode
		if result, ok := vals[6].(string); ok {
			if statusCode, err := strconv.Atoi(result); err == nil {
				e.StatusCode = statusCode
			}
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
	e.OgType = article.OgType
	e.Body = article.Body
	e.StatusCode = article.StatusCode
	e.Fv = util.RemoveDuplicate(example.ExtractFeatures(*e))
}

func SetExample(example example.Example) error {
	key := redisPrefix + ":" + example.Url

	vals := make(map[string]interface{})
	vals["Label"] = &example.Label
	vals["Fv"] = &example.Fv
	vals["Url"] = example.Url
	vals["FinalUrl"] = example.FinalUrl
	vals["Title"] = example.Title
	vals["Description"] = example.Description
	vals["OgDescription"] = example.OgDescription
	vals["OgType"] = example.OgType
	vals["Body"] = example.Body
	vals["Score"] = example.Score
	vals["IsNew"] = example.IsNew
	vals["StatusCode"] = example.StatusCode

	if err := client.HMSet(key, vals).Err(); err != nil {
		return err
	}
	if err := client.Expire(key, time.Hour*240).Err(); err != nil {
		return err
	}
	return nil
}

func AttachMetadata(examples example.Examples, fetchNewExamples bool, useLightMetadata bool) {
	batchSize := 100
	examplesList := make([]example.Examples, 0)
	n := len(examples)

	for i := 0; i < n; i += batchSize {
		max := int(math.Min(float64(i+batchSize), float64(n)))
		examplesList = append(examplesList, examples[i:max])
	}
	for _, l := range examplesList {
		if useLightMetadata {
			if err := attachLightMetadata(l); err != nil {
				log.Println(err.Error())
			}
		} else {
			if err := attachMetadata(l); err != nil {
				log.Println(err.Error())
			}
		}
		if !fetchNewExamples {
			continue
		}
		examplesWithEmptyMetaData := example.Examples{}
		for _, e := range l {
			if e.StatusCode != 200 {
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
				if err := SetExample(*e); err != nil {
					log.Println(err.Error())
				}
				<-sem
			}(e, idx)
		}
		wg.Wait()
	}
}

var listPrefix = "list:"

func AddExamplesToList(listName string, examples example.Examples) error {
	if err := client.Del(listPrefix + listName).Err(); err != nil {
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
	err := client.ZAdd(listPrefix+listName, result...).Err()
	if err != nil {
		return err
	}
	return nil
}

func GetUrlsFromList(listName string, from int64, to int64) ([]string, error) {
	sliceCmd := client.ZRevRange(listPrefix+listName, from, to)
	if sliceCmd.Err() != nil {
		return nil, sliceCmd.Err()
	}
	return sliceCmd.Val(), nil
}
