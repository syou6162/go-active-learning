package model

import (
	"database/sql/driver"
	"encoding/json"
	"strings"
	"time"
)

type Tags []string

type HatenaBookmarkTime struct {
	*time.Time
}

// ref: https://dev.classmethod.jp/go/struct-json/
func (hbt *HatenaBookmarkTime) UnmarshalJSON(data []byte) error {
	t, err := time.Parse("\"2006/01/02 15:04\"", string(data))
	*hbt = HatenaBookmarkTime{&t}
	return err
}

func (hbt HatenaBookmarkTime) MarshalJSON() ([]byte, error) {
	return json.Marshal(hbt.Format("2006/01/02 15:04"))
}

// ref: https://qiita.com/roothybrid7/items/52623bedb45ff0c26a8a
func (hbt *HatenaBookmarkTime) Scan(value interface{}) error {
	v := value.(time.Time)
	hbt.Time = &v
	return nil
}

func (hbt HatenaBookmarkTime) Value() (driver.Value, error) {
	return *hbt.Time, nil
}

func (tags *Tags) Scan(value interface{}) error {
	s := value.(string)
	if s == "" {
		return nil
	}
	v := strings.Split(s, "\t")
	*tags = append(*tags, v...)
	return nil
}

func (tags Tags) Value() (driver.Value, error) {
	return strings.Join(tags, "\t"), nil
}

type Bookmark struct {
	HatenaBookmarkId int                `db:"hatena_bookmark_id"`
	Timestamp        HatenaBookmarkTime `json:"timestamp" db:"timestamp"`
	User             string             `json:"user" db:"user"`
	Tags             Tags               `json:"tags" db:"tags"`
	Comment          string             `json:"comment" db:"comment"`
}

type HatenaBookmark struct {
	Id         int         `db:"id"`
	ExampleId  int         `db:"example_id"`
	Title      string      `json:"title" db:"title"`
	Bookmarks  []*Bookmark `json:"bookmarks"`
	Screenshot string      `json:"screenshot" db:"screenshot"`
	EntryUrl   string      `json:"entry_url" db:"entry_url"`
	Count      int         `json:"count" db:"count"`
	Url        string      `json:"url" db:"url"`
	EId        string      `json:"eid" db:"eid"`
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
