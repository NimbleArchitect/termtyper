package main

import (
	"bufio"
	"encoding/json"
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
	if len(data) <= 0 {
		return ""
	}

	snips := dbfind("name", data)
	for _, itm := range snips {
		itmarg := getArguments(itm.Code)
		itm.Argument = itmarg
	}
	str, _ := json.Marshal(snips)
	return string(str)
}

func snip_write(hash string, vars ...string) error {
	var code []string
	//fmt.Println("** " + hash)string, values []
	snips := dbgetID(hash)
	data := snips.Code

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
