package financedatabase

import (
	"database/sql"
	"fmt"
	"strings"

	// SQLite3 needs to be initialized before being used
	_ "github.com/mattn/go-sqlite3"
)

// Expence is a struct containing the data of an expence in similar fashion of the database
type Expence struct {
	id       int
	name     string
	sum      float32
	start    string
	end      string
	category string
}

// Constructor fills the struct excpet for the id field
// The id field is managed by the database
func (e *Expence) Constructor(name string, sum float32, start string, end string, category string) {
	e.name = name
	e.sum = sum
	e.start = start
	e.end = end
	e.category = category
}

// Equals will compare that each field of Expence are equals
// Except for the field ID
func (e Expence) Equals(e2 Expence) bool {
	if strings.Compare(e.name, e2.name) == 0 && e.sum == e2.sum && strings.Compare(e.start, e2.start) == 0 && strings.Compare(e.end, e2.end) == 0 {
		return true
	}
	return false
}

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
	db = connectdatabase(name)
	createtablecategory()
	createtableexpences()
}

// startdatabase initiate a connection with the database that already exist
// if the database doesn't exist, it creates a new one
// the path to the database is determined by filepath
func connectdatabase(filepath string) (db *sql.DB) {
	db, err := sql.Open("sqlite3", filepath)
	if err != nil {
		panic(fmt.Errorf("Erorr connecting to the database: %v", err))
	}
	if db == nil {
		panic(fmt.Errorf("Erorr connecting to the database: db connection is nil"))
	}
	return db
}

// createtable creates table, provided the righ string
func createtable(table string, tableName string) {
	_, err := db.Exec(table)
	if err != nil {
		panic(fmt.Errorf("Erorr creating table %v: %v", tableName, err))
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

// createtablecategory creates the table category if it doensn't already exist
func createtablecategory() {
	categoryTable := `
	CREATE TABLE IF NOT EXISTS categories(
		category TEXT PRIMARY KEY
	);
	`
	createtable(categoryTable, "categories")
}

// StoreCategory add a category to the table categories
func StoreCategory(category string) {
	sqlCategory := fmt.Sprintf(`
	INSERT INTO categories(
		category
	) VALUES ('%v');
	`, category)
	storeItem(sqlCategory, "category")
}

// createtableexpences creates the table exences if it doensn't already exist
func createtableexpences() {
	// Creting a string with the grave accent because the string contains \n
	// YYYY-MM-DD
	// Start at 1
	expencesTable := `
	CREATE TABLE IF NOT EXISTS expences(
		id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		sum REAL NOT NULL,
		start TEXT NOT NULL,
		end TEXT DAFAULT NULL,
		category TEXT REFERENCES categories(category)
	);
	`
	createtable(expencesTable, "expences")
}

// StoreExpence adds rows to the expences tables
func StoreExpence(ex Expence) {
	sqlExpence := fmt.Sprintf(`
	INSERT INTO expences(
		name,
		sum,
		start,
		end,
		category
	) VALUES ('%v', %v, '%v', '%v', '%v');
	`, ex.name, ex.sum, ex.start, ex.end, ex.category)
	storeItem(sqlExpence, "expence")
}

// ReadExpenceByid given an id it returns all its values
func ReadExpenceByid(id int) Expence {
	sqlExpence := fmt.Sprintf(`
	SELECT * FROM expences
	WHERE id = %v
	AND date(start)
	`, id)
	// return the result of a select
	row, err := db.Query(sqlExpence)
	if err != nil {
		panic(fmt.Errorf("Error reading expences by id: %v", err))
	}
	defer row.Close()
	// reading all element in the row
	var ex Expence
	row.Next()
	err2 := row.Scan(&ex.id, &ex.name, &ex.sum, &ex.start, &ex.end, &ex.category)
	if err2 != nil {
		panic(fmt.Errorf("Error reading expence by id: %v", err2))
	}
	return ex
}

// ReadAllExpence Return all the expences contained in the db
// Ordered by start date
func ReadAllExpence() []Expence {
	sqlExpence := `
	SELECT * FROM expences
	ORDER BY date(start)
	`
	// return the result of a select
	row, err := db.Query(sqlExpence)
	if err != nil {
		panic(fmt.Errorf("Error reading all expences: %v", err))
	}
	defer row.Close()
	// reading all element in the row
	var ex []Expence
	for row.Next() {
		var e Expence
		err2 := row.Scan(&e.id, &e.name, &e.sum, &e.start, &e.end, &e.category)
		if err2 != nil {
			panic(fmt.Errorf("Error reading all expence: %v", err2))
		}
		ex = append(ex, e)
	}
	return ex
}

// TODO: modifica expences già esistenti
