package financetest

import (
	"os"
	"testing"

	"github.com/Mot93/personalfinance/financedatabase"
)

// >sqlite3 testfinance.db
// >.exit

// TestAmount test all function working with the Amount struct
func TestAmount(t *testing.T) {
	var a1, a2 financedatabase.Amount
	a1.Constructor("netflix", 15.0, "2017-06-22", "", "work")
	a2.Constructor("netflix", 15.0, "2017-06-22", "", "work")
	if !a1.Equals(a2) {
		t.Errorf("Equals of Expence doesn't work should be equals:\n   %v\n   %v", a1, a2)
	}
	a2.Constructor("Rent", 400.0, "2006-05-15", "", "utilities")
	if a1.Equals(a2) {
		t.Errorf("Equals of Expence doesn't work should be differen:\n   %v\n   %v", a1, a2)
	}
}

// TestMain run all the test that require the database to exist and be initialized
// When all test are done, erase the test databse
func TestDatabase(t *testing.T) {
	// Catches any panic launched during any part of the testing
	dbName := "testfinance.db"
	defer func() {
		if rec := recover(); rec != nil {
			closeDB(t, dbName)
			t.Errorf("Error in db testing: %v\n", rec)
		}
	}()
	// Initializing the database
	financedatabase.InitDB(dbName)
	// Erasing the test database after it's usage
	defer closeDB(t, dbName)
	// TESTS
	t.Run("Category", func(t *testing.T) { DBCategory(t) })
	t.Run("Expences", func(t *testing.T) { DBExpences(t) })
}

//
func closeDB(t *testing.T, dbName string) {
	err := os.Remove(dbName)
	if err != nil {
		t.Errorf("Error while erasing the test database: %v\n", err)
	}
}

// DBCategory
func DBCategory(t *testing.T) {
	financedatabase.StoreCategory("work")
	financedatabase.StoreCategory("utilities")
	categories := financedatabase.ReadAllCategories()
	if categories[0] != "utilities" || categories[1] != "work" {
		t.Errorf("List of utlities incorrect:\n   utilities %v\n   work %v", categories[0], categories[1])
	}
}

// DBExpences tests all fuction working with the Expense struct and the database
// Launched by TestDatabase
func DBExpences(t *testing.T) {
	var e []financedatabase.Amount
	var e1, e2 financedatabase.Amount
	// Checking that the expences are stored and retrived corectly
	e1.Constructor("netflix", 15.0, "2017-06-22", "", "work")
	financedatabase.StoreExpence(e1)
	e2 = financedatabase.ReadExpenceByid(1)
	if !e1.Equals(e2) {
		t.Errorf("Expence wasn't stored corectly:\n   %v\n   %v", e1, e2)
	}
	// Checking if all the expences are retrived corectly
	e2.Constructor("Rent", 400.0, "2006-05-15", "", "utilities")
	financedatabase.StoreExpence(e2)
	e = financedatabase.ReadAllExpences()
	if !e1.Equals(e[1]) || !e2.Equals(e[0]) {
		t.Errorf("All Expences weren't returned in a proper way:\n   e1   = %v\n   e[1] = %v\n   e2   = %v\n   e[0] = %v", e1, e[1], e2, e[0])
	}
	if e1.Equals(e[0]) || e2.Equals(e[1]) {
		t.Errorf("All Expences weren't returned in a proper way:\n   e1   = %v\n   e[1] = %v\n   e2   = %v\n   e[0] = %v", e1, e[1], e2, e[0])
	}
	// Checking that expences are updated corectly
	e1 = financedatabase.ReadExpenceByid(1)
	e1.Constructor("new expence", 66, "1993-06-02", "", "work")
	financedatabase.UpdateExpence(e1)
	e2 = financedatabase.ReadExpenceByid(1)
	if e1 != e2 {
		t.Errorf("Expences updated failed:\n   %v\n   %v", e1, e2)
	}
}
