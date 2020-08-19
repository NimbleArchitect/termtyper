package main

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"strings"
	"time"
)

var localDbList []*sql.DB

func dbOpen(dbpath string) (*sql.DB, bool) {
	logInfo("* open: " + dbpath)
	db, err := sql.Open("sqlite3", dbpath)

	if err != nil {
		logError("unable to open database")
	}
	ok := db.Ping()
	if ok != nil {
		panic(ok)
	}
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS "snips" (
		"hash"	    TEXT UNIQUE,
		"created"	INTEGER,
		"name"	    TEXT,
		"code"	    TEXT,
		"cmdtype"	TEXT,
		"summary"   TEXT,
		"dead"      INTEGER DEFAULT 0 NOT NULL,
		PRIMARY KEY("hash")
	);`,
	)

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS "popular" (
		"hash"	TEXT UNIQUE,
		"timesused"	INTEGER,
		"lastused"	INTEGER,
		PRIMARY KEY("hash")
	);`,
	)

	_, _ = db.Exec(`ALTER TABLE snips ADD summary TEXT;`)

	_, _ = db.Exec(`ALTER TABLE snips ADD dead INTEGER DEFAULT 0 NOT NULL;`)

	_, _ = db.Exec(`DROP VIEW IF EXISTS textsearch; CREATE VIEW textsearch
	AS
	SELECT 
		hash,
		name || ' ' || code AS search,
		created,	
		name,
		code,
		cmdtype,
		summary
	FROM snips WHERE dead=0;
	`)
	return db, true
}

// returns a Snipitem that represents the hashid from the database table
func dbGetID(database *sql.DB, hash string) (snipItem, int) {
	var snip snipItem

	var name string = ""
	var code string = ""
	var created string = ""
	var cmdtype string = ""

	qry := string("SELECT hash,created,name,code,cmdtype FROM snips WHERE hash=? AND dead=0")
	rows, err := database.Query(qry, hash)
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
		t, err := time.Parse(time.RFC3339, created)
		if err != nil {
			t = time.Now()
		}
		snip = snipItem{
			Hash:    hash,
			Time:    t,
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
func dbFind(database *sql.DB, field string, searchfor string, rowStart int) []snipItem {
	//TODO: search for matching tags
	var snip []snipItem
	var hash string
	var name string = ""
	var code string = ""
	var created string
	var cmdtype string
	var dbSummary sql.NullString
	var strSummary string
	var qry string = ""
	var strRowStart string = ""
	var err error
	var rows *sql.Rows

	if rowStart > 0 {
		strRowStart = "OFFSET ?"
	}
	// query
	if field == "name" {
		qry = string(`
		SELECT
			hash,created,name,code,cmdtype,summary 
		FROM textsearch 
		WHERE name LIKE ? LIMIT ?` + strRowStart)
	}
	if field == "code" {
		qry = string(`
		SELECT
			hash,created,name,code,cmdtype,summary 
		FROM textsearch 
		WHERE code LIKE ? LIMIT ?` + strRowStart)
	}
	if field == "all" {
		qry = string(`
		SELECT 
			hash,created,name,code,cmdtype,summary 
		FROM textsearch 
		WHERE search LIKE ? LIMIT ?` + strRowStart)
	}

	//FIXME:  the next two lines are really silly and need a proper fix, im just too lazy atm
	searchfor = strings.ReplaceAll(searchfor, " ", "%")
	logDebug("F:dbFind:searchfor =", searchfor)

	if strRowStart == "" {
		rows, err = database.Query(qry, "%"+searchfor+"%", settings.Termtyper.maxRows)
	} else {
		rows, err = database.Query(qry, "%"+searchfor+"%", settings.Termtyper.maxRows, rowStart)
	}

	if err != nil {
		logError("ERROR: unable to query db")
		panic(err)
	}

	for rows.Next() {
		err = rows.Scan(&hash, &created, &name, &code, &cmdtype, &dbSummary)
		if err != nil {
			panic(err)
		}
		if dbSummary.Valid {
			strSummary = dbSummary.String
		} else {
			strSummary = ""
		}
		//tags := len(getVars(code))
		//convert created string to time object
		t, err := time.Parse(time.RFC3339, created)
		if err != nil {
			t = time.Now()
		}

		item := snipItem{
			Hash:    hash,
			Time:    t,
			Name:    name,
			Code:    code,
			CmdType: cmdtype,
			Summary: strSummary,
		}

		snip = append(snip, item)
	}

	rows.Close() //good habit to close
	return snip
}

func dbWrite(database *sql.DB, hash string, created time.Time, title string, code string, cmdtype string, summary string) error {
	//TODO: check values are valid before saving
	t := string(created.Format(time.RFC3339))

	tx, _ := database.Begin()
	stmt, _ := tx.Prepare("insert into snips (hash,created,name,code,cmdtype,summary,dead) values (?,?,?,?,?,?,0)")
	_, err := stmt.Exec(hash, t, title, code, cmdtype, summary)
	if err != nil {
		logError("error saving")
	}
	tx.Commit()

	return nil
}

func dbDelete(database *sql.DB, hash string) error {
	//TODO: check values are valid before saving

	tx, _ := database.Begin()
	stmt, _ := tx.Prepare("UPDATE snips set dead=1 WHERE hash=?")
	_, err := stmt.Exec(hash)
	if err != nil {
		logError("error deleteing record")
	}
	tx.Commit()

	return nil
}

func dbUpdate(database *sql.DB, hash string, created time.Time, title string, code string, cmdtype string, summary string) error {
	//TODO: check values are valid before saving
	t := string(created.Format(time.RFC3339))

	tx, _ := database.Begin()
	stmt, _ := tx.Prepare("UPDATE snips set created=?,name=?,code=?,cmdtype=?,summary=?,dead=0 WHERE hash=?")
	_, err := stmt.Exec(t, title, code, cmdtype, summary, hash)
	if err != nil {
		logError("error updating record")
	}
	tx.Commit()

	return nil
}

func dbUpdatePopular(database *sql.DB, hash string) error {
	var matches int = 0
	var timesused int = 0
	var lastused string = ""

	qry := string("SELECT * FROM popular WHERE hash=? AND dead=0")
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
		//more than one match for a given hash meas something fishy is happening...
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

func dbGetPopular(database *sql.DB, count int) []snipItem {
	var snip []snipItem
	var item snipItem
	var hash string
	var name string
	var code string
	var created string
	var cmdtype string
	var dbSummary sql.NullString
	var strSummary string

	qry := string(
		`SELECT popular.hash, lastused, name, code, cmdtype, summary
		FROM popular
		INNER JOIN snips ON popular.hash = snips.hash 
		ORDER BY timesused DESC LIMIT ?;`)

	if count <= 0 {
		return snip
	}

	rows, err := database.Query(qry, count)
	if err != nil {
		logError("ERROR: unable to query db")
		panic(err)
	}

	for rows.Next() {
		err = rows.Scan(&hash, &created, &name, &code, &cmdtype, &dbSummary)
		if err != nil {
			panic(err)
		}
		if dbSummary.Valid {
			strSummary = dbSummary.String
		} else {
			strSummary = ""
		}
		t, err := time.Parse(time.RFC3339, created)
		if err != nil {
			t = time.Now()
		}
		item = snipItem{
			Hash:    hash,
			Time:    t,
			Name:    name,
			Code:    code,
			CmdType: cmdtype,
			Summary: strSummary,
		}
		snip = append(snip, item)
	}
	rows.Close() //good habit to close
	return snip

}

func dbGetAll(database *sql.DB) []snipItem {
	var snip []snipItem
	var item snipItem
	var hash string
	var name string
	var code string
	var created string
	var cmdtype string
	var dbSummary sql.NullString
	var strSummary string

	qry := string("SELECT hash,created,name,code,cmdtype,summary FROM snips WHERE dead<>1")
	rows, err := database.Query(qry)
	if err != nil {
		logError("ERROR: unable to query db")
		panic(err)
	}

	for rows.Next() {
		err = rows.Scan(&hash, &created, &name, &code, &cmdtype, &dbSummary)
		if err != nil {
			panic(err)
		}
		if dbSummary.Valid {
			strSummary = dbSummary.String
		} else {
			strSummary = ""
		}
		//tags := len(getVars(code))
		t, err := time.Parse(time.RFC3339, created)
		if err != nil {
			t = time.Now()
		}
		item = snipItem{
			Hash:    hash,
			Time:    t,
			Name:    name,
			Code:    code,
			CmdType: cmdtype,
			Summary: strSummary,
		}
		snip = append(snip, item)
	}
	rows.Close() //good habit to close
	return snip
}
