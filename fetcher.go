package main

import (
	"io/ioutil"
	"net/http"
	"time"

	"github.com/advancedlogic/GoOse"
)

type Article struct {
	Url         string
	Title       string
	Description string
	Body        string
	StatusCode  int
}

func GetArticle(url string) Article {
	g := goose.New()
	timeout := time.Duration(5 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}

	resp, err := client.Get(url)
	if err != nil {
		return Article{}
	}
	defer resp.Body.Close()

	html, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return Article{StatusCode: resp.StatusCode}
	}

	article, err := g.ExtractFromRawHTML(url, string(html))
	if err != nil {
		return Article{StatusCode: resp.StatusCode}
	}

	return Article{url, article.Title, article.MetaDescription, article.CleanedText, resp.StatusCode}
}
