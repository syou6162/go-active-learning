package apply_test

import (
	"testing"

	"github.com/codegangsta/cli"
	"github.com/syou6162/go-active-learning/lib/command"
)

func TestDoApply(t *testing.T) {
	app := cli.NewApp()
	app.Commands = command.Commands
	args := []string{
		"go-active-learning",
		"apply",
		"--input-filename=../../tech_input_example.txt",
		"--filter-status-code-ok",
		"--json-output",
		"--subset-selection",
		"-r=0.75",
		"--size-constraint=20",
		"--score-threshold=0.1",
	}

	if err := app.Run(args); err != nil {
		t.Error(err)
	}
}
