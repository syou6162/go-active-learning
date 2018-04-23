package web

import (
	"fmt"
	"net/http"

	"bufio"
	"io/ioutil"
	"os"
	"strings"

	"github.com/codegangsta/cli"
	_ "github.com/lib/pq"
	"github.com/syou6162/go-active-learning/lib/db"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, World")
}

func checkAuth(r *http.Request) bool {
	username, password, ok := r.BasicAuth()
	if ok == false {
		return false
	}
	return username == os.Getenv("BASIC_AUTH_USERNAME") && password == os.Getenv("BASIC_AUTH_PASSWORD")
}

func registerTrainingData(w http.ResponseWriter, r *http.Request) {
	if checkAuth(r) == false {
		w.WriteHeader(401)
		w.Write([]byte("401 Unauthorized\n"))
		return
	} else {
		buf, _ := ioutil.ReadAll(r.Body)
		scanner := bufio.NewScanner(strings.NewReader(string(buf)))

		conn, err := db.CreateDBConnection()
		if err != nil {
			fmt.Fprintln(w, err.Error())
		}
		defer conn.Close()

		for scanner.Scan() {
			_, err := db.InsertExampleFromScanner(conn, scanner)
			if err != nil {
				fmt.Fprintln(w, err.Error())
			}
		}
		if err := scanner.Err(); err != nil {
			fmt.Fprintln(w, err.Error())
		}
	}
}

func doServe(c *cli.Context) error {
	http.HandleFunc("/", handler) // ハンドラを登録してウェブページを表示させる
	http.HandleFunc("/register_training_data", registerTrainingData)
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
