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
	Time     time.Time  `json:"time"`
	Name     string     `json:"name"`
	Code     string     `json:"code"`
	Argument []SnipArgs `json:"argument"`
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

func getArgumentPos(text string) ([][]int, bool) {
	var ok bool
	var matches [][]int

	if len(text) > 0 {
		regexstring := regexp.MustCompile("{:[A-Za-z! ]+?:}")
		matches = regexstring.FindAllStringIndex(text, -1)
		ok = true
	} else {
		ok = false
	}
	if len(matches) <= 0 {
		ok = false
	}

	return matches, ok
}

func getArgumentList(text string) ([]string, bool) {
	var ok bool
	var matches []string

	if len(text) > 0 {
		regexstring := regexp.MustCompile("{:[A-Za-z! ]+?:}")
		matches = regexstring.FindAllString(text, -1)
		ok = true
	} else {
		ok = false
	}
	if len(matches) <= 0 {
		ok = false
	}

	return matches, ok
}

func getArguments(text string) []SnipArgs {
	var namelist []SnipArgs
	var varlist []string

	varlist, ok := getArgumentList(text)
	if ok == false {
		return namelist
	}

	var varitem SnipArgs
	for _, varpos := range varlist {
		//var pos is start and end locations in array
		vars := strings.Split(varpos, ":")
		varname := strings.Split(vars[1], "!")
		if len(varname) == 1 {
			varitem = SnipArgs{
				Name:  strings.TrimSpace(varname[0]),
				Value: "",
			}
		} else if len(varname) == 2 {
			varitem = SnipArgs{
				Name:  strings.TrimSpace(varname[0]),
				Value: strings.TrimSpace(varname[1]),
			}
		}
		namelist = append(namelist, varitem)
	}
	return namelist
}

//TODO: check vars is valid, check snips.code has something
func argumentReplace(vars []SnipArgs, code string) string {
	var newcode string
	var val string

	if len(code) <= 0 {
		return ""
	}
	itmarg := getArguments(code)
	argPos, _ := getArgumentPos(code)
	//spin through all arguments and replace variables as needed
	itmlen := len(itmarg) - 1

	newcode = code
	for i := itmlen; i >= 0; i-- {
		itm := itmarg[i]
		//make sure the incomming argument name matches the
		if itm.Name != vars[i].Name {
			return ""
		}

		if len(vars[i].Value) > 0 {
			val = vars[i].Value //incomming value is valid so use that
		} else if len(itm.Value) > 0 {
			val = itm.Value //incoming value not valid but we have a default value so use it
		} else {
			val = "{" + itm.Name + "}" //nothing is valid so we default to the name in braces
		}

		itmpos := argPos[i] //start and end pos of txt to replace
		s := itmpos[0]
		e := itmpos[1]

		newcode = newcode[:s] + val + newcode[e:]
	}

	return newcode
}
