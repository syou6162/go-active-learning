package command

import (
	"github.com/codegangsta/cli"
	"github.com/syou6162/go-active-learning/lib/annotation"
	"github.com/syou6162/go-active-learning/lib/apply"
	"github.com/syou6162/go-active-learning/lib/diagnosis"
	"github.com/syou6162/go-active-learning/lib/expand_url"
)

var Commands = []cli.Command{
	annotation.CommandAnnotate,
	apply.CommandApply,
	expand_url.CommandExpandURL,
	diagnosis.CommandDiagnose,
}
