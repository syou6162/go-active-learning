package main

import (
	"time"
	"github.com/advancedlogic/GoOse"
)

type Article struct {
	Url         string
	Title       string
	Description string
	Body        string
}

func getArticle(url string) Article {
	g := goose.New()
	article, err := g.ExtractFromURL(url)
	if err != nil {
		return Article{}
	}
	return Article{url, article.Title, article.MetaDescription, article.CleanedText}
}

func GetArticle(url string) Article {
	ch := make(chan Article, 1)
	go func(url string) { ch <- getArticle(url) }(url)
	select {
	case article := <-ch:
		return article
	case <-time.After(10 * time.Second):
		return Article{}
	}
}
