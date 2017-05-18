package main

import (
	"os"

	"github.com/codegangsta/cli"
)

var Commands = []cli.Command{
	commandAnnotate,
	commandDiagnose,
}

func main() {
	app := cli.NewApp()
	app.Name = "go-active-learning"
	app.Commands = Commands

	app.Run(os.Args)
}
