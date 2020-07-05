import uuid
import sqlite3
from sqlite3 import Error


def create_connection(db_file):
    """ create a database connection to the SQLite database
        specified by the db_file
    :param db_file: database file
    :return: Connection object or None
    """
    conn = None
    try:
        conn = sqlite3.connect(db_file)
    except Error as e:
        print(e)

    return conn


def select_all_tasks(conn_old, conn_new):
    cur = conn_old.cursor()
    cur.execute("SELECT created, name, code FROM snips")

    i = 0
    for created, name, code in cur.fetchall():
        i = i + 1
        write_new_row(conn_new, i, created, name, code)
    print("\nimported", i, "records")

def create_new_db(conn):
    create_table_sql = """ CREATE TABLE "snips" (
	"hash"	TEXT UNIQUE,
	"created"	INTEGER,
	"name"	TEXT,
	"code"	TEXT,
	"cmdtype"	TEXT,
	PRIMARY KEY("hash")
); """

    try:
        c = conn.cursor()
        c.execute(create_table_sql)
    except Error as e:
        print(e)
        quit()
    

def write_new_row(conn, id, created, name, code):
    try:        
        cursor = conn.cursor()

        sqlite_insert_query = """ INSERT INTO snips
                            (hash, created, name, code, cmdtype) 
                            VALUES 
                            (?,?,?,?,?); """

        hash = str(uuid.uuid4())
        cmdtype = "bash"
        count = cursor.execute(sqlite_insert_query, (hash, created, name, code, cmdtype) )
        conn.commit()
        print(".", sep="", end="")

    except sqlite3.Error as error:
        print("Failed to insert data into sqlite table", error)



def main():
    from pathlib import Path
    home = str(Path.home())

    database_old = r"" + home + "/.termtyper/snippets.db"
    database_new = r"" + home + "/.termtyper/termtyper.db"

    # create a database connection
    conn_old = create_connection(database_old)
    conn_new = create_connection(database_new)
    create_new_db(conn_new)

    with conn_old:
        select_all_tasks(conn_old, conn_new)


if __name__ == '__main__':
    main()