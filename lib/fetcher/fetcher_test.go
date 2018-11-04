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
	if expectedFaviconPath != a.Favicon {
		t.Errorf("Favicon: %s should be %s", a.Favicon, expectedFaviconPath)
	}

	url = "https://arxiv.org/abs/1810.08403"
	a, err = GetArticle(url)
	if err != nil {
		t.Error(fmt.Sprintf("Error must not occur for this url: %s", url))
	}
	expectedFaviconPath = "https://arxiv.org/favicon.ico"
	if expectedFaviconPath != a.Favicon {
		t.Errorf("Favicon: %s should be %s", a.Favicon, expectedFaviconPath)
	}

	url = "https://www.lifehacker.jp/2018/11/amazon-impact-absorption-case.html"
	a, err = GetArticle(url)
	if err != nil {
		t.Error(fmt.Sprintf("Error must not occur for this url: %s", url))
	}
	expectedFaviconPath = "https://www.lifehacker.jp/assets/common/img/favicon.ico"
	if expectedFaviconPath != a.Favicon {
		t.Errorf("Favicon: %s should be %s", a.Favicon, expectedFaviconPath)
	}

	url = "https://peterroelants.github.io/"
	a, err = GetArticle(url)
	if err != nil {
		t.Error(fmt.Sprintf("Error must not occur for this url: %s", url))
	}
	expectedFaviconPath = "https://peterroelants.github.io/images/favicon/apple-icon-57x57.png"
	if expectedFaviconPath != a.Favicon {
		t.Errorf("Favicon: %s should be %s", a.Favicon, expectedFaviconPath)
	}

	url = "https://www.getrevue.co/profile/icoxfog417/issues/weekly-machine-learning-79-121292"
	a, err = GetArticle(url)
	if err != nil {
		t.Error(fmt.Sprintf("Error must not occur for this url: %s", url))
	}
	expectedFaviconPath = "https://d3jbm9h03wxzi9.cloudfront.net/assets/favicon-84fc7f228d52c2410eb7aa839e279caeaa491588c7c75229ed33e1c7f69fe75d.ico"
	if expectedFaviconPath != a.Favicon {
		t.Errorf("Favicon: %s should be %s", a.Favicon, expectedFaviconPath)
	}

	url = "https://ai.googleblog.com/2018/11/open-sourcing-bert-state-of-art-pre.html"
	a, err = GetArticle(url)
	if err != nil {
		t.Error(fmt.Sprintf("Error must not occur for this url: %s", url))
	}
	expectedFaviconPath = "https://ai.googleblog.com/favicon.ico"
	if expectedFaviconPath != a.Favicon {
		t.Errorf("Favicon: %s should be %s", a.Favicon, expectedFaviconPath)
	}
}
