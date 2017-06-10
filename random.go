package main

import "math/rand"

func shuffle(examples Examples) {
	n := len(examples)
	for i := n - 1; i >= 0; i-- {
		j := rand.Intn(i + 1)
		examples[i], examples[j] = examples[j], examples[i]
	}
}

func RandomSelectOneExample(examples Examples) *Example {
	shuffle(examples)
	return findFirstUnlabeledExample(examples)
}
