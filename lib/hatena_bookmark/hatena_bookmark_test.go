package hatena_bookmark

import (
	"testing"
)

func TestGetHatenaBookmark(t *testing.T) {
	bookmarks, err := GetHatenaBookmark("https://www.yasuhisay.info")
	if err != nil {
		t.Error(err.Error())
	}

	if bookmarks.Title == "" {
		t.Error("Title must not be empty")
	}
	if bookmarks.Count == 0 {
		t.Error("Count must not be 0")
	}
	if len(bookmarks.Bookmarks) == 0 {
		t.Error("Count must not be 0")
	}
}
