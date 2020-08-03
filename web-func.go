package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"github.com/atotto/clipboard"
	"github.com/pborman/uuid"
	"strings"
	"time"
)

//copy data into clipboard
func snipCopy(data string) error {
	logDebug("F:snip_copy:start")

	clipboard.WriteAll(data)
	return nil
}

//auto type function
// accepts given hash matching snip record and
// a json string representing argument name and value
func snipTyper(hash string, vars ...string) error {
	if len(hash) <= 0 {
		return errors.New("no hash id specified")
	}

	go asyncTyper(hash, vars)

	return nil
}

//save to db
func snipSave(title string, code string, commandtype string, summary string) {
	logDebug("F:snip_save:start")

	hash := uuid.New()

	cmdtype, _ := validCmdType(commandtype)

	dbWrite(hash, time.Now(), title, code, cmdtype, summary)

}

//returns json object representing string passed from stdin
func snipCodeFromArg() string {
	var thissnip snipItem
	logDebug("F:snip_codeFromArg:start")

	thissnip.Code = codefromarg
	str, _ := json.Marshal(thissnip)
	return string(str)
}

//read fromclipboard
func snipGetClipboard() string {
	out, err := clipboard.ReadAll()
	if err == nil {
		return out
	}
	return ""
}

func snipAsyncRequest(hash string, jsonQuery string) error {
	type asyncRequest struct {
		Operation string `json:"operation"`
		Value     string `json:"value"`
	}
	var request asyncRequest
	json.Unmarshal([]byte(jsonQuery), &request)

	switch request.Operation {
	case "search":
		go asyncSearch(hash, "name", request.Value)
	case "searchcode":
		go asyncSearch(hash, "code", request.Value)

	case "get":
		switch request.Value {
		case "popular":
			go getPopular(hash)
		case "all":
			go getAllSnips(hash)
		}
	}
	return nil
}

func asyncTyper(hash string, vars []string) {
	var code []string
	var data string
	var args []snipArgs
	logDebug("F:asyncTyper:start")

	logDebug("F:asyncTyper:hash =", hash)
	snips, _ := dbGetID(hash)
	logDebug("F:asyncTyper:snips =", snips)
	logDebug("F:asyncTyper:len(vars) =", len(vars))

	if len(vars) > 0 {
		json.Unmarshal([]byte(vars[0]), &args)
	}
	data = argumentReplace(args, snips.Code)
	scanner := bufio.NewScanner(strings.NewReader(data))

	for scanner.Scan() {
		singleline := scanner.Text()
		code = append(code, singleline)
	}

	_, sep := validCmdType(snips.CmdType) //get multiline seperator
	logDebug("F:asyncTyper:switching window")
	minimizeWindow()

	time.Sleep(3 * time.Second)
	typeSnippet(sep, code)
	dbUpdatePopular(hash) //update usage counter

	snipClose()
}

func snipClose() error {
	logDebug("F:snip_close:start")

	w.Dispatch(func() {
		w.Terminate()
	})
	return nil
}
