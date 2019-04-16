package financedatabase

import (
	"database/sql"
	"fmt"
	"strings"

	// SQLite3 needs to be initialized before being used
	_ "github.com/mattn/go-sqlite3"
)

// Amount contains the data of both savings and expences
// The two have identical structur but serve two differen purpose
type Amount struct {
	id         int
	name       string
	sum        float32
	start      string
	end        string
	category   string
	recurrency int
}

// Constructor fills the struct Amount except for the id field
// The id field is managed by the database
func (a *Amount) Constructor(name string, sum float32, start string, end string, category string) {
	a.name = name
	a.sum = sum
	a.start = start
	a.end = end
	a.category = category
	a.recurrency = 0
}

// Equals will compare that each field of Expence are equals
// Except for the field ID
func (a Amount) Equals(a2 Amount) bool {
	if strings.Compare(a.name, a2.name) == 0 && a.sum == a2.sum && strings.Compare(a.start, a2.start) == 0 && strings.Compare(a.end, a2.end) == 0 {
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

//
func retriveSingleAmount(sqlQuery string, readingElement string) Amount {
	// return the result of a select
	row, err := db.Query(sqlQuery)
	if err != nil {
		panic(fmt.Errorf("Error reading %v by id: %v", readingElement, err))
	}
	defer row.Close()
	// reading all element in the row
	var am Amount
	row.Next()
	err2 := row.Scan(&am.id, &am.name, &am.sum, &am.start, &am.end, &am.category, &am.recurrency)
	if err2 != nil {
		panic(fmt.Errorf("Error reading saving by id: %v", err2))
	}
	return am
}

//
func retriveAllAmounts(sqlQuery string, readingElement string) []Amount {
	// return the result of a select
	row, err := db.Query(sqlQuery)
	if err != nil {
		panic(fmt.Errorf("Error reading all %v: %v", readingElement, err))
	}
	defer row.Close()
	// reading all element in the row
	var am []Amount
	for row.Next() {
		var a Amount
		err2 := row.Scan(&a.id, &a.name, &a.sum, &a.start, &a.end, &a.category, &a.recurrency)
		if err2 != nil {
			panic(fmt.Errorf("Error reading all expence: %v", err2))
		}
		am = append(am, a)
	}
	return am
}

// createtablecategory creates the table category if it doensn't already exist
func createTableCategory() {
	categoryTable := `
	CREATE TABLE IF NOT EXISTS categories(
		category TEXT PRIMARY KEY
	);
	`
	executeCommand(categoryTable, "create table categories")
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

// TODO: add delete categories

// ReadAllCategories return the complete list of categories
func ReadAllCategories() []string {
	sqlCategories := `
	SELECT * FROM categories
	ORDER BY category
	`
	// return the result of a select
	row, err := db.Query(sqlCategories)
	if err != nil {
		panic(fmt.Errorf("Error reading all categories: %v", err))
	}
	defer row.Close()
	// reading all element in the row
	var categories []string
	for row.Next() {
		var category string
		err2 := row.Scan(&category)
		if err2 != nil {
			panic(fmt.Errorf("Error reading all categories: %v", err2))
		}
		categories = append(categories, category)
	}
	return categories
}

// TODO: delete expence

// createTableExpences creates the table exences if it doensn't already exist
func createTableExpences() {
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
		category TEXT REFERENCES categories(category),
		recurrency INTEGER DEFAULT 0
	);
	`
	executeCommand(expencesTable, "create table expences")
}

// StoreExpence adds rows to the expences tables
func StoreExpence(am Amount) {
	sqlExpence := fmt.Sprintf(`
	INSERT INTO expences(
		name,
		sum,
		start,
		end,
		category
	) VALUES ('%v', %v, '%v', '%v', '%v');
	`, am.name, am.sum, am.start, am.end, am.category)
	storeItem(sqlExpence, "expence")
}

// ReadExpenceByid given an id it returns all its values
func ReadExpenceByid(id int) Amount {
	sqlExpence := fmt.Sprintf(`
	SELECT * FROM expences
	WHERE id = %v
	AND date(start)
	`, id)
	return retriveSingleAmount(sqlExpence, "expence")
}

// ReadAllExpences return all the expences contained in the db
// Ordered by start from the eldest to the most recent
func ReadAllExpences() []Amount {
	sqlExpences := `
	SELECT * FROM expences
	ORDER BY date(start)
	`
	return retriveAllAmounts(sqlExpences, "expences")
}

// UpdateExpence allows to update and expence
// The only way to set the id of an expence is to get an expence from the database
// This is intended to make sure that we are operating on the expences specified by the user
func UpdateExpence(a Amount) {
	sqlExpence := fmt.Sprintf(`
	UPDATE expences 
	SET name = '%v', sum = %v, start = '%v', end = '%v', category = '%v'
	WHERE id = %v;
	`, a.name, a.sum, a.start, a.end, a.category, a.id)
	executeCommand(sqlExpence, "updating expences")
}

// TODO: add savings & recurent savings

// createTableSavings creates the table savings if it doensn't already exist
func createTableSavings() {
	savingsTable := `
	CREATE TABLE IF NOT EXISTS savings(
		id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		sum REAL NOT NULL,
		start TEXT NOT NULL,
		end TEXT DAFAULT NULL,
		category TEXT REFERENCES categories(category),
		recurrency INTEGER DEFAULT 0
	);
	`
	executeCommand(savingsTable, "create table savings")
}

// ReadSavingByid return the saving rfered by id
func ReadSavingByid(id int) Amount {
	sqlSaving := fmt.Sprintf(`
	SELECT * FROM savings
	WHERE id = %v
	`, id)
	return retriveSingleAmount(sqlSaving, "saving")
}

// ReadAllSavings return all the expences contained in the db
// Ordered by start from the eldest to the most recent
func ReadAllSavings() []Amount {
	sqlSavings := `
	SELECT * FROM savings
	ORDER BY date(start)
	`
	return retriveAllAmounts(sqlSavings, "savings")
}
