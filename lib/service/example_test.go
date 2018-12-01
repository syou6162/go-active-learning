package service_test

import (
	"testing"

	"github.com/syou6162/go-active-learning/lib/example"
	"github.com/syou6162/go-active-learning/lib/model"
	"github.com/syou6162/go-active-learning/lib/service"
)

func TestAttachMetaData(t *testing.T) {
	app, err := service.NewDefaultApp()
	if err != nil {
		t.Error(err)
	}
	defer app.Close()
	if err := app.DeleteAllExamples(); err != nil {
		t.Error("Cannot delete examples")
	}

	e1 := example.NewExample("http://b.hatena.ne.jp", model.POSITIVE)
	e2 := example.NewExample("https://www.yasuhisay.info", model.NEGATIVE)
	e3 := example.NewExample("https://github.com", model.UNLABELED)
	examples := model.Examples{e1, e2, e3}

	if examples[0].Title != "" {
		t.Errorf("Title must be empty for %s", examples[0].Url)
	}
	if len(examples[0].Fv) != 0 {
		t.Errorf("Feature vector must be empty for %s", examples[0].Url)
	}
	app.AttachMetadata(examples)

	if examples[0].Title != "" {
		t.Errorf("Title must be empty for %s", examples[0].Url)
	}
	if len(examples[0].Fv) != 0 {
		t.Errorf("Feature vector must be empty for %s", examples[0].Url)
	}

	if examples[1].OgType != "" {
		t.Errorf("OgType must be empty for %s", examples[1].Url)
	}

	app.Fetch(examples)
	err = app.UpdateExamplesMetadata(examples)
	if err != nil {
		t.Error(err)
	}
	if examples[0].Title == "" {
		t.Errorf("Title must not be empty for %s", examples[0].Url)
	}
	if len(examples[0].Fv) == 0 {
		t.Errorf("Feature vector must not be empty for %s", examples[0].Url)
	}

	if examples[1].OgType != "blog" {
		t.Errorf("OgType must be blog for %s", examples[1].Url)
	}

	examples, err = app.SearchExamplesByUlrs([]string{
		"http://b.hatena.ne.jp",
		"https://www.yasuhisay.info",
		"https://github.com",
	})
	if err != nil {
		t.Error(err)
	}
	err = app.AttachMetadata(examples)
	if err != nil {
		t.Error(err)
	}

	if examples[0].Title == "" {
		t.Errorf("Title must be empty for %s", examples[0].Url)
	}
	if len(examples[0].Fv) == 0 {
		t.Errorf("Feature vector must not be empty for %s", examples[0].Url)
	}

	if examples[1].OgType != "blog" {
		t.Errorf("OgType must be blog for %s", examples[1].Url)
	}
}
