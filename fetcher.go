package main

import (
	"time"
	"github.com/msoap/html2data"
)

const UNKNOWN_TITLE = "UNKNOWN_TITLE"

func getTitle(url string) string {
	doc := html2data.FromURL(url)
	if doc.Err != nil {
		return UNKNOWN_TITLE
	}

	title, _ := doc.GetDataSingle("title")
	return title
}

func GetTitle(url string) string {
	ch := make(chan string, 1)
	go func(url string) { ch <- getTitle(url) }(url)
	select {
	case title := <-ch:
		return title
	case <-time.After(10 * time.Second):
		return UNKNOWN_TITLE
	}
}
