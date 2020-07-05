package main

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"time"
)

var database *sql.DB

func opendb(dbpath string) (*sql.DB, bool) {
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

	return db, true
}

// returns a Snipitem that represents the hashid from the database table
func dbgetID(hash string) Snipitem {
	var snip Snipitem

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

	for rows.Next() {
		err = rows.Scan(&hash, &created, &name, &code, &cmdtype)
		if err != nil {
			panic(err)
		}
		//tags := len(getVars(code))
		snip = Snipitem{
			Hash:    hash,
			Time:    time.Now(),
			Name:    name,
			Code:    code,
			CmdType: cmdtype,
		}
	}
	rows.Close() //good habit to close
	return snip
}

// search for a record name that has a wildcard match to field and return a Snipitem that represents the match
func dbfind(field string, searchfor string) []Snipitem {
	//TODO: search for matching tags
	var snip []Snipitem
	var hash string
	var name string
	var code string
	var created string
	var cmdtype string
	// query
	qry := string("SELECT * FROM snips WHERE " + field + " LIKE '%" + searchfor + "%'")
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
		//TODO: convert created string to time object
		snipitem := Snipitem{
			Hash:    hash,
			Time:    time.Now(),
			Name:    name,
			Code:    code,
			CmdType: cmdtype,
		}

		snip = append(snip, snipitem)
	}

	rows.Close() //good habit to close
	return snip
}
