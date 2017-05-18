package main

import (
	"io/ioutil"
	"net/http"
	"github.com/advancedlogic/GoOse"
)

type Article struct {
	Url         string
	Title       string
	Description string
	Body        string
	RawHTML     string
	StatusCode  int
}

func GetArticle(url string) Article {
	g := goose.New()
	resp, err := http.Get(url)
	if err != nil {
		return Article{}
	}
	defer resp.Body.Close()

	html, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return Article{StatusCode:resp.StatusCode}
	}

	article, err := g.ExtractFromRawHTML(url, string(html))
	if err != nil {
		return Article{StatusCode:resp.StatusCode}
	}

	return Article{url, article.Title, article.MetaDescription, article.CleanedText, article.RawHTML, resp.StatusCode}
}
