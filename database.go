package main

import (
	"database/sql"
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

	return db, true
}

// returns a Snipitem that represents the hashid from the database table
func dbGetID(hash string) (Snipitem, int) {
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
	count := 0
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
		count += 1
	}
	rows.Close() //good habit to close

	return snip, count
}

// search for a record name that has a wildcard match to field and return a Snipitem that represents the match
func dbFind(field string, searchfor string) []Snipitem {
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

func dbGetAll() []Snipitem {
	var snip []Snipitem
	var snipitem Snipitem
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
		snipitem = Snipitem{
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
