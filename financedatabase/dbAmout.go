package financedatabase

import (
	"fmt"
	"strconv"
	"strings"
)

// Amount is a collection of data about a movement of money
// It rapresent a row in bothe the table Expences and Savings
// Amount contain the basic data needed by both Expences and Savings, where:
// 		Expences move money outside of the money pool
//		Savings move money inside of the money pool
type Amount struct {
	// name name of the transaction
	name string
	// sum of the transaction
	sum float32
	// start rapresent when the recurrency has to start
	// alternatively it rapresent the date of the transaction
	start string
	// end rapresent when the transaction stopts it's repetition
	end string
	// categories are stored from the database, never set it manually
	category string
	// recurrency rapresnet the amount of days
	// set to -1 for a monthly reccurence
	recurrency int
}

// equals will assert that all parameters of two Amounts are equal
func (a Amount) equals(a2 Amount) bool {
	if strings.Compare(a.name, a2.name) == 0 && a.sum == a2.sum && strings.Compare(a.start, a2.start) == 0 && strings.Compare(a.end, a2.end) == 0 {
		return true
	}
	return false
}

// Constructor fills the struct Amount imposing the constraint the data has to follow
func (a *Amount) Constructor(name string, sum float32, start string, end string, category string, recurrency int) {
	a.name = name
	// TODO: check the sum is > 0
	a.sum = sum
	// TODO: use date struct to get the right formatting
	a.start = start
	a.end = end
	// TODO: check if category exists in the database
	a.category = category
	// TODO: check if the value is > -2
	a.recurrency = recurrency
}

// storeAmount stores an amount inside a specified table (Expences or Savings)
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

// updateAmount update one row of Expences or Saving with the new values
// Since both the table Expences and Saving require an extra parameter called id, this parameter has to be provided by the Saving or Expence using this method, enche why the method is not public
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

// deleteAmount deletes a row from either Savings or Expence table
// Since both the table Expences and Saving require an extra parameter called id, this parameter has to be provided by the Saving or Expence using this method, enche why the method is not public
func (a Amount) deleteAmount(tableName string, itemName string, id int) {
	sqlExpence := fmt.Sprintf(`
	DELETE FROM %v
	WHERE id = %v
	`, tableName, id)
	executeCommand(sqlExpence, fmt.Sprintf("delete %v", itemName))
}

// createTableAmount creates a table containg data in the form of amount and an id (primary key)
// Used for the tables savings or expences
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

// addAmount adds a row in either an Expences or a Savings table
// tableName and itemName are needed to distinguish betwen savings and expences
func (a Amount) addAmount(tableName string, itemName string) {
	a.storeAmount(tableName, itemName)
}

// retriveSingleAmmount returns a sigle ammount from either the table Expence or Savings
// Can only be used by Expences and Savings that have to build the sql query to beforhand
// Since there could be duplicate, the id is the only way to get a single instance of Amount
// readingElement specifies if savings or expences are used
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
		panic(fmt.Errorf("Error reading %v by id: %v", readingElement, err2))
	}
	return id, a
}

// retriveManyAmounts returns all the rows stored in either Savings or Expences table
// readingElement specifies if savings or expences are used
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

// String returns a string with all the data contained in the amount
func (a Amount) String() string {
	return fmt.Sprintf("%v %v %v %v %v %v", a.name, a.sum, a.start, a.end, a.category, a.recurrency)
}

// GetStart returns the year, month and day of the starting date
func (a Amount) GetStart() (year int, month int, day int) {
	if strings.Compare(a.start, "") != 0 {
		s := strings.SplitAfter(a.start, "-")
		year, _ = strconv.Atoi(s[0])
		month, _ = strconv.Atoi(s[1])
		day, _ = strconv.Atoi(s[2])
	}
	// TODO: an else statement managing the possibility that there isn't a start date
	return year, month, day
}
