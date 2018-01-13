package fetcher

import (
	"fmt"
	"testing"
)

func TestGetArticle(t *testing.T) {
	a := GetArticle("http://www.yasuhisay.info/entry/20090516/1242480413")
	fmt.Println(a)

	if a.Title == "" {
		t.Error("Title must not be empty")
	}
	if a.Description == "" {
		t.Error("Description must not be empty")
	}
	if a.StatusCode != 200 {
		t.Error("StatusCode must be 200")
	}
}
