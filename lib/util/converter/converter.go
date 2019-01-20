package converter

import "github.com/syou6162/go-active-learning/lib/model"
import "github.com/syou6162/go-active-learning/lib/classifier"

func ConvertExamplesToLearningInstances(examples model.Examples) classifier.LearningInstances {
	instances := classifier.LearningInstances{}
	for _, e := range examples {
		instances = append(instances, e)
	}
	return instances
}
