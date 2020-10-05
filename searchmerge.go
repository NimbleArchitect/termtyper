package main

import (
	"database/sql"
	"encoding/json"

	"sync"
	"time"
)

//async search given a search id and query
// perform a search on seperate threads
func asyncSearch(hash string, searchfield string, query string) {
	var requestList []searchRequest

	wg := sync.WaitGroup{}
	if settings.Termtyper.EnableRemote == true && len(query) >= 2 {

		ch := make(chan []snipItem)
		newRequest := searchRequest{
			hash:        hash,
			searchfield: searchfield,
			query:       query,
			channel:     ch,
		}
		requestList = append(requestList, newRequest)
		wg.Add(1)
		go remoteSearch(&wg, newRequest)
	}
	//loop through each local db
	for _, database := range localDbList {
		//make a new channel that will hold the search response
		ch := make(chan []snipItem)
		newRequest := searchRequest{
			hash:        hash,
			searchfield: searchfield,
			query:       query,
			channel:     ch,
		}
		//add the channel to the requestList array
		requestList = append(requestList, newRequest)
		wg.Add(1) //increment wait group
		//and run the search
		go localSearch(&wg, database, newRequest)
	}

	//now we sit an wait for everyone to get back to us
	go waitAndMerge(&wg, requestList)

}

//loop through each searchRequest item, wait for
// channels to timeout and close or recieve data, then merge data
// filter out duplicates based on item hash id and send the data back with its request hash using
// sendResultsToJS(hash, string(str)) this must be a json string though
func waitAndMerge(wg *sync.WaitGroup, requestList []searchRequest) {
	totalSnips := make(map[string]snipItem)

	for _, request := range requestList {
		items := <-request.channel
		for _, singleItem := range items {
			_, ok := totalSnips[singleItem.Hash]
			if ok == false {
				totalSnips[singleItem.Hash] = singleItem
			}
		}
	}

	wg.Wait() //wait for all search functions to finish

	hash := requestList[0].hash
	str, _ := json.Marshal(totalSnips)

	sendResultsToJS(hash, string(str))

}

func localSearch(wg *sync.WaitGroup, database *sql.DB, request searchRequest) {
	defer wg.Done() //update the wait counter on function exit
	defer close(request.channel)

	var foundSnips []snipItem
	var snips []snipItem

	logDebug("F:localSearch:start")

	if len(request.query) <= 0 {
		return
	}

	snips = dbFind(database, request.searchfield, request.query, 0) //search the name field in the snip table

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

//returns a list of items sorted by popularity, defaults to top 20
func getPopular(hash string) {
	totalSnips := make(map[string]snipItem)
	items := dbGetPopular(localDbList[0], 20)

	for _, singleItem := range items {
		_, ok := totalSnips[singleItem.Hash]
		if ok == false {
			itmarg := getArguments(singleItem.Code)
			singleItem.Argument = itmarg
			totalSnips[singleItem.Hash] = singleItem
		}
	}
	str, _ := json.Marshal(totalSnips)

	sendResultsToJS(hash, string(str))
}

func getAllSnips(hash string) {
	totalSnips := make(map[string]snipItem)

	items := dbGetAll(localDbList[0])
	for _, singleItem := range items {
		_, ok := totalSnips[singleItem.Hash]
		if ok == false {
			itmarg := getArguments(singleItem.Code)
			singleItem.Argument = itmarg
			totalSnips[singleItem.Hash] = singleItem
		}
	}

	str, _ := json.Marshal(totalSnips)

	sendResultsToJS(hash, string(str))
}
