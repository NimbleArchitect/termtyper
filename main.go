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
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"termtyper/key"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

const appName string = "termtyper"

// Config struct for webapp config
type config struct {
	termtyper struct {
		Debug        bool   `yaml:"debug"`
		LogLevel     int    `yaml:"loglevel"`
		EnableRemote bool   `yaml:"remote"`
		CmdType      string `yaml:"cmdtype"`
		maxRows      int    `yaml:"maxrows"`
		Path         string `yaml:"path"`
		TypeDelay    int    `yaml:"typedelay"` //time to wait between key presses, in milliseconds
		LineDelay    int64  `yaml:"linedelay"` //time to wait after pressing enter, in milliseconds
		// Timeout      struct {
		// 	Server time.Duration `yaml:"server"`
		// 	Write  time.Duration `yaml:"write"`
		// } `yaml:"timeout"`
	} `yaml:"termtyper"`
}

var settings config = config{}
var codefromarg string = ""

type snipItem struct {
	Hash     string     `json:"hash"`
	Time     time.Time  `json:"-"` // - hides the output from json
	Name     string     `json:"value"`
	Code     string     `json:"code"`
	Argument []snipArgs `json:"argument"`
	CmdType  string     `json:"cmdtype"`
	Summary  string     `json:"summary"`
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

func loadSettings() config {
	//set defaults
	settings.termtyper.Debug = false
	settings.termtyper.LogLevel = 1
	settings.termtyper.EnableRemote = false
	settings.termtyper.CmdType = "bash"
	settings.termtyper.maxRows = 20
	settings.termtyper.Path = ""
	settings.termtyper.TypeDelay = 20
	settings.termtyper.LineDelay = 200

}

func main() {
	var argument string
	const html = `
	<html><head></head><body>
	Move along nothing to see here
	</body></html>`

	loadSettings()

	if settings.termtyper.Path == "" {
		var pathSep = "/"
		appFolder := "." + appName
		datapath, err := os.UserHomeDir()
		if err != nil {
			panic("Unable to get users profile folder")
		}
		settings.termtyper.Path = datapath + pathSep + appFolder
	}
	fldrName := settings.termtyper.Path
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

func typeSnippet(lineSeperator string, text []string) {
	logDebug("F:typeSnippet:start")
	delay := time.Duration(settings.termtyper.LineDelay)
	if len(text) == 0 {
		logError("no text avaliable to type")
		//messages <- true
		return
	}
	if lineSeperator == "" {
		_, lineSeperator = validCmdType(settings.termtyper.CmdType)
	}

	logDebug("F:typeSnippet:sending keys =", text)
	time.Sleep(2 * time.Second)
	//send keys to type to stdin of python script :(
	count := len(text)
	for i := 0; i < count; i++ {
		singleline := text[i]
		if i < (count - 1) { //more than one line and we are not on the last
			key.SendLine(singleline+lineSeperator+"\n", settings.termtyper.TypeDelay) //sent line of text
			time.Sleep(delay * time.Millisecond)                                      //sleep after pressing enter
		} else {
			key.SendLine(singleline, settings.termtyper.TypeDelay) //write the last or only line
		}
	}

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
		out, sep = validCmdType(settings.termtyper.CmdType)
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
			"summary": itm.Summary,
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
			dbWrite(items[i].Hash, items[i].Time, items[i].Name, items[i].Code, items[i].CmdType, items[i].Summary)
			written++
		} else {
			skipped++
		}
	}

	fmt.Println(len(items), "total items to import,", written, "items imported successfully and", skipped, "items skipped")
}
