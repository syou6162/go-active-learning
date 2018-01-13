package evaluation

import (
	"github.com/syou6162/go-active-learning/lib/example"
)

func GetAccuracy(gold []example.LabelType, predict []example.LabelType) float64 {
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

func GetPrecision(gold []example.LabelType, predict []example.LabelType) float64 {
	tp := 0.0
	fp := 0.0
	for i, v := range gold {
		if v == example.POSITIVE && predict[i] == example.POSITIVE {
			tp += 1.0
		}
		if v == example.NEGATIVE && predict[i] == example.POSITIVE {
			fp += 1.0
		}
	}
	return tp / (tp + fp)
}

func GetRecall(gold []example.LabelType, predict []example.LabelType) float64 {
	tp := 0.0
	fn := 0.0
	for i, v := range gold {
		if v == example.POSITIVE && predict[i] == example.POSITIVE {
			tp += 1.0
		}
		if v == example.POSITIVE && predict[i] == example.NEGATIVE {
			fn += 1.0
		}
	}
	return tp / (tp + fn)
}
