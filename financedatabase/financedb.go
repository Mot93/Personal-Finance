package financedatabase

import (
	"database/sql"

	// SQLite3 needs to be initialized before being used
	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB
var dberror error

func init() {
	InitDB("")
}

// InitDB makes the necessaries initialization to use the database, whithout launching this function, the db cannot work
// Those steps have been seprated from the init() because:
//		1) it has to be possible to reconnect the databse if the connection is lost
// 		2) it has to be possible to try a new connection to the database if a failure occours
//		3) both operation have to be avaliable even to the main
func InitDB(name string) {
	if name == "" {
		name = "finance.db"
	}
	db, dberror = connectdatabase(name)
	createtableexpences(db)
}

// startdatabase initiate a connection with the database that already exist
// if the database doesn't exist, it creates a new one
func connectdatabase(filepath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", filepath)
	return db, err
}

// createtableexpences create the table in the database if it doensn't already exist
func createtableexpences(db *sql.DB) {
	// Creting a string with the grave accent because the stringcontains \n
	// YYYY-MM-DD
	expencesTable := `
	CREATE TABLE IF NOT EXISTS expences(
		id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		sum INTEGER NOT NULL,
		start DATE NOT NULL,
		end DATE DAFAULT NULL,
	)
	`
	_, err := db.Exec(expencesTable)
	if err != nil {
		panic(err)
	}
}
