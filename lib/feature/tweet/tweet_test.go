package tweet_feature

import (
	"reflect"
	"testing"

	"github.com/syou6162/go-active-learning/lib/feature"
	"github.com/syou6162/go-active-learning/lib/model"
)

func TestExtractHostFeature(t *testing.T) {
	e := model.Example{}
	e.Title = "Hello world"
	tweet := model.Tweet{}
	tweet.ScreenName = "syou6162"
	tweet.FullText = "Hello world @syou6162 @syou6163 #hashtag1 #hashtag2"
	tweet.FavoriteCount = 7
	tweet.RetweetCount = 7

	et := GetExampleAndTweet(&e, &tweet)
	fv := et.GetFeatureVector()
	expect := feature.FeatureVector{
		"BIAS",
		"LCSLenFeature:25",
		"CleanedLCSLenFeature:25",
		"LCSRatioFeature:0.25",
		"CleanedLCSRatioFeature:0.25",
		"TextLengthFeature:100",
		"CleanedTextLengthFeature:25",
		"ScreenNameFeature:syou6162",
		"FavoriteCountFeature:10",
		"RetweetCountFeature:10",
		"AtMarksCountFeature:3",
		"HashTagsCountFeature:3",
	}
	if !reflect.DeepEqual(expect, fv) {
		t.Error("feature must be wrong")
	}
}
