// yum install webkit2gtk3
// to build install sudo dnf install gtk3-devel webkit2gtk3-devel
// go get github.com/zserge/webview
// go get github.com/atotto/clipboard
// go get github.com/mattn/go-sqlite3
// sudo dnf install libxkbcommon-devel libXtst-devel libxkbcommon-x11-devel xorg-x11-xkb-utils-devel libpng-devel xsel xclip

package main

import (
	"bufio"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"io"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"termtyper/key"
	"time"
)

const webdebug bool = true
const loglevel int = 5
const appName string = "termtyper"
const regexMatch string = "{:[A-Za-z_-]+?.*:}"

var codefromarg string = ""

type Snipitem struct {
	ID       int        `json:"hash"`
	Time     time.Time  `json:"-"` // - hides the output from json
	Name     string     `json:"name"`
	Code     string     `json:"code"`
	Argument []SnipArgs `json:"argument"`
}

type SnipArgs struct {
	Name  string `json:"name,omitempty"`
	Value string `json:"value,omitempty"`
}

var action int
var datapath string

func main() {
	var argument string
	const html = `
	<html><head></head><body>
	Move along nothing to see here
	</body></html>`

	var pathSep = "/"
	appFolder := "." + appName

	datapath, err := os.UserHomeDir()
	if err != nil {
		panic("Unable to get users profile folder")
	}
	fldrName := datapath + pathSep + appFolder
	if _, err := os.Stat(fldrName); err != nil {
		err = os.Mkdir(fldrName, 0770)
		if err != nil {
			panic("unable to create folder " + fldrName)
		}
	}

	database, _ = opendb(fldrName + "/snippets.db")
	// if ok == true {
	// 	//defer database.Close()
	// }
	logDebug("F:main:db ping")
	if database.Ping() != nil {
		logError("F:main:unable to ping db")
	}
	logDebug("F:main:call getprogPath")
	execpath := getprogPath()
	//TODO: set up to support arguments to show the search window, I can then show a managment window by default

	logDebug("F:main:get args =", os.Args)
	progargs := os.Args[1:]
	if len(progargs) >= 1 {
		argument = progargs[0]
	} else {
		argument = ""
	}
	switch argument {
	case "-n":
		fmt.Println(progargs)
		codefromarg = readStdin()
		newfromcommand(execpath)

	case "2":
		fmt.Println("two")
	default:
		searchandpaste(execpath)
	}

	database.Close()
}

//return path to this running program
func getprogPath() string {
	logDebug("F:getprogPath:start")
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
	logDebug("F:getprogPath:return =", dirAbsPath)
	return dirAbsPath
}

// returns a list of start and end positions of each argument, found in text
func getArgumentPos(text string) ([][]int, bool) {
	var ok bool
	var matches [][]int
	logDebug("F:getArgumentPos:start")
	if len(text) > 0 {
		regexstring := regexp.MustCompile(regexMatch)
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
		regexstring := regexp.MustCompile(regexMatch)
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
		logDebug("F:getArguments:varname =", varname)
		logDebug("F:getArguments:len(varname) =", len(varname))
		strName := cleanString(varname[0], "[^A-Za-z_.-]") //remove invalid chars from name
		if len(varname) == 1 {                             // ! is optional so check if argument dosent have a default value
			varitem = SnipArgs{
				Name:  strings.TrimSpace(strName),
				Value: "",
			}
		} else if len(varname) == 2 { //argument has a default value
			varitem = SnipArgs{
				Name:  strings.TrimSpace(strName),
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

func cleanString(data string, regex string) string {
	logDebug("F:cleanString:start")
	logDebug("F:cleanString:data =", data)

	reg, err := regexp.Compile(regex)
	if err != nil {
		log.Fatal(err)
	}
	newstr := reg.ReplaceAllString(data, "")
	logDebug("F:cleanString:return -", newstr)
	return newstr
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

func typeSnippet(text []string) {
	lineSeperator := " \\"

	key.SwitchWindow()

	//send keys to type to stdin of python script :(
	count := len(text)
	for i := 0; i < count; i++ {
		singleline := text[i]
		if i < (count - 1) { //more than one line and we are not on the last
			key.SendLine(singleline + lineSeperator + "\n") //sent line of text
		} else {
			key.SendLine(singleline) //write the last or only line
		}
	}

	w.Terminate()
}

func readStdin() string {
	var retstr string

	info, err := os.Stdin.Stat()
	if err != nil {
		panic(err)
	}

	if info.Mode()&os.ModeCharDevice != 0 {
		fmt.Println("No Pipe found")
		//return
	}

	reader := bufio.NewReader(os.Stdin)
	var output []string

	for {
		input, _, err := reader.ReadRune()
		if err != nil && err == io.EOF {
			output = append(output, string(input))
			break
		}
		output = append(output, string(input))
	}

	for j := 0; j < len(output); j++ {
		retstr += output[j]
	}
	return strings.TrimSpace(retstr)
}

func exportAll() {

}

func importAll() {

}
