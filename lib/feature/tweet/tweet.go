package tweet_feature

import (
	"fmt"
	"regexp"

	"github.com/syou6162/go-active-learning/lib/feature"
	"github.com/syou6162/go-active-learning/lib/model"
	"gopkg.in/vmarkovtsev/go-lcss.v1"
)

type ExampleAndTweet struct {
	example       *model.Example
	tweet         *model.Tweet
	lcsLen        int
	atMarksCnt    int
	hashTagsCnt   int
	cleanedText   string
	cleanedLcsLen int
}

func (et *ExampleAndTweet) GetLabel() model.LabelType {
	return et.tweet.Label
}

func (et *ExampleAndTweet) GetFeatureVector() feature.FeatureVector {
	return et.GetFeatureVector()
}

func GetExampleAndTweet(e *model.Example, t *model.Tweet) ExampleAndTweet {
	result := ExampleAndTweet{example: e, tweet: t}
	result.lcsLen = GetLCSLen(e.Title, t.FullText)

	atRegexp := regexp.MustCompile(`@[^ ]+`)
	result.atMarksCnt = len(atRegexp.FindAllStringSubmatch(t.FullText, -1))
	str := atRegexp.ReplaceAllString(t.FullText, "")
	hashRegexp := regexp.MustCompile(`#[^ ]+`)
	result.hashTagsCnt = len(hashRegexp.FindAllStringSubmatch(t.FullText, -1))
	result.cleanedText = hashRegexp.ReplaceAllString(str, "")
	result.cleanedLcsLen = GetLCSLen(e.Title, result.cleanedText)
	return result
}

func GetLCSLen(str1 string, str2 string) int {
	return len(string(lcss.LongestCommonSubstring([]byte(str1), []byte(str2))))
}

func LCSLenFeature(et ExampleAndTweet) string {
	prefix := "LCSLenFeature"
	len := et.lcsLen
	switch {
	case len == 0:
		return fmt.Sprintf("%s:0", prefix)
	case len < 5:
		return fmt.Sprintf("%s:5", prefix)
	case len < 10:
		return fmt.Sprintf("%s:10", prefix)
	case len < 25:
		return fmt.Sprintf("%s:25", prefix)
	case len < 50:
		return fmt.Sprintf("%s:50", prefix)
	case len < 100:
		return fmt.Sprintf("%s:100", prefix)
	default:
		return fmt.Sprintf("%s:INF", prefix)
	}
}

func CleanedLCSLenFeature(et ExampleAndTweet) string {
	prefix := "CleanedLCSLenFeature"
	len := et.cleanedLcsLen
	switch {
	case len == 0:
		return fmt.Sprintf("%s:0", prefix)
	case len < 5:
		return fmt.Sprintf("%s:5", prefix)
	case len < 10:
		return fmt.Sprintf("%s:10", prefix)
	case len < 25:
		return fmt.Sprintf("%s:25", prefix)
	case len < 50:
		return fmt.Sprintf("%s:50", prefix)
	case len < 100:
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

func CleanedLCSRatioFeature(et ExampleAndTweet) string {
	prefix := "CleanedLCSRatioFeature"
	ratio := float64(et.cleanedLcsLen) / float64(len(et.tweet.FullText))
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
	case cnt <= 3:
		return fmt.Sprintf("%s:3", prefix)
	case cnt <= 5:
		return fmt.Sprintf("%s:5", prefix)
	case cnt <= 10:
		return fmt.Sprintf("%s:10", prefix)
	case cnt <= 25:
		return fmt.Sprintf("%s:25", prefix)
	case cnt <= 50:
		return fmt.Sprintf("%s:50", prefix)
	case cnt <= 100:
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
	case cnt <= 3:
		return fmt.Sprintf("%s:3", prefix)
	case cnt <= 5:
		return fmt.Sprintf("%s:5", prefix)
	case cnt <= 10:
		return fmt.Sprintf("%s:10", prefix)
	case cnt <= 25:
		return fmt.Sprintf("%s:25", prefix)
	case cnt <= 50:
		return fmt.Sprintf("%s:50", prefix)
	case cnt <= 100:
		return fmt.Sprintf("%s:100", prefix)
	default:
		return fmt.Sprintf("%s:INF", prefix)
	}
}

func AtMarksCountFeature(et ExampleAndTweet) string {
	prefix := "AtMarksCountFeature"
	cnt := et.atMarksCnt
	switch {
	case cnt == 0:
		return fmt.Sprintf("%s:0", prefix)
	case cnt == 1:
		return fmt.Sprintf("%s:1", prefix)
	case cnt <= 3:
		return fmt.Sprintf("%s:3", prefix)
	case cnt <= 5:
		return fmt.Sprintf("%s:5", prefix)
	case cnt <= 10:
		return fmt.Sprintf("%s:10", prefix)
	default:
		return fmt.Sprintf("%s:INF", prefix)
	}
}

func HashTagsCountFeature(et ExampleAndTweet) string {
	prefix := "HashTagsCountFeature"
	cnt := et.atMarksCnt
	switch {
	case cnt == 0:
		return fmt.Sprintf("%s:0", prefix)
	case cnt == 1:
		return fmt.Sprintf("%s:1", prefix)
	case cnt <= 3:
		return fmt.Sprintf("%s:3", prefix)
	case cnt <= 5:
		return fmt.Sprintf("%s:5", prefix)
	case cnt <= 10:
		return fmt.Sprintf("%s:10", prefix)
	default:
		return fmt.Sprintf("%s:INF", prefix)
	}
}

func TextLengthFeature(et ExampleAndTweet) string {
	prefix := "TextLengthFeature"
	cnt := len(et.tweet.FullText)
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

func CleanedTextLengthFeature(et ExampleAndTweet) string {
	prefix := "CleanedTextLengthFeature"
	cnt := len(et.cleanedText)
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

func GetTweetFeature(et ExampleAndTweet) feature.FeatureVector {
	var fv feature.FeatureVector

	fv = append(fv, "BIAS")
	fv = append(fv, LCSLenFeature(et))
	fv = append(fv, CleanedLCSLenFeature(et))
	fv = append(fv, LCSRatioFeature(et))
	fv = append(fv, CleanedLCSRatioFeature(et))
	fv = append(fv, TextLengthFeature(et))
	fv = append(fv, CleanedTextLengthFeature(et))

	fv = append(fv, ScreenNameFeature(et))
	fv = append(fv, FavoriteCountFeature(et))
	fv = append(fv, RetweetCountFeature(et))
	fv = append(fv, AtMarksCountFeature(et))
	fv = append(fv, HashTagsCountFeature(et))
	return fv
}
