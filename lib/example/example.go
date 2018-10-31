package example

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

type Example struct {
	Label         LabelType `json:"Label"`
	Fv            feature.FeatureVector
	Url           string `json:"Url"`
	FinalUrl      string `json:"FinalUrl"`
	Title         string `json:"Title"`
	Description   string `json:"Description"`
	OgDescription string `json:"OgDescription"`
	OgType        string `json:"OgType"`
	OgImage       string `json:"OgImage"`
	Body          string `json:"Body"`
	Score         float64
	IsNew         bool
	StatusCode    int       `json:"StatusCode"`
	Favicon       string    `json:"Favicon"`
	CreatedAt     time.Time `json:"CreatedAt"`
	UpdatedAt     time.Time `json:"UpdatedAt"`
}

type Examples []*Example

func NewExample(url string, label LabelType) *Example {
	IsNew := false
	if label == UNLABELED {
		IsNew = true
	}
	now := time.Now()
	return &Example{
		Label:         label,
		Fv:            []string{},
		Url:           url,
		FinalUrl:      url,
		Title:         "",
		Description:   "",
		OgDescription: "",
		OgType:        "",
		OgImage:       "",
		Body:          "",
		Score:         0.0,
		IsNew:         IsNew,
		StatusCode:    0,
		Favicon:       "",
		CreatedAt:     now,
		UpdatedAt:     now,
	}
}

func (example *Example) Annotate(label LabelType) {
	example.Label = label
}

func (example *Example) IsLabeled() bool {
	return example.Label != UNLABELED
}

func (example *Example) isTwitterUrl() bool {
	twitterUrl := "https://twitter.com"
	return strings.Contains(example.Url, twitterUrl) || strings.Contains(example.FinalUrl, twitterUrl)
}

func (example *Example) IsArticle() bool {
	// twitterはarticleと返ってくるが除外
	return example.OgType == "article" && !example.isTwitterUrl()
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

func GetStat(examples Examples) map[string]int {
	stat := make(map[string]int)
	for _, e := range examples {
		switch e.Label {
		case POSITIVE:
			stat["positive"]++
		case NEGATIVE:
			stat["negative"]++
		case UNLABELED:
			stat["unlabeled"]++
		}
	}
	return stat
}

func ExtractFeatures(e Example) feature.FeatureVector {
	var fv feature.FeatureVector
	fv = append(fv, "BIAS")
	fv = append(fv, feature.ExtractHostFeature(e.FinalUrl))
	fv = append(fv, feature.ExtractJpnNounFeatures(feature.ExtractPath(e.FinalUrl), "URL")...)
	fv = append(fv, feature.ExtractNounFeatures(e.Title, "TITLE")...)
	fv = append(fv, feature.ExtractNounFeatures(e.Description, "DESCRIPTION")...)
	fv = append(fv, feature.ExtractNounFeatures(e.Body, "BODY")...)
	return fv
}
