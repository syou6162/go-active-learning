package main

type LabelType int

const (
	POSITIVE LabelType = iota
	NEGATIVE
	UNLABELED
)

type Example struct {
	label LabelType
	fv    FeatureVector
	url   string
}

func NewExample(url string, label LabelType) *Example {
	return &Example{label, []string{}, url}
}
