package service

import (
	"bufio"
	"io"
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

func (app *goActiveLearningApp) SearchExamplesByUlr(url string) (*model.Example, error) {
	return app.repo.SearchExamplesByUlr(url)
}

func (app *goActiveLearningApp) SearchExamplesByUlrs(urls []string) (model.Examples, error) {
	return app.repo.SearchExamplesByUlrs(urls)
}

func (app *goActiveLearningApp) DeleteAllExamples() error {
	return app.repo.DeleteAllExamples()
}
