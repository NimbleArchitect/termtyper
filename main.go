// yum install webkit2gtk3
// to build install sudo dnf install gtk3-devel webkit2gtk3-devel
// go get github.com/zserge/webview
// go get github.com/atotto/clipboard
// go get github.com/mattn/go-sqlite3
// sudo dnf install libxkbcommon-devel libXtst-devel libxkbcommon-x11-devel xorg-x11-xkb-utils-devel libpng-devel xsel xclip

package main

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"github.com/zserge/webview"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

const debug bool = true
const loglevel int = 1
const appName string = "termtyper"

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
	fldrName := datapath + "/." + appName
	if _, err := os.Stat(fldrName); err != nil {
		err = os.Mkdir(fldrName, 0770)
		if err != nil {
			panic("unable to create folder ~/." + appName)
		}
	}

	database, _ = opendb(fldrName + "/snippets.db")
	// if ok == true {
	// 	//defer database.Close()
	// }
	if database.Ping() != nil {
		fmt.Println("99")
	}
	execpath := getprogPath()
	//TODO: set up to support arguments to show the search window, I can then show a managment window by default
	searchandpaste(execpath)
	database.Close()
}

func logError(msg ...interface{}) {
	if loglevel >= 1 {
		log.Print("[ERROR] ", msg)

	}
}

func logWarn(msg ...interface{}) {
	if loglevel >= 2 {
		log.Print("[WARN]", msg)
	}
}

func logInfo(msg ...interface{}) {
	if loglevel >= 3 {
		log.Print("[INFO] ", msg)
	}
}

func logDebug(msg ...interface{}) {
	if loglevel >= 4 {
		log.Print("[DEBUG] ", msg)
	}
}

//return path to this running program
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
	w.SetTitle(appName)
	w.SetSize(800, 600, webview.HintNone)
	//w.Navigate("data:text/html," + html)
	w.Navigate("file://" + datapath + "/searchpage.html")
	w.Bind("snipSearch", snip_search)
	w.Bind("toclipboard", snip_copy)
	w.Bind("snipWrite", snip_write)
	w.Bind("snipClose", snip_close)
	w.Bind("snipSave", snip_save)
	w.Run()
}

func opendb(dbpath string) (*sql.DB, bool) {
	logInfo("* open: " + dbpath)
	db, err := sql.Open("sqlite3", dbpath)

	if err != nil {
		logError("unablet o open database")
	}
	ok := db.Ping()
	if ok != nil {
		panic(err)
	}
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS snips (id INTEGER PRIMARY KEY, created INTEGER, name TEXT, code TEXT)")

	return db, true
}

// returns a Snipitem that represents the hashid from the database table
func dbgetID(hash string) Snipitem {
	var snip Snipitem
	var id int
	var name string
	var code string
	var created string

	qry := string("SELECT * FROM snips WHERE ID = " + hash)
	rows, err := database.Query(qry)
	if err != nil {
		logError("ERROR: unable to query db")
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

// search for a record name that has a wildcard match to field and return a Snipitem that represents the match
func dbfind(field string, searchfor string) []Snipitem {
	//TODO: search search for matching tags
	var snip []Snipitem
	var id int
	var name string
	var code string
	var created string

	// query
	qry := string("SELECT * FROM snips WHERE " + field + " LIKE '%" + searchfor + "%'")
	rows, err := database.Query(qry)

	if err != nil {
		logError("ERROR: unable to query db")
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

// returns a list of start and end positions of each argument, found in text
func getArgumentPos(text string) ([][]int, bool) {
	var ok bool
	var matches [][]int
	logDebug("F:getArgumentPos:start")
	if len(text) > 0 {
		regexstring := regexp.MustCompile("{:[A-Za-z0-9! ]+?:}")
		matches = regexstring.FindAllStringIndex(text, -1)
		ok = true
	} else {
		ok = false
	}
	if len(matches) <= 0 {
		ok = false
	}
	logDebug("F:getArgumentPos:return =", matches, ",", ok)
	return matches, ok
}

// return string array of arguments found in text
func getArgumentList(text string) ([]string, bool) {
	var ok bool
	var matches []string
	logDebug("F:getArgumentList:start =", text)

	if len(text) > 0 {
		regexstring := regexp.MustCompile("{:[A-Za-z0-9! ]+?:}")
		matches = regexstring.FindAllString(text, -1)
		ok = true
	} else {
		ok = false
	}
	if len(matches) <= 0 {
		ok = false
	}
	logDebug("F:getArgumentList:return =", matches, ",", ok)
	return matches, ok
}

//search text looking for arguments returns array of SnipArgs
func getArguments(text string) []SnipArgs {
	var namelist []SnipArgs
	var varlist []string
	var varitem SnipArgs
	logDebug("F:getArguments:start")
	varlist, ok := getArgumentList(text)
	if ok == false { //no arguments found in text
		logDebug("F:getArguments:return =", namelist)
		return namelist
	}

	for _, varpos := range varlist {
		//var pos is start and end locations in array
		vars := strings.Split(varpos, ":")     //arguments are enclosed in : so we remove those first
		varname := strings.Split(vars[1], "!") // default values for arguments can be found after !
		if len(varname) == 1 {                 // ! is optional so check if argument dosent have a default value
			varitem = SnipArgs{
				Name:  strings.TrimSpace(varname[0]),
				Value: "",
			}
		} else if len(varname) == 2 { //argument has a default value
			varitem = SnipArgs{
				Name:  strings.TrimSpace(varname[0]),
				Value: strings.TrimSpace(varname[1]),
			}
		} else {
			// multipule defaults values have been suppilied so write warning
			logWarn("multipule default values detected.")
		}

		namelist = append(namelist, varitem)
	}
	logDebug("F:getArguments:return =", namelist)
	return namelist
}

//search code looing for arguments, replace with values from SnipArgs
func argumentReplace(vars []SnipArgs, code string) string {
	var newcode string
	var val string
	logDebug("F:argumentReplace:start")
	if len(code) <= 0 {
		return ""
	}
	itmarg := getArguments(code)      // get array of arguments from code
	argPos, _ := getArgumentPos(code) // get array or argument start/end positions
	//spin through all arguments and replace variables as needed
	itmlen := len(itmarg)
	varlen := len(vars)
	if varlen < 0 {
		varlen = 0
	}

	if varlen != itmlen { // itmlen is not the same length as varlen
		emptyarg := SnipArgs{Name: "", Value: ""}
		for c := varlen; c <= itmlen; c++ {
			vars = append(vars, emptyarg) //so add enough empty values to our vars array this helps the for loop below
		}
	}
	newcode = code
	logDebug("V:itmlen =", itmlen)
	logDebug("V:varlen =", len(vars))

	startlen := itmlen - 1
	for i := (startlen); i >= 0; i-- {
		logDebug("V:i =", i)
		itm := itmarg[i]
		//logDebug("F:argumentReplace:vars[i] =", vars[i])
		logDebug("F:argumentReplace:itm =", itm)

		//make sure the incomming argument name matches the variable in the code
		if itm.Name != vars[i].Name {
			if vars[i].Name != "" { //vars name could of been added above so can be empty
				logDebug("F:argumentReplace:return = \"\"")
				return "" // refuse and exit function if name wasn't empty
			}
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
	logDebug("F:argumentReplace:return =", newcode)
	return newcode
}
