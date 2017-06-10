package main

import (
	"fmt"
	"testing"
)

func TestGetAccuracy(t *testing.T) {
	gold := []LabelType{POSITIVE, POSITIVE, NEGATIVE, NEGATIVE}
	predict := []LabelType{POSITIVE, POSITIVE, NEGATIVE, POSITIVE}
	accuracy := 0.75

	if GetAccuracy(gold, predict) != accuracy {
		t.Error(fmt.Printf("Accuracy should be %f", accuracy))
	}
}

func TestGetPrecision(t *testing.T) {
	gold := []LabelType{POSITIVE, POSITIVE, NEGATIVE, NEGATIVE}
	predict := []LabelType{POSITIVE, NEGATIVE, NEGATIVE, POSITIVE}
	precision := 0.5

	if GetPrecision(gold, predict) != precision {
		t.Error(fmt.Printf("Precision should be %f", precision))
	}
}

func TestGetRecall(t *testing.T) {
	gold := []LabelType{POSITIVE, POSITIVE, NEGATIVE, NEGATIVE}
	predict := []LabelType{POSITIVE, NEGATIVE, NEGATIVE, POSITIVE}
	recall := 0.5

	if GetRecall(gold, predict) != recall {
		t.Error(fmt.Printf("Recall should be %f", recall))
	}
}
