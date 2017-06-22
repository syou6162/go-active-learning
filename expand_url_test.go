package main

import (
	"testing"

	"github.com/codegangsta/cli"
)

func TestDoExpandUrl(t *testing.T) {
	app := cli.NewApp()
	app.Commands = Commands
	args := []string{
		"go-active-learning",
		"expand-url",
		"--input-filename=tech_input_example.txt",
	}

	if err := app.Run(args); err != nil {
		t.Error(err)
	}
}
