package service

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/syou6162/go-active-learning/lib/model"
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

func (app *goActiveLearningApp) DeleteAllExamples() error {
	return app.repo.DeleteAllExamples()
}

func (app *goActiveLearningApp) UpdateExampleMetadata(e model.Example) error {
	if err := app.cache.UpdateExampleMetadata(e); err != nil {
		return err
	}
	if err := app.repo.InsertOrUpdateExample(&e); err != nil {
		return err
	}
	if err := app.repo.UpdateFeatureVector(&e); err != nil {
		return err
	}
	if err := app.repo.UpdateHatenaBookmark(&e); err != nil {
		return err
	}
	return nil
}

func (app *goActiveLearningApp) UpdateExamplesMetadata(examples model.Examples) error {
	if err := app.cache.UpdateExamplesMetadata(examples); err != nil {
		return err
	}
	for _, e := range examples {
		if err := app.repo.InsertOrUpdateExample(e); err != nil {
			log.Println(fmt.Sprintf("Error occured proccessing %s %s", e.Url, err.Error()))
		}
		if err := app.repo.UpdateFeatureVector(e); err != nil {
			log.Println(fmt.Sprintf("Error occured updating feature vector %s %s", e.Url, err.Error()))
		}
		if err := app.repo.UpdateHatenaBookmark(e); err != nil {
			log.Println(fmt.Sprintf("Error occured updating feature vector %s %s", e.Url, err.Error()))
		}
	}
	return nil
}

func (app *goActiveLearningApp) UpdateExampleExpire(e model.Example, duration time.Duration) error {
	return app.cache.UpdateExampleExpire(e, duration)
}

func (app *goActiveLearningApp) AttachMetadata(examples model.Examples) error {
	fvList, err := app.repo.SearchFeatureVector(examples)
	if err != nil {
		return err
	}

	for idx, e := range examples {
		e.Fv = fvList[idx]
	}

	hatenaBookmarks, err := app.repo.SearchHatenaBookmarks(examples)
	if err != nil {
		return err
	}
	for idx, e := range examples {
		e.HatenaBookmark = hatenaBookmarks[idx]
	}
	return nil
}

func (app *goActiveLearningApp) AttachLightMetadata(examples model.Examples) error {
	hatenaBookmarks, err := app.repo.SearchHatenaBookmarks(examples)
	if err != nil {
		return err
	}
	for idx, e := range examples {
		e.HatenaBookmark = hatenaBookmarks[idx]
	}
	return nil
}

func (app *goActiveLearningApp) Fetch(examples model.Examples) {
	app.cache.Fetch(examples)
}

func (app *goActiveLearningApp) AddExamplesToList(listName string, examples model.Examples) error {
	return app.cache.AddExamplesToList(listName, examples)
}

func (app *goActiveLearningApp) GetUrlsFromList(listName string, from int64, to int64) ([]string, error) {
	return app.cache.GetUrlsFromList(listName, from, to)
}
