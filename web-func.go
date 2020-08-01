package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"strings"
	"sync"
	"time"

	"github.com/atotto/clipboard"
	"github.com/pborman/uuid"
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
func snipWrite(hash string, vars ...string) error {
	if len(hash) <= 0 {
		return errors.New("no hash id specified")
	}

	go asyncWrite(hash, vars)

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

//async search given a search id and query
// perform a search on seperate threads
func snipAsyncSearch(hash string, query string) error {
	var requestList []searchRequest

	wg := sync.WaitGroup{}
	if settings.Termtyper.EnableRemote == true && len(query) >= 2 {

		ch := make(chan []snipItem)
		newRequest := searchRequest{
			hash:    hash,
			query:   query,
			channel: ch,
		}
		requestList = append(requestList, newRequest)
		wg.Add(1)
		go remoteSearch(&wg, newRequest)
	}
	ch := make(chan []snipItem)
	newRequest := searchRequest{
		hash:    hash,
		query:   query,
		channel: ch,
	}
	requestList = append(requestList, newRequest)
	wg.Add(1)
	go localSearch(&wg, newRequest)

	go waitAndMerge(&wg, requestList)
	return nil
}

func asyncWrite(hash string, vars []string) {
	var code []string
	var data string
	var args []snipArgs
	logDebug("F:snip_write:start")

	logDebug("F:snip_write:hash =", hash)
	snips, _ := dbGetID(hash)
	logDebug("F:snip_write:snips =", snips)
	logDebug("F:snip_write:len(vars) =", len(vars))

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
	logDebug("F:asyncWrite:switching window")
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
