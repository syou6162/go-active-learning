package main

type LabelType int

const (
	POSITIVE  LabelType = 1
	NEGATIVE  LabelType = -1
	UNLABELED LabelType = 0
)

type Example struct {
	label LabelType
	fv    FeatureVector
	url   string
	title string
	score float64
}

type Examples []*Example

func NewExample(url string, label LabelType) *Example {
	return &Example{label, []string{}, url, "", 0.0}
}

func (example *Example) Annotate(label LabelType) {
	example.label = label
}

func (example *Example) IsLabeled() bool {
	return example.label != UNLABELED
}
