package fetcher

import (
	"fmt"
	"testing"
)

func TestGetArticle(t *testing.T) {
	a, err := GetArticle("http://www.yasuhisay.info/entry/20090516/1242480413")
	if err != nil {
		t.Error(err.Error())
	}

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

func TestGetArticleWithInvalidEncoding(t *testing.T) {
	url := "http://www.atmarkit.co.jp/ait/articles/1702/20/news021.html"
	_, err := GetArticle(url)
	if err == nil {
		t.Error(fmt.Sprintf("Error must occur for this url: %s", url))
	}
}

func TestFavicon(t *testing.T) {
	url := "https://twitter.com/facebookai/status/1057764513582215168"
	a, err := GetArticle(url)
	if err != nil {
		t.Error(fmt.Sprintf("Error must not occur for this url: %s", url))
	}
	expectedFaviconPath := "https://abs.twimg.com/favicons/favicon.ico"
	println(a.Favicon)
	if expectedFaviconPath != a.Favicon {
		t.Errorf("Favicon: %s should be %s", a.Favicon, expectedFaviconPath)
	}
}
