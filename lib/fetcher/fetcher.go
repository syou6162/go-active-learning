package fetcher

import (
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"net/url"
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
	OgImage       string
	Body          string
	StatusCode    int
	Favicon       string
}

var articleFetcher = http.Client{
	Transport: &http.Transport{
		MaxIdleConns:        0,
		MaxIdleConnsPerHost: 100,
	},
	Timeout: time.Duration(5 * time.Second),
}

func updateMetaDescriptionIfArxiv(article *goose.Article, origUrl string, finalUrl string, html []byte) error {
	arxivUrl := "https://arxiv.org/abs/"
	if strings.Contains(origUrl, arxivUrl) || strings.Contains(finalUrl, arxivUrl) {
		// article.Docでもいけそうだが、gooseが中で書き換えていてダメ。Documentを作りなおす
		doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(html)))
		if err != nil {
			return err
		}
		article.MetaDescription = doc.Find(".abstract").Text()
	}
	return nil
}

func GetArticle(origUrl string) (*Article, error) {
	g := goose.New()
	resp, err := articleFetcher.Get(origUrl)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	html, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if !utf8.Valid(html) {
		return nil, err
	}

	article, err := g.ExtractFromRawHTML(resp.Request.URL.String(), string(html))
	if err != nil {
		return nil, err
	}

	finalUrl := article.CanonicalLink
	if finalUrl == "" {
		finalUrl = resp.Request.URL.String()
	}

	updateMetaDescriptionIfArxiv(article, origUrl, finalUrl, html)

	favicon := ""
	if u, err := url.Parse(article.MetaFavicon); err == nil {
		if u.IsAbs() {
			favicon = article.MetaFavicon
		}
	}

	return &Article{
		Url:           finalUrl,
		Title:         article.Title,
		Description:   article.MetaDescription,
		OgDescription: article.MetaOgDescription,
		OgType:        article.MetaOgType,
		OgImage:       article.MetaOgImage,
		Body:          article.CleanedText,
		StatusCode:    resp.StatusCode,
		Favicon:       favicon,
	}, nil
}
