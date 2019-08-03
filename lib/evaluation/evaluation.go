package evaluation

import (
	"github.com/syou6162/go-active-learning/lib/model"
)

func GetAccuracy(gold []model.LabelType, predict []model.LabelType) float64 {
	if len(gold) != len(predict) {
		return 0.0
	}
	sum := 0.0
	for i, v := range gold {
		if v == predict[i] {
			sum += 1.0
		}
	}
	return sum / float64(len(gold))
}

func GetPrecision(gold []model.LabelType, predict []model.LabelType) float64 {
	tp := 0.0
	fp := 0.0
	for i, v := range gold {
		if v == model.POSITIVE && predict[i] == model.POSITIVE {
			tp += 1.0
		}
		if v == model.NEGATIVE && predict[i] == model.POSITIVE {
			fp += 1.0
		}
	}
	return tp / (tp + fp)
}

func GetRecall(gold []model.LabelType, predict []model.LabelType) float64 {
	tp := 0.0
	fn := 0.0
	for i, v := range gold {
		if v == model.POSITIVE && predict[i] == model.POSITIVE {
			tp += 1.0
		}
		if v == model.POSITIVE && predict[i] == model.NEGATIVE {
			fn += 1.0
		}
	}
	return tp / (tp + fn)
}

func GetConfusionMatrix(gold []model.LabelType, predict []model.LabelType) (int, int, int, int) {
	tp := 0
	fp := 0
	fn := 0
	tn := 0
	for i, v := range gold {
		if v == model.POSITIVE && predict[i] == model.POSITIVE {
			tp += 1
		}
		if v == model.NEGATIVE && predict[i] == model.POSITIVE {
			fp += 1
		}
		if v == model.POSITIVE && predict[i] == model.NEGATIVE {
			fn += 1
		}
		if v == model.NEGATIVE && predict[i] == model.NEGATIVE {
			tn += 1
		}
	}
	return tp, fp, fn, tn
}
