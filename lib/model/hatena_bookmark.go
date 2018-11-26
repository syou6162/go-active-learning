package model

import "encoding/json"

type Tags []string

type Bookmark struct {
	// Timestamp time.Time `json:"timestamp"`
	// 使わないし、unmarshalが面倒なのでひとまずなしで
	User    string `json:"user"`
	Tags    Tags   `json:"tags"`
	Comment string `json:"comment"`
}

type HatenaBookmark struct {
	Title      string      `json:"title"`
	Bookmarks  []*Bookmark `json:"bookmarks"`
	Screenshot string      `json:"screenshot"`
	EntryUrl   string      `json:"entry_url"`
	Count      int         `json:"count"`
	Url        string      `json:"url"`
	EId        string      `json:"eid"`
}

func (bookmarks *HatenaBookmark) MarshalBinary() ([]byte, error) {
	json, err := json.Marshal(bookmarks)
	if err != nil {
		return nil, err
	}
	return []byte(json), nil
}

func (bookmarks *HatenaBookmark) UnmarshalBinary(data []byte) error {
	err := json.Unmarshal(data, bookmarks)
	if err != nil {
		return err
	}
	return nil
}
