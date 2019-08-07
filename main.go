package main

import (
	"fmt"
	"os"

	"github.com/syou6162/go-active-learning/lib/command"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "go-active-learning"
	app.Commands = command.Commands

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}
