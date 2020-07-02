package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"github.com/atotto/clipboard"
	"strings"
	"time"
)

func snip_copy(data string) error {
	logDebug("F:snip_copy:start")

	clipboard.WriteAll(data)
	return nil
}

func snip_close() error {
	logDebug("F:snip_close:start")

	w.Terminate()
	return nil
}

func snip_search(data string) string {
	var foundSnips []Snipitem
	logDebug("F:snip_search:start")

	if len(data) <= 0 {
		return ""
	}

	snips := dbfind("name", data) //search the name field in the snip table
	for _, itm := range snips {
		itmarg := getArguments(itm.Code)
		itm.Argument = itmarg
		foundSnips = append(foundSnips, itm)
		//fmt.Println(foundSnips)
	}
	str, _ := json.Marshal(foundSnips)
	return string(str)
}

func snip_write(hash string, vars ...string) error {
	var code []string
	var data string
	var args []SnipArgs
	logDebug("F:snip_write:start")

	if len(hash) <= 0 {
		return errors.New("no hash id specified")
	}
	logDebug("F:snip_write:hash =", hash)
	snips := dbgetID(hash)
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

	go typeSnippet(code)
	return nil
}

func snip_save(title string, code string) {
	logDebug("F:snip_save:start")

	tx, _ := database.Begin()
	stmt, _ := tx.Prepare("insert into snips (created,name,code) values (?,?,?)")
	_, err := stmt.Exec(time.Now(), title, code)
	if err != nil {
		logError("error saving")
	}
	tx.Commit()
}

func snip_codeFromArg() string {
	var thissnip Snipitem
	logDebug("F:snip_codeFromArg:start")

	thissnip.Code = codefromarg
	str, _ := json.Marshal(thissnip)
	return string(str)
}
