package main

import (
	"encoding/json"
	"sync"
	"time"
)

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

}

func localSearch(wg *sync.WaitGroup, request searchRequest) {
	defer wg.Done() //update the wait counter on function exit
	defer close(request.channel)

	var foundSnips []snipItem
	var snips []snipItem

	logDebug("F:localSearch:start")

	if len(request.query) < 0 {
		return
	} else if len(request.query) == 0 {
		snips = dbGetAll()
	} else {
		snips = dbFind("name", request.query, 0) //search the name field in the snip table
	}

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

	var isLoggedIn bool = false
	if isLoggedIn != true {
		return
	}
	singlesnip := snipItem{
		Name:     "Z who on the web",
		Hash:     "64c42bc9-87e2-4771-85fe-07d05f9c0042",
		Code:     "curl google.com",
		Argument: nil,
		CmdType:  "bash",
	}
	foundSnips = append(foundSnips, singlesnip)

	singlesnip = snipItem{
		Name:     "A who in the house",
		Hash:     "64bcc429-87e2-4771-85fe-07d05f9c0042",
		Code:     "curl google.com",
		Argument: nil,
		CmdType:  "bash",
	}
	foundSnips = append(foundSnips, singlesnip)
	time.Sleep(150 * time.Millisecond)

	request.channel <- foundSnips

}
