package cache

import (
	"fmt"

	"log"
	"math"
	"math/rand"
	"os"
	"runtime"
	"strconv"
	"sync"
	"time"

	"github.com/go-redis/redis"
	"github.com/syou6162/go-active-learning/lib/example"
	"github.com/syou6162/go-active-learning/lib/feature"
	"github.com/syou6162/go-active-learning/lib/fetcher"
	"github.com/syou6162/go-active-learning/lib/model"
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

func Ping() error {
	return client.Ping().Err()
}

func Close() error {
	if client != nil {
		return client.Close()
	} else {
		return nil
	}
}

func attachMetadata(examples model.Examples) error {
	for _, e := range examples {
		key := redisPrefix + ":" + e.Url
		vals, err := client.HMGet(key,
			"Fv",              // 0
			"FinalUrl",        // 1
			"Title",           // 2
			"Description",     // 3
			"OgDescription",   // 4
			"OgType",          // 5
			"OgImage",         // 6
			"Body",            // 7
			"Score",           // 8
			"IsNew",           // 9
			"StatusCode",      // 10
			"Favicon",         // 11
			"ReferringTweets", // 12
			"HatenaBookmark",  // 13
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
		// OgImage
		if result, ok := vals[6].(string); ok {
			e.OgImage = result
		}
		// Body
		if result, ok := vals[7].(string); ok {
			e.Body = result
		}
		// Score
		if result, ok := vals[8].(string); ok {
			if score, err := strconv.ParseFloat(result, 64); err == nil {
				e.Score = score
			}
		}
		// IsNew
		if result, ok := vals[9].(string); ok {
			if isNew, err := strconv.ParseBool(result); err == nil {
				e.IsNew = isNew
			}
		}
		// StatusCode
		if result, ok := vals[10].(string); ok {
			if statusCode, err := strconv.Atoi(result); err == nil {
				e.StatusCode = statusCode
			}
		}
		// Favicon
		if result, ok := vals[11].(string); ok {
			e.Favicon = result
		}
		// ReferringTweets
		if result, ok := vals[12].(string); ok {
			tweets := model.ReferringTweets{}
			if err := tweets.UnmarshalBinary([]byte(result)); err == nil {
				e.ReferringTweets = tweets
			}
		}
		// HatenaBookmark
		if result, ok := vals[13].(string); ok {
			bookmarks := model.HatenaBookmark{}
			if err := bookmarks.UnmarshalBinary([]byte(result)); err == nil {
				e.HatenaBookmark = bookmarks
			}
		}
	}
	return nil
}

func attachLightMetadata(examples model.Examples) error {
	url2Cmd := make(map[string]*redis.SliceCmd)
	url2Example := make(map[string]*model.Example)
	pipe := client.Pipeline()

	for _, e := range examples {
		key := redisPrefix + ":" + e.Url
		url2Cmd[key] = pipe.HMGet(key,
			"FinalUrl",        // 0
			"Title",           // 1
			"Description",     // 2
			"OgDescription",   // 3
			"OgType",          // 4
			"OgImage",         // 5
			"Score",           // 6
			"StatusCode",      // 7
			"Favicon",         // 8
			"ReferringTweets", // 9
			"HatenaBookmark",  // 10
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
			e.OgType = result
		}
		// OgImage
		if result, ok := vals[5].(string); ok {
			e.OgImage = result
		}
		// Score
		if result, ok := vals[6].(string); ok {
			if score, err := strconv.ParseFloat(result, 64); err == nil {
				e.Score = score
			}
		}
		// StatusCode
		if result, ok := vals[7].(string); ok {
			if statusCode, err := strconv.Atoi(result); err == nil {
				e.StatusCode = statusCode
			}
		}
		// Favicon
		if result, ok := vals[8].(string); ok {
			e.Favicon = result
		}
		// ReferringTweets
		if result, ok := vals[9].(string); ok {
			tweets := model.ReferringTweets{}
			if err := tweets.UnmarshalBinary([]byte(result)); err == nil {
				e.ReferringTweets = tweets
			}
		}
		// HatenaBookmark
		if result, ok := vals[10].(string); ok {
			bookmarks := model.HatenaBookmark{}
			if err := bookmarks.UnmarshalBinary([]byte(result)); err == nil {
				e.HatenaBookmark = bookmarks
			}
		}
	}
	return nil
}

var errorCountPrefix = "errorCountPrefix:"

func incErrorCount(url string) error {
	key := errorCountPrefix + url
	exist, err := client.Exists(key).Result()
	if err != nil {
		return err
	}
	if exist == 0 {
		hour := 24 * 10
		client.Set(key, 1, time.Hour*time.Duration(hour))
		return nil
	} else {
		if _, err = client.Incr(key).Result(); err != nil {
			return err
		}
	}
	return nil
}

func getErrorCount(url string) (int, error) {
	key := errorCountPrefix + url
	c, err := client.Exists(key).Result()
	if err != nil {
		return 0, err
	}
	if c == 0 {
		return 0, nil
	}

	cntStr, err := client.Get(key).Result()
	if err != nil {
		return 0, err
	}
	cnt, err := strconv.Atoi(cntStr)
	if err != nil {
		return 0, err
	}
	return cnt, nil
}

func fetchMetaData(e *model.Example) error {
	article, err := fetcher.GetArticle(e.Url)
	if err != nil {
		return err
	}

	e.Title = article.Title
	e.FinalUrl = article.Url
	e.Description = article.Description
	e.OgDescription = article.OgDescription
	e.OgType = article.OgType
	e.OgImage = article.OgImage
	e.Body = article.Body
	e.StatusCode = article.StatusCode
	e.Favicon = article.Favicon
	e.Fv = util.RemoveDuplicate(example.ExtractFeatures(*e))

	return nil
}

func SetExample(e model.Example) error {
	key := redisPrefix + ":" + e.Url

	vals := make(map[string]interface{})
	vals["Label"] = &e.Label
	vals["Fv"] = &e.Fv
	vals["Url"] = e.Url
	vals["FinalUrl"] = e.FinalUrl
	vals["Title"] = e.Title
	vals["Description"] = e.Description
	vals["OgDescription"] = e.OgDescription
	vals["OgType"] = e.OgType
	vals["OgImage"] = e.OgImage
	vals["Body"] = e.Body
	vals["Score"] = e.Score
	vals["IsNew"] = e.IsNew
	vals["StatusCode"] = e.StatusCode
	vals["Favicon"] = e.Favicon
	vals["ReferringTweets"] = &e.ReferringTweets
	vals["HatenaBookmark"] = &e.HatenaBookmark

	if err := client.HMSet(key, vals).Err(); err != nil {
		return err
	}

	// 一度にexpireされるとクロールも一度に走ってOOMが発生するので、多少ばらしてそれを避ける
	hour := int64(240 * rand.Float64())
	return SetExampleExpire(e, time.Hour*time.Duration(hour))
}

func SetExampleExpire(e model.Example, duration time.Duration) error {
	key := redisPrefix + ":" + e.Url
	if err := client.Expire(key, duration).Err(); err != nil {
		return err
	}
	return nil
}

func AttachMetadata(examples model.Examples, fetchNewExamples bool, useLightMetadata bool) {
	batchSize := 100
	examplesList := make([]model.Examples, 0)
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
		examplesWithEmptyMetaData := model.Examples{}
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
			go func(e *model.Example, idx int) {
				defer wg.Done()
				cnt, err := getErrorCount(e.Url)
				if err != nil {
					log.Println(err.Error())
				}
				if cnt < 5 {
					fmt.Fprintln(os.Stderr, "Fetching("+strconv.Itoa(idx)+"): "+e.Url)
					if err := fetchMetaData(e); err != nil {
						incErrorCount(e.Url)
						log.Println(err.Error())
					}
					if err := SetExample(*e); err != nil {
						log.Println(err.Error())
					}
				}
				<-sem
			}(e, idx)
		}
		wg.Wait()
	}
}

var listPrefix = "list:"

func AddExamplesToList(listName string, examples model.Examples) error {
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
