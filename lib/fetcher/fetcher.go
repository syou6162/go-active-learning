package fetcher

import (
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"unicode/utf8"

	"github.com/PuerkitoBio/goquery"
	"github.com/syou6162/GoOse"
)

type Article struct {
	Url           string
	Title         string
	Description   string
	OgDescription string
	OgType        string
	Body          string
	StatusCode    int
}

var articleFetcher = http.Client{
	Transport: &http.Transport{
		MaxIdleConns:        0,
		MaxIdleConnsPerHost: 100,
	},
	Timeout: time.Duration(5 * time.Second),
}

func GetArticle(url string) Article {
	g := goose.New()
	resp, err := articleFetcher.Get(url)
	if err != nil {
		return Article{}
	}
	defer resp.Body.Close()

	html, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return Article{StatusCode: resp.StatusCode}
	}

	if !utf8.Valid(html) {
		return Article{Url: resp.Request.URL.String(), StatusCode: resp.StatusCode}
	}

	article, err := g.ExtractFromRawHTML(resp.Request.URL.String(), string(html))
	if err != nil {
		return Article{StatusCode: resp.StatusCode}
	}

	finalUrl := article.CanonicalLink
	if finalUrl == "" {
		finalUrl = resp.Request.URL.String()
	}

	arxivUrl := "https://arxiv.org/abs/"
	if strings.Contains(url, arxivUrl) || strings.Contains(finalUrl, arxivUrl) {
		// article.Docでもいけそうだが、gooseが中で書き換えていてダメ。Documentを作りなおす
		doc, _ := goquery.NewDocumentFromReader(strings.NewReader(string(html)))
		article.MetaDescription = doc.Find(".abstract").Text()
	}

	return Article{
		finalUrl,
		article.Title,
		article.MetaDescription,
		article.MetaOgDescription,
		article.MetaOgType,
		article.CleanedText,
		resp.StatusCode,
	}
}
