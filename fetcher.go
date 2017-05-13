package main

import (
	"github.com/PuerkitoBio/goquery"
)

func GetTitle(url string) (string, error) {
	var title string
	doc, err := goquery.NewDocument(url)
	if err != nil {
		return title, err
	}

	doc.Find("title").Each(func(index int, item *goquery.Selection) {
		title = item.Text()
	})
	return title, nil
}

