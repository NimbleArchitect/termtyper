// yum install webkit2gtk3
// to build install sudo dnf install gtk3-devel webkit2gtk3-devel
// go get github.com/zserge/webview
// go get github.com/atotto/clipboard
// go get github.com/go-vgo/robotgo
// sudo dnf install libxkbcommon-devel libXtst-devel libxkbcommon-x11-devel xorg-x11-xkb-utils-devel libpng-devel xsel xclip

package main

import (
	"database/sql"
	"fmt"
	"github.com/NimbleArchitect/webview"
	_ "github.com/mattn/go-sqlite3"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

const debug = true

type Snipitem struct {
	ID       int        `json:"hash"`
	Time     time.Time  `json:"time,omitempty"`
	Name     string     `json:"name,omitempty"`
	Code     string     `json:"code,omitempty"`
	Argument []SnipArgs `json:"argument,omitempty"`
}

type SnipArgs struct {
	Name  string `json:"name,omitempty"`
	Value string `json:"value,omitempty"`
}

var w webview.WebView
var action int
var database *sql.DB
var datapath string

func main() {

	const html = `
	<html><head></head><body>
	Move along nothing to see here
	</body></html>`

	datapath, err := os.UserHomeDir()
	if err != nil {
		panic("Unable to get users profile folder")
	}

	if _, err := os.Stat(datapath + "/.snippets"); err != nil {
		err = os.Mkdir(datapath+"/.snippets", 0770)
		if err != nil {
			panic("unable to create folder ~/.snippets")
		}
	}

	database, _ = opendb(datapath + "/.snippets/snippets.db")
	// if ok == true {
	// 	//defer database.Close()
	// }
	if database.Ping() != nil {
		fmt.Println("99")
	}
	execpath := getprogPath()
	searchandpaste(execpath)
	database.Close()
}

func getprogPath() string {
	var dirAbsPath string
	ex, err := os.Executable()
	if err == nil {
		dirAbsPath = filepath.Dir(ex)
		fmt.Println(dirAbsPath)
		return dirAbsPath
	}

	exReal, err := filepath.EvalSymlinks(ex)
	if err != nil {
		panic(err)
	}
	dirAbsPath = filepath.Dir(exReal)
	fmt.Println(dirAbsPath)
	return dirAbsPath
}

func searchandpaste(datapath string) {
	w = webview.New(debug)
	defer w.Destroy()
	w.SetTitle("snip search")
	w.SetSize(600, 400, webview.HintNone)
	//w.Navigate("data:text/html," + html)
	w.Navigate("file://" + datapath + "/frontpage.html")
	w.Bind("snipSearch", snip_search)
	w.Bind("toclipboard", snip_copy)
	w.Bind("snipWrite", snip_write)
	w.Bind("snipClose", snip_close)
	w.Bind("snipSave", snip_save)
	//w.Bind("snipGetVarList", snip_getvars)
	w.Run()
}

func opendb(dbpath string) (*sql.DB, bool) {
	fmt.Println("* open: " + dbpath)
	db, err := sql.Open("sqlite3", dbpath)

	if err != nil {
		fmt.Println("ERROR opening database")
	}
	ok := db.Ping()
	if ok != nil {
		panic(err)
	}
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS snips (id INTEGER PRIMARY KEY, created INTEGER, name TEXT, code TEXT)")

	return db, true
}

func dbgetID(hash string) Snipitem {
	var snip Snipitem
	var id int
	var name string
	var code string
	var created string

	qry := string("SELECT * FROM snips WHERE ID = " + hash)
	rows, err := database.Query(qry)
	if err != nil {
		fmt.Println("ERROR: unable to query db")
		panic(err)
	}

	for rows.Next() {
		err = rows.Scan(&id, &created, &name, &code)
		if err != nil {
			panic(err)
		}
		//tags := len(getVars(code))
		snip = Snipitem{
			ID:   id,
			Time: time.Now(),
			Name: name,
			Code: code,
		}
	}

	rows.Close() //good habit to close
	return snip
}

func dbfind(field string, searchfor string) []Snipitem {
	var snip []Snipitem
	var id int
	var name string
	var code string
	var created string

	// query
	qry := string("SELECT * FROM snips WHERE " + field + " LIKE '%" + searchfor + "%'")
	rows, err := database.Query(qry)

	if err != nil {
		fmt.Println("ERROR: unable to query db")
		panic(err)
	}

	for rows.Next() {
		err = rows.Scan(&id, &created, &name, &code)
		if err != nil {
			panic(err)
		}
		//tags := len(getVars(code))
		snipitem := Snipitem{
			ID:   id,
			Time: time.Now(),
			Name: name,
			Code: code,
		}

		snip = append(snip, snipitem)
	}

	rows.Close() //good habit to close
	return snip
}

func getArgumentList(text string) []string {
	var matches []string

	if len(text) > 0 {
		regexstring := regexp.MustCompile("{:[A-Za-z! ]+?:}")
		matches = regexstring.FindAllString(text, -1)
	}

	return matches
}

func getArguments(text string) []SnipArgs {
	var namelist []SnipArgs
	var varlist []string

	if len(text) > 0 {
		regexstring := regexp.MustCompile("{:[A-Za-z! ]+?:}")
		varlist = regexstring.FindAllString(text, -1)
	} else {
		return namelist
	}

	if len(varlist) <= 0 {
		return namelist
	}

	var varitem SnipArgs
	for _, varpos := range varlist {
		//var pos is start and end locations in array
		vars := strings.Split(varpos, ":")
		varname := strings.Split(vars[1], "!")
		fmt.Println(varname)
		if len(varname) == 1 {
			varitem = SnipArgs{
				Name:  varname[0],
				Value: "",
			}
		} else if len(varname) == 2 {
			varitem = SnipArgs{
				Name:  varname[0],
				Value: varname[1],
			}
		}
		namelist = append(namelist, varitem)
	}
	return namelist
}
