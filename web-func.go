package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/atotto/clipboard"
	"strings"
	"time"
)

func snip_copy(data string) error {
	clipboard.WriteAll(data)
	return nil
}

func snip_close() error {
	go w.Terminate()
	return nil
}

func snip_search(data string) string {
	var foundSnips []Snipitem

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

	if len(hash) <= 0 {
		return errors.New("no hash id specified")
	}

	snips := dbgetID(hash)
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
	tx, _ := database.Begin()
	stmt, _ := tx.Prepare("insert into snips (created,name,code) values (?,?,?)")
	_, err := stmt.Exec(time.Now(), title, code)
	if err != nil {
		fmt.Print("error saving")
	}
	tx.Commit()
}
