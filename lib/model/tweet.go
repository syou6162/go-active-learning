package model

import (
	"time"
)

type Tweet struct {
	Id        int `db:"id"`
	ExampleId int `db:"example_id"`

	CreatedAt     time.Time `json:"CreatedAt" db:"created_at"`
	IdStr         string    `json:"IdStr" db:"id_str"`
	FullText      string    `json:"FullText" db:"full_text"`
	FavoriteCount int       `json:"FavoriteCount" db:"favorite_count"`
	RetweetCount  int       `json:"RetweetCount" db:"retweet_count"`
	Lang          string    `json:"Lang" db:"lang"`

	ScreenName      string    `json:"ScreenName" db:"screen_name"`
	Name            string    `json:"Name" db:"name"`
	ProfileImageUrl string    `json:"ProfileImageUrl" db:"profile_image_url"`
	Label           LabelType `json:"Label" db:"label"`
}

type ReferringTweets []*Tweet
