// yum install webkit2gtk3
// to build install sudo dnf install gtk3-devel webkit2gtk3-devel
// go get github.com/zserge/webview
// go get github.com/atotto/clipboard
// go get github.com/go-vgo/robotgo
// sudo dnf install libxkbcommon-devel libXtst-devel libxkbcommon-x11-devel xorg-x11-xkb-utils-devel libpng-devel xsel xclip

package main

import (
	"bufio"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/NimbleArchitect/webview"
	"github.com/atotto/clipboard"
	_ "github.com/mattn/go-sqlite3"
	"strings"
	"time"
)

const debug = true

var w webview.WebView
var action int
var database *sql.DB

func main() {

	const html = `
	<html><head></head><body>
	Move along nothing to see here
	</body></html>`

	database, _ = opendb()
	// if ok == true {
	// 	//defer database.Close()
	// }
	if database.Ping() != nil {
		fmt.Println("99")
	}

	searchandpaste()
	database.Close()
}

func searchandpaste() {
	w = webview.New(debug)
	defer w.Destroy()
	w.SetTitle("snip search")
	w.SetSize(600, 400, webview.HintNone)
	//w.Navigate("data:text/html," + html)
	w.Navigate("file:///home/rich/data/src/go/src/snippets/frontpage.html")
	w.Bind("searchsnip", searchsnip)
	w.Bind("toclipboard", copysnip)
	w.Bind("writesnip", writesnip)
	w.Bind("closesnip", closesnip)
	w.Bind("savesnip", savesnip)
	w.Run()
}

func searchsnip(data string) string {
	if len(data) <= 0 {
		return ""
	}
	//time.Sleep(4 * time.Second)
	//println("running from js: " + data)

	snips := dbfind("name", data)
	for _, itm := range snips {
		itm.Text = ""
	}
	str, _ := json.Marshal(snips)
	//fmt.Println("json: " + string(str))
	return string(str)
}

func copysnip(data string) error {
	clipboard.WriteAll(data)
	return nil
}

func writesnip(hash string) error {
	var code []string
	//fmt.Println("** " + hash)
	snips := dbgetID(hash)
	data := snips.Text

	scanner := bufio.NewScanner(strings.NewReader(data))

	for scanner.Scan() {
		singleline := scanner.Text()
		code = append(code, singleline)
	}

	go typeSnippet(code)
	return nil
}

func closesnip() error {
	go w.Terminate()
	return nil
}

func opendb() (*sql.DB, bool) {
	db, err := sql.Open("sqlite3", "/home/rich/data/src/go/src/snippets/snippets.db")
	if err != nil {
		fmt.Println("ERROR opening database")
	}
	//ok := db.Ping()
	//panic(err)
	statement, _ := db.Prepare("CREATE TABLE IF NOT EXISTS snips (id INTEGER PRIMARY KEY, created INTEGER, name TEXT, code TEXT)")
	statement.Exec()

	return db, true
}

func savesnip(title string, code string) {
	tx, _ := database.Begin()
	stmt, _ := tx.Prepare("insert into snips (created,name,code) values (?,?,?)")
	_, err := stmt.Exec(time.Now(), title, code)
	if err != nil {
		fmt.Print("error saving")
	}
	tx.Commit()
}

type Snipitem struct {
	ID   int       `json:"hash"`
	Time time.Time `json:"time,omitempty"`
	Name string    `json:"name,omitempty"`
	Text string    `json:"text,omitempty"`
}

//type Snips []Snipitem
func dbgetID(hash string) Snipitem {

	qry := string("SELECT * FROM snips WHERE ID = " + hash)
	rows, err := database.Query(qry)
	if err != nil {
		fmt.Println("ERROR: unable to query db")
		panic(err)
	}

	var snip Snipitem
	var id int
	var name string
	var code string
	var created string

	for rows.Next() {
		err = rows.Scan(&id, &created, &name, &code)
		if err != nil {
			panic(err)
		}

		snip = Snipitem{
			ID:   id,
			Time: time.Now(),
			Name: name,
			Text: code,
		}
	}

	rows.Close() //good habit to close
	return snip
}

func dbfind(field string, searchfor string) []Snipitem {

	// query
	qry := string("SELECT * FROM snips WHERE " + field + " LIKE '%" + searchfor + "%'")
	rows, err := database.Query(qry)

	if err != nil {
		fmt.Println("ERROR: unable to query db")
		panic(err)
	}

	var snip []Snipitem
	var id int
	var name string
	var code string
	var created string

	for rows.Next() {
		err = rows.Scan(&id, &created, &name, &code)
		if err != nil {
			panic(err)
		}

		snipitem := Snipitem{
			ID:   id,
			Time: time.Now(),
			Name: name,
			Text: code,
		}

		snip = append(snip, snipitem)
	}

	rows.Close() //good habit to close
	return snip
}
