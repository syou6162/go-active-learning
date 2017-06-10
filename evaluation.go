package main

func GetAccuracy(gold []LabelType, predict []LabelType) float64 {
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

func GetPrecision(gold []LabelType, predict []LabelType) float64 {
	tp := 0.0
	fp := 0.0
	for i, v := range gold {
		if v == POSITIVE && predict[i] == POSITIVE {
			tp += 1.0
		}
		if v == NEGATIVE && predict[i] == POSITIVE {
			fp += 1.0
		}
	}
	return tp / (tp + fp)
}

func GetRecall(gold []LabelType, predict []LabelType) float64 {
	tp := 0.0
	fn := 0.0
	for i, v := range gold {
		if v == POSITIVE && predict[i] == POSITIVE {
			tp += 1.0
		}
		if v == POSITIVE && predict[i] == NEGATIVE {
			fn += 1.0
		}
	}
	return tp / (tp + fn)
}
