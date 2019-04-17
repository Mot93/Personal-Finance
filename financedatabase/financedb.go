package financedatabase

import (
	"database/sql"
	"fmt"
	"strings"

	// SQLite3 needs to be initialized before being used
	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

// InitDB makes the necessaries initialization to use the database, whithout launching this function, the db cannot work
// Those steps have been seprated from the init() because:
//		1) it has to be possible to reconnect the databse if the connection is lost
// 		2) it has to be possible to try a new connection to the database if a failure occours
//		3) both operation have to be avaliable even to the main
// if the name is an empty string "" the default db name will be finance.db
func InitDB(name string) {
	fmt.Printf("db = %v", db)
	if strings.Compare(name, "") == 0 {
		name = "finance.db"
	}
	db = connectDatabase(name)
	createTableCategory()
	createTableExpences()
	createTableSavings()
}

// startdatabase initiate a connection with the database that already exist
// if the database doesn't exist, it creates a new one
// the path to the database is determined by filepath
func connectDatabase(filepath string) (db *sql.DB) {
	db, err := sql.Open("sqlite3", filepath)
	if err != nil {
		panic(fmt.Errorf("Erorr connecting to the database: %v", err))
	}
	if db == nil {
		panic(fmt.Errorf("Erorr connecting to the database: db connection is nil"))
	}
	return db
}

// CLose TODO:
func Close() {
	db = nil
}

// executeCommand execute an sql command
func executeCommand(command string, operation string) {
	_, err := db.Exec(command)
	if err != nil {
		panic(fmt.Errorf("Erorr %v: %v", operation, err))
	}
}

// storeItem stores an item in a table, given a proper string
func storeItem(item string, itemName string) {
	statement, err := db.Prepare(item)
	if err != nil {
		panic(fmt.Errorf("Error adding item %v, preparing: %v", itemName, err))
	}
	defer statement.Close()
	_, errexe := statement.Exec()
	if errexe != nil {
		panic(fmt.Errorf("Error adding item %v, executing: %v", itemName, errexe))
	}
}
