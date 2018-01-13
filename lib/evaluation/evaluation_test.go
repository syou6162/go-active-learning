package evaluation

import (
	"fmt"
	"github.com/syou6162/go-active-learning/lib/example"
	"testing"
)

func TestGetAccuracy(t *testing.T) {
	gold := []example.LabelType{example.POSITIVE, example.POSITIVE, example.NEGATIVE, example.NEGATIVE}
	predict := []example.LabelType{example.POSITIVE, example.POSITIVE, example.NEGATIVE, example.POSITIVE}
	accuracy := 0.75

	if GetAccuracy(gold, predict) != accuracy {
		t.Error(fmt.Printf("Accuracy should be %f", accuracy))
	}
}

func TestGetPrecision(t *testing.T) {
	gold := []example.LabelType{example.POSITIVE, example.POSITIVE, example.NEGATIVE, example.NEGATIVE}
	predict := []example.LabelType{example.POSITIVE, example.NEGATIVE, example.NEGATIVE, example.POSITIVE}
	precision := 0.5

	if GetPrecision(gold, predict) != precision {
		t.Error(fmt.Printf("Precision should be %f", precision))
	}
}

func TestGetRecall(t *testing.T) {
	gold := []example.LabelType{example.POSITIVE, example.POSITIVE, example.NEGATIVE, example.NEGATIVE}
	predict := []example.LabelType{example.POSITIVE, example.NEGATIVE, example.NEGATIVE, example.POSITIVE}
	recall := 0.5

	if GetRecall(gold, predict) != recall {
		t.Error(fmt.Printf("Recall should be %f", recall))
	}
}
