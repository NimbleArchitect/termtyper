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
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS snips (id INTEGER PRIMARY KEY, created INTEGER, name TEXT, code TEXT)")

	return db, true
}

// returns a Snipitem that represents the hashid from the database table
func dbgetID(hash string) Snipitem {
	var snip Snipitem
	var id int
	var name string
	var code string
	var created string

	qry := string("SELECT * FROM snips WHERE ID = " + hash)
	rows, err := database.Query(qry)
	if err != nil {
		logError("ERROR: unable to query db")
		panic(err)
	}

	for rows.Next() {
		err = rows.Scan(&id, &created, &name, &code)
		if err != nil {
			panic(err)
		}
		//tags := len(getVars(code))
		snip = Snipitem{
			ID:   id,
			Time: time.Now(),
			Name: name,
			Code: code,
		}
	}
	rows.Close() //good habit to close
	return snip
}

// search for a record name that has a wildcard match to field and return a Snipitem that represents the match
func dbfind(field string, searchfor string) []Snipitem {
	//TODO: search search for matching tags
	var snip []Snipitem
	var id int
	var name string
	var code string
	var created string

	// query
	qry := string("SELECT * FROM snips WHERE " + field + " LIKE '%" + searchfor + "%'")
	rows, err := database.Query(qry)

	if err != nil {
		logError("ERROR: unable to query db")
		panic(err)
	}

	for rows.Next() {
		err = rows.Scan(&id, &created, &name, &code)
		if err != nil {
			panic(err)
		}
		//tags := len(getVars(code))
		snipitem := Snipitem{
			ID:   id,
			Time: time.Now(),
			Name: name,
			Code: code,
		}

		snip = append(snip, snipitem)
	}

	rows.Close() //good habit to close
	return snip
}
