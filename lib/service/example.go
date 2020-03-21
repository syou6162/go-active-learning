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

func (app *goActiveLearningApp) UpdateOrCreateExample(e *model.Example) error {
	return app.repo.UpdateOrCreateExample(e)
}

func (app *goActiveLearningApp) UpdateScore(e *model.Example) error {
	return app.repo.UpdateScore(e)
}

func (app *goActiveLearningApp) InsertExampleFromScanner(scanner *bufio.Scanner) (*model.Example, error) {
	return app.repo.InsertExampleFromScanner(scanner)
}

func (app *goActiveLearningApp) InsertExamplesFromReader(reader io.Reader) error {
	return app.repo.InsertExamplesFromReader(reader)
}

func (app *goActiveLearningApp) SearchExamples() (model.Examples, error) {
	return app.repo.SearchExamples()
}

func (app *goActiveLearningApp) SearchRecentExamples(from time.Time, limit int) (model.Examples, error) {
	return app.repo.SearchRecentExamples(from, limit)
}

func (app *goActiveLearningApp) SearchRecentExamplesByHost(host string, from time.Time, limit int) (model.Examples, error) {
	return app.repo.SearchRecentExamplesByHost(host, from, limit)
}

func (app *goActiveLearningApp) SearchExamplesByLabel(label model.LabelType, limit int) (model.Examples, error) {
	return app.repo.SearchExamplesByLabel(label, limit)
}

func (app *goActiveLearningApp) SearchLabeledExamples(limit int) (model.Examples, error) {
	return app.repo.SearchLabeledExamples(limit)
}

func (app *goActiveLearningApp) SearchPositiveExamples(limit int) (model.Examples, error) {
	return app.repo.SearchPositiveExamples(limit)
}

func (app *goActiveLearningApp) SearchNegativeExamples(limit int) (model.Examples, error) {
	return app.repo.SearchNegativeExamples(limit)
}

func (app *goActiveLearningApp) SearchUnlabeledExamples(limit int) (model.Examples, error) {
	return app.repo.SearchUnlabeledExamples(limit)
}

func (app *goActiveLearningApp) SearchPositiveScoredExamples(limit int) (model.Examples, error) {
	return app.repo.SearchPositiveScoredExamples(limit)
}

func (app *goActiveLearningApp) FindExampleByUlr(url string) (*model.Example, error) {
	return app.repo.FindExampleByUlr(url)
}

func (app *goActiveLearningApp) FindExampleById(id int) (*model.Example, error) {
	return app.repo.FindExampleById(id)
}

func (app *goActiveLearningApp) SearchExamplesByUlrs(urls []string) (model.Examples, error) {
	return app.repo.SearchExamplesByUlrs(urls)
}

func (app *goActiveLearningApp) SearchExamplesByIds(ids []int) (model.Examples, error) {
	return app.repo.SearchExamplesByIds(ids)
}

func (app *goActiveLearningApp) SearchExamplesByKeywords(keywords []string, aggregator string, limit int) (model.Examples, error) {
	return app.repo.SearchExamplesByKeywords(keywords, aggregator, limit)
}

func (app *goActiveLearningApp) DeleteAllExamples() error {
	return app.repo.DeleteAllExamples()
}

func (app *goActiveLearningApp) CountPositiveExamples() (int, error) {
	return app.repo.CountPositiveExamples()
}

func (app *goActiveLearningApp) CountNegativeExamples() (int, error) {
	return app.repo.CountNegativeExamples()
}

func (app *goActiveLearningApp) CountUnlabeledExamples() (int, error) {
	return app.repo.CountUnlabeledExamples()
}

func (app *goActiveLearningApp) UpdateFeatureVector(e *model.Example) error {
	return app.repo.UpdateFeatureVector(e)
}

func (app *goActiveLearningApp) UpdateHatenaBookmark(e *model.Example) error {
	return app.repo.UpdateHatenaBookmark(e)
}

func (app *goActiveLearningApp) UpdateOrCreateReferringTweets(e *model.Example) error {
	return app.repo.UpdateOrCreateReferringTweets(e)
}

func (app *goActiveLearningApp) UpdateTweetLabel(exampleId int, idStr string, label model.LabelType) error {
	return app.repo.UpdateTweetLabel(exampleId, idStr, label)
}

func (app *goActiveLearningApp) SearchReferringTweets(limit int) (model.ReferringTweets, error) {
	return app.repo.SearchReferringTweets(limit)
}

func (app *goActiveLearningApp) SearchPositiveReferringTweets(limit int) (model.ReferringTweets, error) {
	return app.repo.SearchPositiveReferringTweets(limit)
}

func (app *goActiveLearningApp) SearchNegativeReferringTweets(limit int) (model.ReferringTweets, error) {
	return app.repo.SearchNegativeReferringTweets(limit)
}

func (app *goActiveLearningApp) SearchUnlabeledReferringTweets(limit int) (model.ReferringTweets, error) {
	return app.repo.SearchUnlabeledReferringTweets(limit)
}

func (app *goActiveLearningApp) SearchRecentReferringTweetsWithHighScore(from time.Time, scoreThreshold float64, limit int) (model.ReferringTweets, error) {
	return app.repo.SearchRecentReferringTweetsWithHighScore(from, scoreThreshold, limit)
}

func hatenaBookmarkByExampleId(hatenaBookmarks []*model.HatenaBookmark) map[int]*model.HatenaBookmark {
	result := make(map[int]*model.HatenaBookmark)
	for _, hb := range hatenaBookmarks {
		result[hb.ExampleId] = hb
	}
	return result
}

func (app *goActiveLearningApp) AttachMetadataIncludingFeatureVector(examples model.Examples, bookmarkLimit int, tweetLimit int) error {
	// make sure that example id must be filled
	for _, e := range examples {
		if e.Id == 0 {
			tmp, err := app.FindExampleByUlr(e.Url)
			if err != nil {
				return err
			}
			e.Id = tmp.Id
		}
	}

	fvList, err := app.repo.SearchFeatureVector(examples)
	if err != nil {
		return err
	}

	for _, e := range examples {
		if fv, ok := fvList[e.Id]; ok {
			e.Fv = fv
		}
	}

	return app.AttachMetadata(examples, bookmarkLimit, tweetLimit)
}

func (app *goActiveLearningApp) AttachMetadata(examples model.Examples, bookmarkLimit int, tweetLimit int) error {
	hatenaBookmarks, err := app.repo.SearchHatenaBookmarks(examples, bookmarkLimit)
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

	referringTweetsById, err := app.repo.SearchReferringTweetsList(examples, tweetLimit)
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

func (app *goActiveLearningApp) UpdateRelatedExamples(related model.RelatedExamples) error {
	return app.repo.UpdateRelatedExamples(related)
}

func (app *goActiveLearningApp) SearchRelatedExamples(e *model.Example) (model.Examples, error) {
	related, err := app.repo.FindRelatedExamples(e)
	if err != nil {
		return nil, err
	}
	return app.repo.SearchExamplesByIds(related.RelatedExampleIds)
}

func (app *goActiveLearningApp) UpdateTopAccessedExampleIds(exampleIds []int) error {
	return app.repo.UpdateTopAccessedExampleIds(exampleIds)
}

func (app *goActiveLearningApp) SearchTopAccessedExamples() (model.Examples, error) {
	exampleIds, err := app.repo.SearchTopAccessedExampleIds()
	if err != nil {
		return nil, err
	}
	return app.repo.SearchExamplesByIds(exampleIds)
}

func (app *goActiveLearningApp) UpdateRecommendation(listName string, examples model.Examples) error {
	listType, err := model.GetRecommendationListType(listName)
	if err != nil {
		return err
	}

	exampleIds := make([]int, 0)
	for _, e := range examples {
		exampleIds = append(exampleIds, e.Id)
	}

	rec := model.Recommendation{RecommendationListType: listType, ExampleIds: exampleIds}
	return app.repo.UpdateRecommendation(rec)
}

func (app *goActiveLearningApp) GetRecommendation(listName string) (model.Examples, error) {
	listType, err := model.GetRecommendationListType(listName)
	if err != nil {
		return nil, err
	}
	rec, err := app.repo.FindRecommendation(listType)
	return app.repo.SearchExamplesByIds(rec.ExampleIds)
}

func (app *goActiveLearningApp) splitExamplesByStatusOK(examples model.Examples) (model.Examples, model.Examples, error) {
	urls := make([]string, 0)
	exampleByurl := make(map[string]*model.Example)
	for _, e := range examples {
		exampleByurl[e.Url] = e
		urls = append(urls, e.Url)
	}
	tmpExamples, err := app.SearchExamplesByUlrs(urls)
	if err != nil {
		return nil, nil, err
	}

	examplesWithMetaData := model.Examples{}
	examplesWithEmptyMetaData := model.Examples{}
	for _, e := range tmpExamples {
		if e.StatusCode == http.StatusOK {
			examplesWithMetaData = append(examplesWithMetaData, exampleByurl[e.Url])
			delete(exampleByurl, e.Url)
		} else {
			examplesWithEmptyMetaData = append(examplesWithEmptyMetaData, exampleByurl[e.Url])
			delete(exampleByurl, e.Url)
		}
	}
	for _, e := range exampleByurl {
		examplesWithEmptyMetaData = append(examplesWithEmptyMetaData, e)
	}
	return examplesWithMetaData, examplesWithEmptyMetaData, nil
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

	now := time.Now()
	tooOldDate := time.Date(2000, time.January, 1, 1, 1, 0, 0, time.UTC)
	if article.PublishDate != nil && (now.After(*article.PublishDate) || tooOldDate.Before(*article.PublishDate)) {
		e.CreatedAt = *article.PublishDate
		e.UpdatedAt = *article.PublishDate
	}

	fv := util.RemoveDuplicate(example.ExtractFeatures(*e))
	if len(fv) > 100000 {
		return fmt.Errorf("too large features (N = %d) for %s", len(fv), e.FinalUrl)
	}
	e.Fv = fv

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
		examplesWithMetaData, examplesWithEmptyMetaData, err := app.splitExamplesByStatusOK(l)
		if err != nil {
			log.Println(err.Error())
		}
		// ToDo: 本当に必要か考える
		app.AttachMetadataIncludingFeatureVector(examplesWithMetaData, 0, 0)

		wg := &sync.WaitGroup{}
		cpus := runtime.NumCPU()
		runtime.GOMAXPROCS(cpus)
		sem := make(chan struct{}, batchSize)
		for idx, e := range examplesWithEmptyMetaData {
			wg.Add(1)
			sem <- struct{}{}
			go func(e *model.Example, idx int) {
				defer wg.Done()
				cnt, err := app.repo.GetErrorCount(e)
				if err != nil {
					log.Println(err.Error())
				}
				if cnt < 5 {
					fmt.Fprintln(os.Stderr, "Fetching("+strconv.Itoa(idx)+"): "+e.Url)
					if err := fetchMetaData(e); err != nil {
						app.repo.IncErrorCount(e)
						log.Println(err.Error())
					}
				}
				<-sem
			}(e, idx)
		}
		wg.Wait()
	}
}
