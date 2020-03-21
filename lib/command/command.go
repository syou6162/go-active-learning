package command

import (
	"github.com/syou6162/go-active-learning/lib/add"
	"github.com/syou6162/go-active-learning/lib/annotation"
	"github.com/syou6162/go-active-learning/lib/diagnosis"
	"github.com/syou6162/go-active-learning/lib/related_example"
	"github.com/syou6162/go-active-learning/lib/top_accessed_example"
	"github.com/urfave/cli"
)

var Commands = []cli.Command{
	add.CommandAdd,
	related_example.CommandAddRelatedExamples,
	annotation.CommandAnnotate,
	top_accessed_example.CommandAddTopAccessedExamples,
	diagnosis.CommandDiagnose,
}
