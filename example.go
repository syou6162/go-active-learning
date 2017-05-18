package main

import "math"

type LabelType int

const (
	POSITIVE  LabelType = 1
	NEGATIVE  LabelType = -1
	UNLABELED LabelType = 0
)

type Example struct {
	Label       LabelType `json:"Label"`
	Fv          FeatureVector
	Url         string `json:"Url"`
	Title       string `json:"Title"`
	Description string `json:"Description"`
	Body        string `json:"Body"`
	RawHTML     string `json:"RawHTML"`
	Score       float64
	IsNew       bool
	StatusCode  int    `json:StatusCode`
}

type Examples []*Example

func NewExample(url string, label LabelType) *Example {
	IsNew := false
	if label == UNLABELED {
		IsNew = true
	}
	return &Example{label, []string{}, url, "", "", "", "", 0.0, IsNew, 0}
}

func (example *Example) Annotate(label LabelType) {
	example.Label = label
}

func (example *Example) IsLabeled() bool {
	return example.Label != UNLABELED
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
