package hatena_bookmark

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Tags []string

type Bookmark struct {
	// Timestamp time.Time `json:"timestamp"`
	// 使わないし、unmarshalが面倒なのでひとまずなしで
	User    string `json:"user"`
	Tags    Tags   `json:"tags"`
	Comment string `json:"comment"`
}

type HatenaBookmarks struct {
	Title      string      `json:"title"`
	Bookmarks  []*Bookmark `json:"bookmarks,omitempty"`
	Screenshot string      `json:"screenshot"`
	EntryUrl   string      `json:"entry_url"`
	Count      int         `json:"count"`
	Url        int         `json:"url"`
	EId        int         `json:"eid"`
}

func (bookmarks *HatenaBookmarks) MarshalBinary() ([]byte, error) {
	json, err := json.Marshal(bookmarks)
	if err != nil {
		return nil, err
	}
	return []byte(json), nil
}

func (bookmarks *HatenaBookmarks) UnmarshalBinary(data []byte) error {
	err := json.Unmarshal(data, bookmarks)
	if err != nil {
		return err
	}
	return nil
}

func GetHatenaBookmark(url string) (*HatenaBookmarks, error) {
	res, err := http.Get(fmt.Sprintf("http://b.hatena.ne.jp/entry/jsonlite/?url=%s", url))
	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error: %d", res.StatusCode)
	}

	defer res.Body.Close()
	body, error := ioutil.ReadAll(res.Body)
	if error != nil {
		return nil, err
	}

	bookmarks := HatenaBookmarks{}
	err = json.Unmarshal(body, &bookmarks)
	if error != nil {
		return nil, err
	}
	return &bookmarks, nil
}
