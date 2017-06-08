package main

import (
	"fmt"
	"os"

	"github.com/codegangsta/cli"
)

var Commands = []cli.Command{
	commandAnnotate,
	commandApply,
	commandExpandURL,
	commandDiagnose,
}

func main() {
	app := cli.NewApp()
	app.Name = "go-active-learning"
	app.Commands = Commands

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}
