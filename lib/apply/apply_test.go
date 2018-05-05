package apply_test

import (
	"testing"

	"fmt"

	"github.com/codegangsta/cli"
	"github.com/syou6162/go-active-learning/lib/command"
	"github.com/syou6162/go-active-learning/lib/db"
	"github.com/syou6162/go-active-learning/lib/util/file"
)

func TestDoApply(t *testing.T) {
	inputFilename := "../../tech_input_example.txt"
	train, err := file.ReadExamples(inputFilename)
	if err != nil {
		t.Error(err)
	}

	conn, err := db.CreateDBConnection()
	if err != nil {
		t.Error(err)
	}
	defer conn.Close()

	_, err = db.DeleteAllExamples(conn)
	if err != nil {
		t.Error(err)
	}

	for _, example := range train {
		_, err = db.InsertOrUpdateExample(conn, example)
		if err != nil {
			t.Error(err)
		}
	}

	app := cli.NewApp()
	app.Commands = command.Commands
	args := []string{
		"go-active-learning",
		"apply",
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
