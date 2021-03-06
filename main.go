// yum install webkit2gtk3
// to build install sudo dnf install gtk3-devel webkit2gtk3-devel
// go get github.com/zserge/webview
// go get github.com/atotto/clipboard
// go get github.com/mattn/go-sqlite3
// sudo dnf install libxkbcommon-devel libXtst-devel libxkbcommon-x11-devel xorg-x11-xkb-utils-devel libpng-devel xsel xclip

package main

import (
	"bufio"
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"termtyper/key"
	"time"
)

const appName string = "termtyper"

// Config struct for webapp config
type config struct {
	Termtyper struct {
		Debug        bool   `yaml:"debug"`
		WebLocal     bool   `yaml:"weblocal"`
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
		Database []struct {
			Name string `yaml:"name"`
			Path string `yaml:"path"`
			Type string `yaml:"type"`
		}
	} `yaml:"termtyper"`
}

var settings config
var codefromarg string = ""

type snipItem struct {
	Hash     string     `json:"-"` // we dont need this in json output as the map key is the hash id also
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
	hash        string
	searchfield string
	query       string
	channel     chan []snipItem
}

var action int
var datapath string

var queryQueue chan command

func loadSettings(filename string) {
	//set defaults
	settings.Termtyper.Debug = false
	settings.Termtyper.WebLocal = false
	settings.Termtyper.LogLevel = 1
	settings.Termtyper.EnableRemote = false
	settings.Termtyper.CmdType = "bash"
	settings.Termtyper.maxRows = 20
	settings.Termtyper.Path = "~/.termtyper/"
	settings.Termtyper.TypeDelay = 20
	settings.Termtyper.LineDelay = 200

	settings.Termtyper.Path = convertPath(settings.Termtyper.Path)
	data, err := ioutil.ReadFile(settings.Termtyper.Path + filename)

	if err != nil {
		logWarn("unable to open config file", settings.Termtyper.Path+filename)
	}
	//var conf config
	err = yaml.Unmarshal([]byte(data), &settings)
	if err != nil {
		//TODO: panic here also
		logWarn("unable to process config file")
	}
	//fmt.Println(conf)
	if len(settings.Termtyper.Database) <= 0 {
		//this is a bit hacky, could do with fixing
		settings.Termtyper.Database = append(
			settings.Termtyper.Database, struct {
				Name string "yaml:\"name\""
				Path string "yaml:\"path\""
				Type string "yaml:\"type\""
			}{"default", "~/.termtyper/termtyper.db", "local"},
		)
	}
	//loop through each database correcting the path
	for i, database := range settings.Termtyper.Database {
		settings.Termtyper.Database[i].Path = convertPath(database.Path)
	}

}

func main() {
	var argument string
	const html = `
	<html><head></head><body>
	Move along nothing to see here
	</body></html>`

	fldrName := settings.Termtyper.Path
	loadSettings("config.yaml")
	fldrName = settings.Termtyper.Path

	if _, err := os.Stat(fldrName); err != nil {
		err = os.Mkdir(fldrName, 0770)
		if err != nil {
			panic("unable to create folder " + fldrName)
		}
	}

	count := openDatabases()
	if count <= 0 {
		panic("unable to open any listed databases")
	}
	// database, _ = dbOpen(settings.Termtyper.Database[0].Path)
	// // if ok == true {
	// // 	//defer database.Close()
	// // }
	// logDebug("F:main:db ping")
	// if database.Ping() != nil {
	// 	logError("F:main:unable to ping db")
	// }
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

	closeDatabases()
	//database.Close()
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
	delay := time.Duration(settings.Termtyper.LineDelay)
	if len(text) == 0 {
		logError("no text avaliable to type")
		//messages <- true
		return
	}
	if lineSeperator == "" {
		_, lineSeperator = validCmdType(settings.Termtyper.CmdType)
	}

	logDebug("F:typeSnippet:sending keys =", text)
	//send keys to type to stdin of python script :(
	count := len(text)
	for i := 0; i < count; i++ {
		singleline := text[i]
		if i < (count - 1) { //more than one line and we are not on the last
			key.SendLine(singleline+lineSeperator+"\n", settings.Termtyper.TypeDelay) //sent line of text
			time.Sleep(delay * time.Millisecond)                                      //sleep after pressing enter
		} else {
			key.SendLine(singleline, settings.Termtyper.TypeDelay) //write the last or only line
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
		out, sep = validCmdType(settings.Termtyper.CmdType)
	}

	return out, sep
}

func exportAll(filename string) {
	//FIXME: currently this will only export the first database in the list
	foundSnips := dbGetAll(localDbList[0])
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
	//FIXME: currently this will only import into the first database in the list

	file, _ := ioutil.ReadFile(filename)

	var items []snipItem
	_ = json.Unmarshal([]byte(file), &items)

	written := 0
	skipped := 0

	for i := 0; i < len(items); i++ {
		_, count := dbGetID(localDbList[0], items[i].Hash)
		if count == 0 {
			//TODO: need to sanity check the data
			dbWrite(localDbList[0], items[i].Hash, items[i].Time, items[i].Name, items[i].Code, items[i].CmdType, items[i].Summary)
			written++
		} else {
			skipped++
		}
	}

	fmt.Println(len(items), "total items to import,", written, "items imported successfully and", skipped, "items skipped")
}

func openDatabases() int {
	var database *sql.DB
	var count int = 0

	for _, db := range settings.Termtyper.Database {
		database, _ = dbOpen(db.Path)
		// if ok == true {
		// 	//defer database.Close()
		// }
		logDebug("F:main:db ping")
		//check db is valid and opens
		if database.Ping() != nil {
			logError("F:main:unable to ping db")
		} else {
			//append to this array localDbList
			localDbList = append(localDbList, database)
			count++
		}
	}
	return count
}

func closeDatabases() {
	for _, database := range localDbList {
		database.Close()
	}
}

func convertPath(path string) string {
	newpath := path
	if strings.HasPrefix(path, "~") {
		datapath, err := os.UserHomeDir()
		if err != nil {
			panic("Unable to get users profile folder")
		}
		newpath = strings.Replace(path, "~", datapath, 1)
	}
	return newpath
}
