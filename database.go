package main

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"time"
)

var database *sql.DB

func dbOpen(dbpath string) (*sql.DB, bool) {
	logInfo("* open: " + dbpath)
	db, err := sql.Open("sqlite3", dbpath)

	if err != nil {
		logError("unable to open database")
	}
	ok := db.Ping()
	if ok != nil {
		panic(err)
	}
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS "snips" (
		"hash"	TEXT UNIQUE,
		"created"	INTEGER,
		"name"	TEXT,
		"code"	TEXT,
		"cmdtype"	TEXT,
		PRIMARY KEY("hash")
	);`,
	)

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS "popular" (
		"hash"	TEXT UNIQUE,
		"timesused"	INTEGER,
		"lastused"	TEXT,
		PRIMARY KEY("hash")
	);`,
	)

	return db, true
}

// returns a Snipitem that represents the hashid from the database table
func dbGetID(hash string) (snipItem, int) {
	var snip snipItem

	var name string
	var code string
	var created string
	var cmdtype string

	qry := string("SELECT * FROM snips WHERE hash = '" + hash + "'")
	rows, err := database.Query(qry)
	if err != nil {
		logError("ERROR: unable to query db")
		panic(err)
	}
	count := 0
	for rows.Next() {
		err = rows.Scan(&hash, &created, &name, &code, &cmdtype)
		if err != nil {
			panic(err)
		}
		//tags := len(getVars(code))
		snip = snipItem{
			Hash:    hash,
			Time:    time.Now(),
			Name:    name,
			Code:    code,
			CmdType: cmdtype,
		}
		count++
	}
	rows.Close() //good habit to close

	return snip, count
}

// search for a record name that has a wildcard match to field and return a Snipitem that represents the match
func dbFind(field string, searchfor string) []snipItem {
	//TODO: search for matching tags
	var snip []snipItem
	var hash string
	var name string
	var code string
	var created string
	var cmdtype string
	var qry string = ""

	// query
	if field == "name" {
		qry = string("SELECT * FROM snips WHERE name LIKE ?")
	}
	if field == "code" {
		qry = string("SELECT * FROM snips WHERE code LIKE ?")
	}
	rows, err := database.Query(qry, "%"+searchfor+"%")

	if err != nil {
		logError("ERROR: unable to query db")
		panic(err)
	}

	for rows.Next() {
		err = rows.Scan(&hash, &created, &name, &code, &cmdtype)
		if err != nil {
			panic(err)
		}
		//tags := len(getVars(code))
		//TODO: convert created string to time object
		item := snipItem{
			Hash:    hash,
			Time:    time.Now(),
			Name:    name,
			Code:    code,
			CmdType: cmdtype,
		}

		snip = append(snip, item)
	}

	rows.Close() //good habit to close
	return snip
}

func dbWrite(hash string, created time.Time, title string, code string, cmdtype string) error {

	tx, _ := database.Begin()
	stmt, _ := tx.Prepare("insert into snips (hash,created,name,code,cmdtype) values (?,?,?,?,?)")
	_, err := stmt.Exec(hash, time.Now(), title, code, cmdtype)
	if err != nil {
		logError("error saving")
	}
	tx.Commit()

	return nil
}

func dbUpdatePopular(hash string) error {
	var matches int = 0
	var timesused int = 0
	var lastused string = ""

	qry := string("SELECT * FROM popular where hash=?")
	rows, err := database.Query(qry, hash)
	if err != nil {
		logError("ERROR: unable to query db")
		panic(err)
	}

	for rows.Next() {
		err = rows.Scan(&hash, &timesused, &lastused)
		if err != nil {
			panic(err)
		}
		matches++
	}
	rows.Close()

	if matches > 1 {
		//more than one match for a given hash mes something fishy is happening...
		fmt.Println("Error: too many matches found for a single hash, the offending hash is", hash, ".")
		panic("Refusing to continue")
	}

	tx, _ := database.Begin()
	if matches == 0 {
		stmt, _ := tx.Prepare("insert into popular (hash,timesused,lastused) values (?,?,?)")
		_, err = stmt.Exec(hash, 1, time.Now())
		if err != nil {
			logError("error saving popular")
		}
	} else if matches == 1 {
		timesused++
		stmt, _ := tx.Prepare("update popular set timesused=?,lastused=? where hash=?")
		_, err = stmt.Exec(timesused, time.Now(), hash)
		if err != nil {
			logError("error updating popular")
		}
	}
	tx.Commit()

	return nil
}

func dbGetAll() []snipItem {
	var snip []snipItem
	var item snipItem
	var hash string
	var name string
	var code string
	var created string
	var cmdtype string

	qry := string("SELECT * FROM snips")
	rows, err := database.Query(qry)
	if err != nil {
		logError("ERROR: unable to query db")
		panic(err)
	}

	for rows.Next() {
		err = rows.Scan(&hash, &created, &name, &code, &cmdtype)
		if err != nil {
			panic(err)
		}
		//tags := len(getVars(code))
		item = snipItem{
			Hash:    hash,
			Time:    time.Now(),
			Name:    name,
			Code:    code,
			CmdType: cmdtype,
		}
		snip = append(snip, item)
	}
	rows.Close() //good habit to close
	return snip
}
