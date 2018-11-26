package model

import (
	"encoding/json"
	"math"
	"strings"
	"time"

	"github.com/syou6162/go-active-learning/lib/feature"
)

type LabelType int

func (lt *LabelType) MarshalBinary() ([]byte, error) {
	return json.Marshal(lt)
}

func (lt *LabelType) UnmarshalBinary(data []byte) error {
	if err := json.Unmarshal(data, &lt); err != nil {
		return err
	}
	return nil
}

const (
	POSITIVE  LabelType = 1
	NEGATIVE  LabelType = -1
	UNLABELED LabelType = 0
)

type ReferringTweets []string

func (tweets *ReferringTweets) MarshalBinary() ([]byte, error) {
	json, err := json.Marshal(tweets)
	if err != nil {
		return nil, err
	}
	return []byte(json), nil
}

func (tweets *ReferringTweets) UnmarshalBinary(data []byte) error {
	err := json.Unmarshal(data, tweets)
	if err != nil {
		return err
	}
	return nil
}

type Example struct {
	Label           LabelType `json:"Label"`
	Fv              feature.FeatureVector
	Url             string `json:"Url"`
	FinalUrl        string `json:"FinalUrl"`
	Title           string `json:"Title"`
	Description     string `json:"Description"`
	OgDescription   string `json:"OgDescription"`
	OgType          string `json:"OgType"`
	OgImage         string `json:"OgImage"`
	Body            string `json:"Body"`
	Score           float64
	IsNew           bool
	StatusCode      int             `json:"StatusCode"`
	Favicon         string          `json:"Favicon"`
	CreatedAt       time.Time       `json:"CreatedAt"`
	UpdatedAt       time.Time       `json:"UpdatedAt"`
	ReferringTweets ReferringTweets `json:"ReferringTweets"`
	HatenaBookmark  HatenaBookmark  `json:"HatenaBookmark"`
}

type Examples []*Example

func (example *Example) Annotate(label LabelType) {
	example.Label = label
}

func (example *Example) IsLabeled() bool {
	return example.Label != UNLABELED
}

func (example *Example) IsTwitterUrl() bool {
	twitterUrl := "https://twitter.com"
	return strings.Contains(example.Url, twitterUrl) || strings.Contains(example.FinalUrl, twitterUrl)
}

func (example *Example) IsArticle() bool {
	// twitterはarticleと返ってくるが除外
	return example.OgType == "article" && !example.IsTwitterUrl()
}

func (slice Examples) Len() int {
	return len(slice)
}

func (slice Examples) Less(i, j int) bool {
	return math.Abs(slice[i].Score) < math.Abs(slice[j].Score)
}

func (slice Examples) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}
