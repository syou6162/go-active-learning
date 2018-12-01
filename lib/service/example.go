package service

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net/http"
	"runtime"
	"time"

	"math"
	"os"
	"strconv"
	"sync"

	"github.com/syou6162/go-active-learning/lib/example"
	"github.com/syou6162/go-active-learning/lib/fetcher"
	"github.com/syou6162/go-active-learning/lib/model"
	"github.com/syou6162/go-active-learning/lib/util"
)

func (app *goActiveLearningApp) InsertOrUpdateExample(e *model.Example) error {
	return app.repo.InsertOrUpdateExample(e)
}

func (app *goActiveLearningApp) InsertExampleFromScanner(scanner *bufio.Scanner) (*model.Example, error) {
	return app.repo.InsertExampleFromScanner(scanner)
}

func (app *goActiveLearningApp) InsertExamplesFromReader(reader io.Reader) error {
	return app.repo.InsertExamplesFromReader(reader)
}

func (app *goActiveLearningApp) ReadExamples() (model.Examples, error) {
	return app.repo.ReadExamples()
}

func (app *goActiveLearningApp) ReadRecentExamples(from time.Time) (model.Examples, error) {
	return app.repo.ReadRecentExamples(from)
}

func (app *goActiveLearningApp) ReadExamplesByLabel(label model.LabelType, limit int) (model.Examples, error) {
	return app.repo.ReadExamplesByLabel(label, limit)
}

func (app *goActiveLearningApp) ReadLabeledExamples(limit int) (model.Examples, error) {
	return app.repo.ReadLabeledExamples(limit)
}

func (app *goActiveLearningApp) ReadPositiveExamples(limit int) (model.Examples, error) {
	return app.repo.ReadPositiveExamples(limit)
}

func (app *goActiveLearningApp) ReadNegativeExamples(limit int) (model.Examples, error) {
	return app.repo.ReadNegativeExamples(limit)
}

func (app *goActiveLearningApp) ReadUnlabeledExamples(limit int) (model.Examples, error) {
	return app.repo.ReadUnlabeledExamples(limit)
}

func (app *goActiveLearningApp) FindExampleByUlr(url string) (*model.Example, error) {
	return app.repo.FindExampleByUlr(url)
}

func (app *goActiveLearningApp) SearchExamplesByUlrs(urls []string) (model.Examples, error) {
	return app.repo.SearchExamplesByUlrs(urls)
}

func (app *goActiveLearningApp) SearchExamplesByKeywords(keywords []string, aggregator string, limit int) (model.Examples, error) {
	return app.repo.SearchExamplesByKeywords(keywords, aggregator, limit)
}

func (app *goActiveLearningApp) DeleteAllExamples() error {
	return app.repo.DeleteAllExamples()
}

func (app *goActiveLearningApp) UpdateExampleMetadata(e model.Example) error {
	if err := app.repo.InsertOrUpdateExample(&e); err != nil {
		log.Println(fmt.Sprintf("Error occured proccessing %s %s", e.Url, err.Error()))
	}
	if err := app.repo.UpdateFeatureVector(&e); err != nil {
		log.Println(fmt.Sprintf("Error occured updating feature vector %s %s", e.Url, err.Error()))
	}
	if err := app.repo.UpdateHatenaBookmark(&e); err != nil {
		log.Println(fmt.Sprintf("Error occured updating bookmark info %s %s", e.Url, err.Error()))
	}
	if err := app.repo.UpdateReferringTweets(&e); err != nil {
		log.Println(fmt.Sprintf("Error occured updating twitter info %s %s", e.Url, err.Error()))
	}
	return nil
}

func (app *goActiveLearningApp) UpdateExamplesMetadata(examples model.Examples) error {
	for _, e := range examples {
		app.UpdateExampleMetadata(*e)
	}
	return nil
}

func hatenaBookmarkByExampleId(hatenaBookmarks []*model.HatenaBookmark) map[int]*model.HatenaBookmark {
	result := make(map[int]*model.HatenaBookmark)
	for _, hb := range hatenaBookmarks {
		result[hb.ExampleId] = hb
	}
	return result
}

func (app *goActiveLearningApp) AttachMetadata(examples model.Examples) error {
	fvList, err := app.repo.SearchFeatureVector(examples)
	if err != nil {
		return err
	}

	// ToDo: Rewrite
	if len(fvList) > 0 {
		for idx, e := range examples {
			e.Fv = fvList[idx]
		}
	}

	hatenaBookmarks, err := app.repo.SearchHatenaBookmarks(examples)
	if err != nil {
		return err
	}
	hbByid := hatenaBookmarkByExampleId(hatenaBookmarks)
	for _, e := range examples {
		if b, ok := hbByid[e.Id]; ok {
			e.HatenaBookmark = b
		} else {
			e.HatenaBookmark = &model.HatenaBookmark{Bookmarks: []*model.Bookmark{}}
		}
	}

	referringTweetsById, err := app.repo.SearchReferringTweetsList(examples)
	if err != nil {
		return err
	}
	for _, e := range examples {
		if t, ok := referringTweetsById[e.Id]; ok {
			e.ReferringTweets = &t
		} else {
			e.ReferringTweets = &model.ReferringTweets{}
		}
	}

	return nil
}

func (app *goActiveLearningApp) AttachLightMetadata(examples model.Examples) error {
	hatenaBookmarks, err := app.repo.SearchHatenaBookmarks(examples)
	if err != nil {
		return err
	}
	hbByid := hatenaBookmarkByExampleId(hatenaBookmarks)
	for _, e := range examples {
		if b, ok := hbByid[e.Id]; ok {
			e.HatenaBookmark = b
		} else {
			e.HatenaBookmark = &model.HatenaBookmark{Bookmarks: []*model.Bookmark{}}
		}
	}

	referringTweetsById, err := app.repo.SearchReferringTweetsList(examples)
	if err != nil {
		return err
	}
	for _, e := range examples {
		if t, ok := referringTweetsById[e.Id]; ok {
			e.ReferringTweets = &t
		} else {
			e.ReferringTweets = &model.ReferringTweets{}
		}
	}
	return nil
}

func (app *goActiveLearningApp) AddExamplesToList(listName string, examples model.Examples) error {
	return app.cache.AddExamplesToList(listName, examples)
}

func (app *goActiveLearningApp) GetUrlsFromList(listName string, from int64, to int64) ([]string, error) {
	return app.cache.GetUrlsFromList(listName, from, to)
}

func (app *goActiveLearningApp) ExistMetadata(e model.Example) bool {
	tmp, err := app.FindExampleByUlr(e.Url)
	if err != nil {
		return false
	}

	if tmp.StatusCode == http.StatusOK {
		return true
	}
	return false
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

func (app *goActiveLearningApp) Fetch(examples model.Examples) {
	batchSize := 100
	examplesList := make([]model.Examples, 0)
	n := len(examples)

	for i := 0; i < n; i += batchSize {
		max := int(math.Min(float64(i+batchSize), float64(n)))
		examplesList = append(examplesList, examples[i:max])
	}
	for _, l := range examplesList {
		examplesWithMetaData := model.Examples{}
		examplesWithEmptyMetaData := model.Examples{}
		for _, e := range l {
			if !app.ExistMetadata(*e) {
				examplesWithEmptyMetaData = append(examplesWithEmptyMetaData, e)
			} else {
				examplesWithMetaData = append(examplesWithMetaData, e)
			}
		}
		app.AttachMetadata(examplesWithMetaData)

		wg := &sync.WaitGroup{}
		cpus := runtime.NumCPU()
		runtime.GOMAXPROCS(cpus)
		sem := make(chan struct{}, batchSize)
		for idx, e := range examplesWithEmptyMetaData {
			wg.Add(1)
			sem <- struct{}{}
			go func(e *model.Example, idx int) {
				defer wg.Done()
				cnt, err := app.cache.GetErrorCount(e.Url)
				if err != nil {
					log.Println(err.Error())
				}
				if cnt < 5 {
					fmt.Fprintln(os.Stderr, "Fetching("+strconv.Itoa(idx)+"): "+e.Url)
					if err := fetchMetaData(e); err != nil {
						app.cache.IncErrorCount(e.Url)
						log.Println(err.Error())
					}
				}
				<-sem
			}(e, idx)
		}
		wg.Wait()
	}
}
