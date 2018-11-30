package example

import (
	"time"

	"github.com/syou6162/go-active-learning/lib/feature"
	"github.com/syou6162/go-active-learning/lib/model"
)

func NewExample(url string, label model.LabelType) *model.Example {
	IsNew := false
	if label == model.UNLABELED {
		IsNew = true
	}
	now := time.Now()
	return &model.Example{
		Label:           label,
		Fv:              feature.FeatureVector{},
		Url:             url,
		FinalUrl:        url,
		Title:           "",
		Description:     "",
		OgDescription:   "",
		OgType:          "",
		OgImage:         "",
		Body:            "",
		Score:           0.0,
		IsNew:           IsNew,
		StatusCode:      0,
		Favicon:         "",
		CreatedAt:       now,
		UpdatedAt:       now,
		ReferringTweets: &model.ReferringTweets{},
		HatenaBookmark:  &model.HatenaBookmark{Bookmarks: make([]*model.Bookmark, 0)},
	}
}

func GetStat(examples model.Examples) map[string]int {
	stat := make(map[string]int)
	for _, e := range examples {
		switch e.Label {
		case model.POSITIVE:
			stat["positive"]++
		case model.NEGATIVE:
			stat["negative"]++
		case model.UNLABELED:
			stat["unlabeled"]++
		}
	}
	return stat
}

func ExtractFeatures(e model.Example) feature.FeatureVector {
	var fv feature.FeatureVector
	fv = append(fv, "BIAS")
	fv = append(fv, feature.ExtractHostFeature(e.FinalUrl))
	fv = append(fv, feature.ExtractJpnNounFeatures(feature.ExtractPath(e.FinalUrl), "URL")...)
	fv = append(fv, feature.ExtractNounFeatures(e.Title, "TITLE")...)
	fv = append(fv, feature.ExtractNounFeatures(e.Description, "DESCRIPTION")...)
	fv = append(fv, feature.ExtractNounFeatures(e.Body, "BODY")...)
	return fv
}
