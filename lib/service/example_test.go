package service_test

import (
	"testing"

	"github.com/syou6162/go-active-learning/lib/example"
	"github.com/syou6162/go-active-learning/lib/model"
	"github.com/syou6162/go-active-learning/lib/service"
)

func findExampleByurl(examples model.Examples, url string) *model.Example {
	for _, e := range examples {
		if e.Url == url {
			return e
		}
	}
	return nil
}

func TestAttachMetaData(t *testing.T) {
	app, err := service.NewDefaultApp()
	if err != nil {
		t.Error(err)
	}
	defer app.Close()
	if err := app.DeleteAllExamples(); err != nil {
		t.Error("Cannot delete examples")
	}

	hatebuUrl := "https://b.hatena.ne.jp"
	myBlogUrl := "https://www.yasuhisay.info"
	githubUrl := "https://github.com"
	e1 := example.NewExample(hatebuUrl, model.POSITIVE)
	e2 := example.NewExample(myBlogUrl, model.NEGATIVE)
	e3 := example.NewExample(githubUrl, model.UNLABELED)
	examples := model.Examples{e1, e2, e3}

	hatebu := findExampleByurl(examples, hatebuUrl)
	if hatebu == nil {
		t.Errorf("Cannot find %s", hatebuUrl)
	}
	if hatebu.Title != "" {
		t.Errorf("Title must be empty for %s", hatebu.Url)
	}
	if len(hatebu.Fv) != 0 {
		t.Errorf("Feature vector must be empty for %s", hatebu.Url)
	}
	app.AttachMetadataIncludingFeatureVector(examples, 10, 10)

	if hatebu.Title != "" {
		t.Errorf("Title must be empty for %s", hatebu.Url)
	}
	if len(hatebu.Fv) != 0 {
		t.Errorf("Feature vector must be empty for %s", hatebu.Url)
	}

	myBlog := findExampleByurl(examples, myBlogUrl)
	if myBlog == nil {
		t.Errorf("Cannot find %s", myBlogUrl)
	}
	if myBlog.OgType != "" {
		t.Errorf("OgType must be empty for %s", myBlog.Url)
	}

	app.Fetch(examples)
	for _, e := range examples {
		err = app.UpdateOrCreateExample(e)
		if err != nil {
			t.Error(err)
		}
		err = app.UpdateFeatureVector(e)
		if err != nil {
			t.Error(err)
		}
	}
	if hatebu.Title == "" {
		t.Errorf("Title must not be empty for %s", hatebu.Url)
	}
	if len(hatebu.Fv) == 0 {
		t.Errorf("Feature vector must not be empty for %s", hatebu.Url)
	}

	if myBlog.OgType != "blog" {
		t.Errorf("OgType must be blog for %s", myBlog.Url)
	}

	examples, err = app.SearchExamplesByIds([]int{e1.Id, e2.Id, e3.Id})
	if err != nil {
		t.Error(err)
	}
	err = app.AttachMetadataIncludingFeatureVector(examples, 10, 10)
	if err != nil {
		t.Error(err)
	}

	if hatebu.Title == "" {
		t.Errorf("Title must be empty for %s", hatebu.Url)
	}
	if len(hatebu.Fv) == 0 {
		t.Errorf("Feature vector must not be empty for %s", hatebu.Url)
	}

	if myBlog.OgType != "blog" {
		t.Errorf("OgType must be blog for %s", myBlog.Url)
	}
}

func TestGetRecommendation(t *testing.T) {
	app, err := service.NewDefaultApp()
	if err != nil {
		t.Error(err)
	}
	defer app.Close()
	if err := app.DeleteAllExamples(); err != nil {
		t.Error("Cannot delete examples")
	}

	e1 := example.NewExample("http://hoge1.com", model.POSITIVE)
	e2 := example.NewExample("http://hoge2.com", model.NEGATIVE)
	e3 := example.NewExample("http://hoge3.com", model.UNLABELED)
	examples := model.Examples{e1, e2, e3}
	for _, e := range examples {
		err = app.UpdateOrCreateExample(e)
		if err != nil {
			t.Error(err)
		}
	}

	listName := "general"
	err = app.UpdateRecommendation(listName, examples)
	if err != nil {
		t.Error(err)
	}
	examples, err = app.GetRecommendation(listName)
	if err != nil {
		t.Error(err)
	}
	if len(examples) != 3 {
		t.Errorf("len(examples) should be 3, but %d", len(examples))
	}
}
