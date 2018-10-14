package fetcher

import (
	"testing"
)

func TestGetArticle(t *testing.T) {
	a := GetArticle("http://www.yasuhisay.info/entry/20090516/1242480413")

	if a.Title == "" {
		t.Error("Title must not be empty")
	}
	if a.Description == "" {
		t.Error("Description must not be empty")
	}
	if a.OgType != "article" {
		t.Error("OgType must be article")
	}
	if a.StatusCode != 200 {
		t.Error("StatusCode must be 200")
	}
}
