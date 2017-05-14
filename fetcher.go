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

}
