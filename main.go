// yum install webkit2gtk3
// to build install sudo dnf install gtk3-devel webkit2gtk3-devel
// go get github.com/zserge/webview
// go get github.com/atotto/clipboard
// go get github.com/mattn/go-sqlite3
// sudo dnf install libxkbcommon-devel libXtst-devel libxkbcommon-x11-devel xorg-x11-xkb-utils-devel libpng-devel xsel xclip

package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	//"runtime"
	"strings"
	"sync"
	"termtyper/key"
	"time"
)

const webdebug bool = true
const loglevel int = 1
const defaultcmdtype string = "bash"
const appName string = "termtyper"
const regexMatch string = "{:[A-Za-z0-9!._ -]+?:}"

var codefromarg string = ""

type snipItem struct {
	Hash     string     `json:"hash"`
	Time     time.Time  `json:"-"` // - hides the output from json
	Name     string     `json:"value"`
	Code     string     `json:"code"`
	Argument []snipArgs `json:"argument"`
	CmdType  string     `json:"cmdtype"`
}

type snipArgs struct {
	Name  string `json:"name,omitempty"`
	Value string `json:"value,omitempty"`
}

type command struct {
	id   string
	data string
}

type searchRequest struct {
	hash    string
	query   string
	channel chan []snipItem
}

var action int
var datapath string

var queryQueue chan command

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

	database, _ = dbOpen(fldrName + "/termtyper.db")
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
	case "-m":
		typemanager(execpath)

	case "-n":
		codefromarg = readStdin()
		newfromcommand(execpath)

	case "--export":
		exportAll(progargs[1])

	case "--import":
		importAll(progargs[1])

	case "--help":
		showHelp()
	case "-h":
		showHelp()

	default:
		searchandpaste(execpath)
	}

	database.Close()
}

func showHelp() {
	fmt.Println("Usage:", appName, "[OPTION...]")
	fmt.Print(`
  -a                    open the autotype window
  -h, --help            show this help message
      --export FILE     export the local database to FILE
      --import FILE     import previously exported FILE into the local database

  -m                    open the managment UI
  -n                    read command from stdin (unix only)
  
`)

}

//return path to this running program
func getprogPath() string {
	logDebug("F:getprogPath:start")
	var dirAbsPath string
	ex, err := os.Executable()
	if err == nil {
		dirAbsPath = filepath.Dir(ex)
		return dirAbsPath
	}

	exReal, err := filepath.EvalSymlinks(ex)
	if err != nil {
		panic(err)
	}
	dirAbsPath = filepath.Dir(exReal)
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

//search text looking for arguments returns array of snipArgs
func getArguments(text string) []snipArgs {
	var namelist []snipArgs
	var varlist []string
	var varitem snipArgs

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
		strName := cleanString(varname[0], "[^A-Za-z0-9_. -]") //remove invalid chars from name
		if len(varname) == 1 {                                 // ! is optional so check if argument dosent have a default value
			varitem = snipArgs{
				Name:  strings.TrimSpace(strName),
				Value: "",
			}
		} else if len(varname) == 2 { //argument has a default value
			varitem = snipArgs{
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

//search code look ing for arguments, replace with values from snipArgs
func argumentReplace(vars []snipArgs, code string) string {
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
		emptyarg := snipArgs{Name: "", Value: ""}
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

func typeSnippet(messages chan bool, lineSeperator string, text []string) {
	if len(text) == 0 {
		logError("no text avaliable to type")
		messages <- true
		return
	}
	if lineSeperator == "" {
		_, lineSeperator = validCmdType(defaultcmdtype)
	}

	logDebug("F:typeSnippet:start")
	//runtime.LockOSThread()
	logDebug("F:typeSnippet:switching window")
	key.SwitchWindow()

	logDebug("F:typeSnippet:sending keys =", text)
	time.Sleep(2 * time.Second)
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
	//runtime.UnlockOSThread()

	messages <- true
}

func readStdin() string {
	var retstr string

	info, err := os.Stdin.Stat()
	if err != nil {
		panic(err)
	}

	if info.Mode()&os.ModeCharDevice != 0 {
		logError("No Pipe found")
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

func validCmdType(cmdtype string) (string, string) {
	var out string
	var sep string

	nowhite := strings.TrimSpace(cmdtype)
	lower := strings.ToLower(nowhite)
	out = lower

	switch lower {
	case "bash":
		sep = " \\"
	case "powershell":
		sep = ""
	case "dos":
		sep = " ^"

	default:
		out, sep = validCmdType(defaultcmdtype)
	}

	return out, sep
}

func exportAll(filename string) {
	foundSnips := dbGetAll()
	// {"hash":"000000-0000-0000-0000-000000000000",
	//  "name":"name for this command",
	//  "code":"actual command to type",
	//  "cmdtype":"bash"
	// }

	var out []map[string]interface{}

	for _, itm := range foundSnips {

		m := map[string]interface{}{
			"hash":    itm.Hash,
			"created": itm.Time,
			"name":    itm.Name,
			"code":    itm.Code,
			"cmdtype": itm.CmdType,
		}

		out = append(out, m)
	}
	strout, err := json.MarshalIndent(out, "", "  ")
	if err != nil {
		logError(err)
	}
	_ = ioutil.WriteFile(filename, []byte(strout), 0644)

}

func importAll(filename string) {

	file, _ := ioutil.ReadFile(filename)

	var items []snipItem
	_ = json.Unmarshal([]byte(file), &items)

	written := 0
	skipped := 0

	for i := 0; i < len(items); i++ {
		_, count := dbGetID(items[i].Hash)
		if count == 0 {
			//TODO: need to sanity check the data
			dbWrite(items[i].Hash, items[i].Time, items[i].Name, items[i].Code, items[i].CmdType)
			written++
		} else {
			skipped++
		}
	}

	fmt.Println(len(items), "total items to import,", written, "items imported successfully and", skipped, "items skipped")
}

//loop through each searchRequest item, wait for
// channels to timeout and close or recieve data, then merge data
// and send the data back with its hash using
// sendResultsToJS(hash, string(str)) this must be a json string though
func waitAndMerge(wg *sync.WaitGroup, requestList []searchRequest) {
	var totalSnips []snipItem

	for _, request := range requestList {
		items := <-request.channel
		totalSnips = append(totalSnips, items...)
	}

	wg.Wait() //wait for all search functions to finish

	hash := requestList[0].hash
	str, _ := json.Marshal(totalSnips)
	sendResultsToJS(hash, string(str))
	//fmt.Println("* Ready")
}

func localSearch(wg *sync.WaitGroup, request searchRequest) {
	defer wg.Done() //update the wait counter on function exit
	defer close(request.channel)

	var foundSnips []snipItem
	logDebug("F:localSearch:start")

	if len(request.query) <= 0 {
		return
	}

	snips := dbFind("name", request.query) //search the name field in the snip table
	for _, itm := range snips {
		itmarg := getArguments(itm.Code)
		itm.Argument = itmarg
		foundSnips = append(foundSnips, itm)
	}

	request.channel <- foundSnips
}

func remoteSearch(wg *sync.WaitGroup, request searchRequest) {
	defer wg.Done() //update the wait counter on function exit
	defer close(request.channel)

	var foundSnips []snipItem

	singlesnip := snipItem{
		Name:     "who on the web",
		Hash:     "64c42bc9-87e2-4771-85fe-07d05f9c0042",
		Code:     "curl google.com",
		Argument: nil,
		CmdType:  "bash",
	}
	foundSnips = append(foundSnips, singlesnip)
	//time.Sleep(150 * time.Millisecond)

	request.channel <- foundSnips

}
