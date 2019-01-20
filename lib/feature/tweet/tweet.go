package tweet_feature

import (
	"fmt"

	"github.com/syou6162/go-active-learning/lib/feature"
	"github.com/syou6162/go-active-learning/lib/model"
	"gopkg.in/vmarkovtsev/go-lcss.v1"
)

type ExampleAndTweet struct {
	example *model.Example
	tweet   *model.Tweet
	lcsLen  int
}

func GetExampleAndTweet(e *model.Example, t *model.Tweet) ExampleAndTweet {
	result := ExampleAndTweet{example: e, tweet: t}
	result.lcsLen = GetLCSLen(result)
	return result
}

func GetLCSLen(et ExampleAndTweet) int {
	return len(string(lcss.LongestCommonSubstring([]byte(et.example.Title), []byte(et.tweet.FullText))))
}

func LCSLenFeature(et ExampleAndTweet) string {
	prefix := "LCSLenFeature"
	switch {
	case et.lcsLen == 0:
		return fmt.Sprintf("%s:0", prefix)
	case et.lcsLen < 5:
		return fmt.Sprintf("%s:5", prefix)
	case et.lcsLen < 10:
		return fmt.Sprintf("%s:10", prefix)
	case et.lcsLen < 25:
		return fmt.Sprintf("%s:25", prefix)
	case et.lcsLen < 50:
		return fmt.Sprintf("%s:50", prefix)
	case et.lcsLen < 100:
		return fmt.Sprintf("%s:100", prefix)
	default:
		return fmt.Sprintf("%s:INF", prefix)
	}
}

func LCSRatioFeature(et ExampleAndTweet) string {
	prefix := "LCSRatioFeature"
	ratio := float64(et.lcsLen) / float64(len(et.tweet.FullText))
	switch {
	case ratio == 0.0:
		return fmt.Sprintf("%s:0.0", prefix)
	case ratio < 0.1:
		return fmt.Sprintf("%s:0.1", prefix)
	case ratio < 0.25:
		return fmt.Sprintf("%s:0.25", prefix)
	case ratio < 0.5:
		return fmt.Sprintf("%s:0.5", prefix)
	case ratio < 0.75:
		return fmt.Sprintf("%s:0.75", prefix)
	case ratio < 0.9:
		return fmt.Sprintf("%s:0.0", prefix)
	default:
		return fmt.Sprintf("%s:1.0", prefix)
	}
}

func FavoriteCountFeature(et ExampleAndTweet) string {
	prefix := "FavoriteCountFeature"
	cnt := et.tweet.FavoriteCount
	switch {
	case cnt == 0:
		return fmt.Sprintf("%s:0", prefix)
	case cnt == 1:
		return fmt.Sprintf("%s:1", prefix)
	case cnt == 3:
		return fmt.Sprintf("%s:3", prefix)
	case cnt < 5:
		return fmt.Sprintf("%s:5", prefix)
	case cnt < 10:
		return fmt.Sprintf("%s:10", prefix)
	case cnt < 25:
		return fmt.Sprintf("%s:25", prefix)
	case cnt < 50:
		return fmt.Sprintf("%s:50", prefix)
	case cnt < 100:
		return fmt.Sprintf("%s:100", prefix)
	default:
		return fmt.Sprintf("%s:INF", prefix)
	}
}

func RetweetCountFeature(et ExampleAndTweet) string {
	prefix := "RetweetCountFeature"
	cnt := et.tweet.RetweetCount
	switch {
	case cnt == 0:
		return fmt.Sprintf("%s:0", prefix)
	case cnt == 1:
		return fmt.Sprintf("%s:1", prefix)
	case cnt == 3:
		return fmt.Sprintf("%s:3", prefix)
	case cnt < 5:
		return fmt.Sprintf("%s:5", prefix)
	case cnt < 10:
		return fmt.Sprintf("%s:10", prefix)
	case cnt < 25:
		return fmt.Sprintf("%s:25", prefix)
	case cnt < 50:
		return fmt.Sprintf("%s:50", prefix)
	case cnt < 100:
		return fmt.Sprintf("%s:100", prefix)
	default:
		return fmt.Sprintf("%s:INF", prefix)
	}
}

func ScreenNameFeature(et ExampleAndTweet) string {
	prefix := "ScreenNameFeature"
	return fmt.Sprintf("%s:%s", prefix, et.tweet.ScreenName)
}

func GetTweetFeature(e *model.Example, t *model.Tweet) feature.FeatureVector {
	var fv feature.FeatureVector
	et := GetExampleAndTweet(e, t)

	fv = append(fv, "BIAS")
	fv = append(fv, LCSLenFeature(et))
	fv = append(fv, LCSRatioFeature(et))
	fv = append(fv, ScreenNameFeature(et))
	fv = append(fv, FavoriteCountFeature(et))
	fv = append(fv, RetweetCountFeature(et))
	return fv
}
