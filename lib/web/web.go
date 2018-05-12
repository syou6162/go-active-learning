package web

import (
	"fmt"
	"net/http"

	"io/ioutil"
	"os"
	"strings"

	"html/template"

	"github.com/codegangsta/cli"
	_ "github.com/lib/pq"
	"github.com/syou6162/go-active-learning/lib/cache"
	"github.com/syou6162/go-active-learning/lib/db"
	"github.com/syou6162/go-active-learning/lib/example"
)

const templateIndexContent = `
<!DOCTYPE html>
<html>
<head>
  <meta charset="UTF-8">
  <title>Recommended entries</title>
</head>
<body>
<h2>General</h2>
<ul>{{range .GeneralList}}
  <li><a href="{{.Url}}">{{or .Title .Url}}</a></li>{{end}}
</ul>
<h2>Twitter</h2>
<ul>{{range .TwitterList}}
  <li><a href="{{.Url}}">{{or .Title .Url}}</a></li>{{end}}
</ul>
<h2>Github</h2>
<ul>{{range .GithubList}}
  <li><a href="{{.Url}}">{{or .Title .Url}}</a></li>{{end}}
</ul>
<h2>Arxiv</h2>
<ul>{{range .ArxivList}}
  <li><a href="{{.Url}}">{{or .Title .Url}}</a></li>{{end}}
</ul>
<h2>Slideshare</h2>
<ul>{{range .SlideShareList}}
  <li><a href="{{.Url}}">{{or .Title .Url}}</a></li>{{end}}
</ul>
<h2>Speaker Deck</h2>
<ul>{{range .SpeakerDeckList}}
  <li><a href="{{.Url}}">{{or .Title .Url}}</a></li>{{end}}
</ul>
</body>
</html>
`

const templateRecentAddedExamplesContent = `
<!DOCTYPE html>
<html>
<head>
  <meta charset="UTF-8">
  <title>Examples recently annotated</title>
</head>
<style>
li {list-style-type: none;}
</style>
<body>
<ul>{{range .}}
  <li><a href="{{.Url}}">{{or .Title .Url}}</a><dd>Label: {{.Label}}</dd></li>{{end}}
</ul>
</body>
</html>
`

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
		err := db.InsertExamplesFromReader(strings.NewReader(string(buf)))
		if err != nil {
			w.WriteHeader(http.StatusBadGateway)
			fmt.Fprintln(w, err.Error())
		}
	}
}

func showRecentAddedExamples(w http.ResponseWriter, r *http.Request) {
	var t *template.Template

	cache, err := cache.NewCache()
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		fmt.Fprintln(w, err.Error())
		return
	}
	defer cache.Close()

	conn, err := db.CreateDBConnection()
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		fmt.Fprintln(w, err.Error())
		return
	}
	defer conn.Close()

	examples, err := db.ReadLabeledExamples(conn, 100)
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		fmt.Fprintln(w, err.Error())
		return
	}
	cache.AttachMetaData(examples)

	t = template.Must(template.New("body").Parse(templateRecentAddedExamplesContent))
	err = t.Execute(w, examples)
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		fmt.Fprintln(w, err.Error())
		return
	}
}

type recommendation struct {
	GeneralList     example.Examples
	TwitterList     example.Examples
	GithubList      example.Examples
	SlideShareList  example.Examples
	ArxivList       example.Examples
	SpeakerDeckList example.Examples
}

func index(w http.ResponseWriter, r *http.Request) {
	var t *template.Template

	cache, err := cache.NewCache()
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		fmt.Fprintln(w, err.Error())
		return
	}
	defer cache.Close()

	conn, err := db.CreateDBConnection()
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		fmt.Fprintln(w, err.Error())
		return
	}
	defer conn.Close()
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		fmt.Fprintln(w, err.Error())
		return
	}

	getUrlsFromList := func(listName string) (example.Examples, error) {
		generalUrls, err := cache.GetUrlsFromList(listName, 0, 100)
		if err != nil {
			return nil, err
		}
		examples, err := db.SearchExamplesByUlrs(conn, generalUrls)
		if err != nil {
			return nil, err
		}
		cache.AttachMetaData(examples)
		return examples, nil
	}

	generalExamples, err := getUrlsFromList("general")
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		fmt.Fprintln(w, err.Error())
		return
	}

	githubExamples, err := getUrlsFromList("github")
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		fmt.Fprintln(w, err.Error())
		return
	}

	slideshareExamples, err := getUrlsFromList("slideshare")
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		fmt.Fprintln(w, err.Error())
		return
	}

	twitterExamples, err := getUrlsFromList("twitter")
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		fmt.Fprintln(w, err.Error())
		return
	}

	arxivExamples, err := getUrlsFromList("arxiv")
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		fmt.Fprintln(w, err.Error())
		return
	}
	speakerdeckExamples, err := getUrlsFromList("speakerdeck")
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		fmt.Fprintln(w, err.Error())
		return
	}

	t = template.Must(template.New("index").Parse(templateIndexContent))
	err = t.Execute(w, recommendation{
		GeneralList:     generalExamples,
		GithubList:      githubExamples,
		SlideShareList:  slideshareExamples,
		TwitterList:     twitterExamples,
		ArxivList:       arxivExamples,
		SpeakerDeckList: speakerdeckExamples,
	})
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		fmt.Fprintln(w, err.Error())
		return
	}
}

func doServe(c *cli.Context) error {
	http.HandleFunc("/", index) // ハンドラを登録してウェブページを表示させる
	http.HandleFunc("/register_training_data", registerTrainingData)
	http.HandleFunc("/show_recent_added_examples", showRecentAddedExamples)
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
