package main

import (
	"github.com/advancedlogic/GoOse"
)

type Article struct {
	Url         string
	Title       string
	Description string
	Body        string
	RawHTML     string
}

func GetArticle(url string) Article {
	g := goose.New()
	article, err := g.ExtractFromURL(url)
	if err != nil {
		return Article{}
	}
	return Article{url, article.Title, article.MetaDescription, article.CleanedText, article.RawHTML}
}
