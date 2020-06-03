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
	//time.Sleep(4 * time.Second)
	//println("running from js: " + data)

	snips := dbfind("name", data)
	for _, itm := range snips {
		itm.Code = ""
	}
	str, _ := json.Marshal(snips)
	//fmt.Println("json: " + string(str))
	return string(str)
}

func snip_getvars(hash string) []SnipVars {
	var namelist []SnipVars

	snips := dbgetID(hash)
	varlist := getVars(snips.Code)

	for _, varpos := range varlist {
		//var pos is start and end locations in array
		vars := strings.Split(varpos, ":")
		varname := strings.Split(vars[1], "!")
		fmt.Println(varname)
		varitem := SnipVars{
			Name:  varname[0],
			Value: varname[1],
		}
		namelist = append(namelist, varitem)
	}
	return namelist
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
