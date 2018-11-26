package evaluation

import (
	"fmt"
	"testing"

	"github.com/syou6162/go-active-learning/lib/model"
)

func TestGetAccuracy(t *testing.T) {
	gold := []model.LabelType{model.POSITIVE, model.POSITIVE, model.NEGATIVE, model.NEGATIVE}
	predict := []model.LabelType{model.POSITIVE, model.POSITIVE, model.NEGATIVE, model.POSITIVE}
	accuracy := 0.75

	if GetAccuracy(gold, predict) != accuracy {
		t.Error(fmt.Printf("Accuracy should be %f", accuracy))
	}
}

func TestGetPrecision(t *testing.T) {
	gold := []model.LabelType{model.POSITIVE, model.POSITIVE, model.NEGATIVE, model.NEGATIVE}
	predict := []model.LabelType{model.POSITIVE, model.NEGATIVE, model.NEGATIVE, model.POSITIVE}
	precision := 0.5

	if GetPrecision(gold, predict) != precision {
		t.Error(fmt.Printf("Precision should be %f", precision))
	}
}

func TestGetRecall(t *testing.T) {
	gold := []model.LabelType{model.POSITIVE, model.POSITIVE, model.NEGATIVE, model.NEGATIVE}
	predict := []model.LabelType{model.POSITIVE, model.NEGATIVE, model.NEGATIVE, model.POSITIVE}
	recall := 0.5

	if GetRecall(gold, predict) != recall {
		t.Error(fmt.Printf("Recall should be %f", recall))
	}
}
