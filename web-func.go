package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/atotto/clipboard"
	"github.com/pborman/uuid"
	"strings"
	"sync"
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
func snipWrite(hash string, vars ...string) error {
	var code []string
	var data string
	var args []snipArgs
	logDebug("F:snip_write:start")

	if len(hash) <= 0 {
		return errors.New("no hash id specified")
	}
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
	//set up channel to wait on, this fixes a crash where the window
	// was closing before the fucntion had finished
	messages := make(chan bool)
	go typeSnippet(messages, sep, code)
	//wait for completion signal
	<-messages
	snipClose()

	return nil
}

//save to db
func snipSave(title string, code string, commandtype string) {
	logDebug("F:snip_save:start")

	hash := uuid.New()

	cmdtype, _ := validCmdType(commandtype)

	dbWrite(hash, time.Now(), title, code, cmdtype)

}

//returns json object representing string passed from stdin
func snipCodeFromArg() string {
	var thissnip snipItem
	logDebug("F:snip_codeFromArg:start")

	thissnip.Code = codefromarg
	str, _ := json.Marshal(thissnip)
	return string(str)
}

func snipSearchRemote() string {
	remoteenabled := false

	if remoteenabled == false {
		return ""
	} else {

		return `function snipSearchRemote( request, response ) {
	if (request.term.length >= 2) {
		$.getJSON("http://localhost:8080/sch?t=123456",request,response)
	}
}`

	}
}

//read fromclipboard
func snipGetClipboard() string {
	out, err := clipboard.ReadAll()
	if err == nil {
		return out
	}
	return ""
}

type searchRequest struct {
	hash    string
	query   string
	channel chan []snipItem
}

//async search given a search id and query
// perform a search on seperate threads
func snipAsyncSearch(hash string, query string) error {

	var requestList []searchRequest

	remoteActive := true

	wg := sync.WaitGroup{}
	if remoteActive == true {
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

func waitAndMerge(wg *sync.WaitGroup, requestList []searchRequest) {
	var totalSnips []snipItem
	//TODO: loop through each searchRequest item, wait for
	// channels to timeout and close or recieve data, then merge data
	// and send the data back with its hash using
	// sendResultsToJS(hash, string(str)) this must be a json string though

	fmt.Println("*** Waiting....")
	for _, request := range requestList {
		items := <-request.channel
		totalSnips = append(totalSnips, items...)
		fmt.Println(items[0].Name)
	}

	wg.Wait() //wait for all search functions to finish

	hash := requestList[0].hash
	str, _ := json.Marshal(totalSnips)
	sendResultsToJS(hash, string(str))
	fmt.Println("* Ready")
}
