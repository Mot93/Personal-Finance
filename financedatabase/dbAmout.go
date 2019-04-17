package financedatabase

import (
	"fmt"
	"strings"
)

// amount contains the data of both savings and expences
// The two have identical structur but serve two differen purpose
type Amount struct {
	name       string
	sum        float32
	start      string
	end        string
	category   string
	recurrency int
}

// equals will compare that each field of Expence are equals
// Except for the field ID
func (a Amount) equals(a2 Amount) bool {
	if strings.Compare(a.name, a2.name) == 0 && a.sum == a2.sum && strings.Compare(a.start, a2.start) == 0 && strings.Compare(a.end, a2.end) == 0 {
		return true
	}
	return false
}

func (a *Amount) Constructor(name string, sum float32, start string, end string, category string, recurrency int) {
	a.name = name
	a.sum = sum
	a.start = start
	a.end = end
	a.category = category
	a.recurrency = recurrency
}

func (a Amount) storeAmount(tableName string, itemName string) {
	sqlAmount := fmt.Sprintf(`
	INSERT INTO %v (
		name, 
		sum, 
		start, 
		end, 
		category, 
		recurrency
	) VALUES ('%v', %v, '%v', '%v', '%v', %v);
	`, tableName, a.name, a.sum, a.start, a.end, a.category, a.recurrency)
	storeItem(sqlAmount, itemName)
}

// updateAmount update one saving or expence
func (a Amount) updateAmount(tableName string, itemName string, id int, a2 Amount) {
	a.name = a2.name
	a.sum = a2.sum
	a.start = a2.start
	a.end = a2.end
	a.category = a2.category
	a.recurrency = a2.recurrency
	sqlUpdate := fmt.Sprintf(`
	UPDATE %v 
	SET name = '%v', sum = %v, start = '%v', end = '%v', category = '%v', recurrency =%v
	WHERE id = %v;
	`, tableName, a.name, a.sum, a.start, a.end, a.category, a.recurrency, id)
	executeCommand(sqlUpdate, fmt.Sprintf("updating %v", itemName))
}

// deleteAmountByid deletes either an expence or a saving
func (a Amount) deleteAmount(tableName string, itemName string, id int) {
	sqlExpence := fmt.Sprintf(`
	DELETE FROM %v
	WHERE id = %v
	`, tableName, id)
	executeCommand(sqlExpence, fmt.Sprintf("delete %v", itemName))
}

// createTableAmount creates either the table for savings or expences
func createTableAmount(tableName string) {
	// Creting a string with the grave accent because the string contains \n
	// YYYY-MM-DD
	// Start at 1
	expencesTable := fmt.Sprintf(`
	CREATE TABLE IF NOT EXISTS %v(
		id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		sum REAL NOT NULL,
		start TEXT DEFAULT NULL,
		end TEXT DAFAULT NULL,
		category TEXT REFERENCES categories(category),
		recurrency INTEGER DEFAULT 0
	);
	`, tableName)
	executeCommand(expencesTable, fmt.Sprintf("create table %v", tableName))
}

// addAmount adds either an Expence or a Saving to the db
func (a Amount) addAmount(tableName string, itemName string) {
	a.storeAmount(tableName, itemName)
}

// retriveSingleAmmount returns a sigle ammount (saving or expence)
func retriveSingleAmount(sqlQuery string, readingElement string) (id int, a Amount) {
	// return the result of a select
	row, err := db.Query(sqlQuery)
	if err != nil {
		panic(fmt.Errorf("Error reading %v by id: %v", readingElement, err))
	}
	defer row.Close()
	// reading all element in the row
	row.Next()
	err2 := row.Scan(&id, &a.name, &a.sum, &a.start, &a.end, &a.category, &a.recurrency)
	if err2 != nil {
		panic(fmt.Errorf("Error reading %d by id: %v", readingElement, err2))
	}
	return id, a
}

// retriveManyAmounts returns all the element stored in either savings or expences
func retriveManyAmounts(sqlQuery string, readingElement string) (ids []int, am []Amount) {
	// return the result of a select
	row, err := db.Query(sqlQuery)
	if err != nil {
		panic(fmt.Errorf("Error reading all %v: %v", readingElement, err))
	}
	defer row.Close()
	// reading all element in the row
	for row.Next() {
		var a Amount
		var id int
		err2 := row.Scan(&id, &a.name, &a.sum, &a.start, &a.end, &a.category, &a.recurrency)
		if err2 != nil {
			panic(fmt.Errorf("Error reading all expence: %v", err2))
		}
		am = append(am, a)
		ids = append(ids, id)
	}
	return ids, am
}

// String to print a string with all teh values
func (a Amount) String() string {
	return fmt.Sprintf("%v %v %v %v %v %v", a.name, a.sum, a.start, a.end, a.category, a.recurrency)
}
