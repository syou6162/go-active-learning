package main

import (
	"github.com/syou6162/go-active-learning/lib/example"
	"math/rand"
)

func shuffle(examples example.Examples) {
	n := len(examples)
	for i := n - 1; i >= 0; i-- {
		j := rand.Intn(i + 1)
		examples[i], examples[j] = examples[j], examples[i]
	}
}
