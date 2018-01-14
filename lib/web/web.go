package web

import (
	"fmt"
	"net/http"

	"github.com/codegangsta/cli"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, World")
}

func doServe(c *cli.Context) error {
	http.HandleFunc("/", handler) // ハンドラを登録してウェブページを表示させる
	return http.ListenAndServe(":7777", nil)
}

var CommandServe = cli.Command{
	Name:  "serve",
	Usage: "Run a server",
	Description: `
Run a web server.
`,
	Action: doServe,
	Flags:  []cli.Flag{},
}
