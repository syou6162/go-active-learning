package model

import "fmt"

type RecommendationListType int

const (
	GENERAL RecommendationListType = 0
	ARTICLE RecommendationListType = 1
	GITHUB  RecommendationListType = 2
	SLIDE   RecommendationListType = 3
	ARXIV   RecommendationListType = 4
)

func GetRecommendationListType(listname string) (RecommendationListType, error) {
	switch listname {
	case "general":
		return GENERAL, nil
	case "article":
		return ARTICLE, nil
	case "github":
		return GITHUB, nil
	case "slide":
		return SLIDE, nil
	case "arxiv":
		return ARXIV, nil
	default:
		return -1, fmt.Errorf("no such RecommendationListType for '%s'", listname)
	}
}

type Recommendation struct {
	RecommendationListType RecommendationListType
	ExampleIds             []int
}
