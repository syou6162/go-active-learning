package add_test

import (
	"testing"

	"github.com/syou6162/go-active-learning/lib/command"
	"github.com/urfave/cli"
)

func TestDoAdd(t *testing.T) {
	app := cli.NewApp()
	app.Commands = command.Commands
	args := []string{
		"go-active-learning-web",
		"add",
		"--input-filename=../../tech_input_example.txt",
	}

	if err := app.Run(args); err != nil {
		t.Error(err)
	}
}
