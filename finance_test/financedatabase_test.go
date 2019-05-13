package financetest

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/Mot93/Personal-Finance/financedatabase"
)

// >sqlite3 testfinance.db
// >.exit

// TestDatabase run all the test that require the database to exist and be initialized
// When all test are done, erase the test databse
func TestDatabase(t *testing.T) {
	// Catches any panic launched during any part of the testing
	dbName := "testfinance.db"
	fmt.Println(dbName)
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
	t.Run("Expence", func(t *testing.T) { DBFound(t, financedatabase.NewExpences(), "Expences") })
	// TODO: Update Saving to Found interface
	//t.Run("Saving", func(t *testing.T) { DBFound(t, financedatabase.NewSavings(), "Savings") })
}

// closeDB erases the tadabase after testing
func closeDB(t *testing.T, dbName string) {
	financedatabase.Close()
	err := os.Remove(dbName)
	if err != nil {
		t.Errorf("Error while erasing the test database: %v\n", err)
	}
}

// DBCategory tests all the function of Categories
func DBCategory(t *testing.T) {
	// Checking if all categories have been stored and are in alphabetical order
	var ca financedatabase.Categories
	ca.Add("work")
	ca.Add("example")
	ca.Add("utilities")
	if len(ca) != 3 || strings.Compare(ca[0].String(), "example") != 0 || strings.Compare(ca[1].String(), "utilities") != 0 {
		t.Errorf("Error while adding categories: %v\n   %v\n   %v", len(ca), ca[0].String(), ca[1].String())
	}
	// Checking if categories are delete
	ca.Delete(ca[0])
	if len(ca) != 2 || strings.Compare(ca[0].String(), "utilities") != 0 || strings.Compare(ca[1].String(), "work") != 0 {
		t.Errorf("Error while erasing categories: %v\n   %v\n   %v", len(ca), ca[0].String(), ca[1].String())
	}
}

// DBExpences tests all fuction of Expense/Saving & Expences/Savings struct and the database
func DBFound(t *testing.T, fo financedatabase.Founds, name string) {
	t.Logf("Testing %v", name)
	if fo.Len() != 0 {
		t.Errorf("Error with the base structure conteining %v", name)
	}
	// Checking that the expences are stored and retrived corectly
	var a1, a2, a3 financedatabase.Amount
	a1.Constructor("IMU", 500.0, "2018-01-16", "", "utilities", 0)
	fo.Add(a1)
	if fo.Len() != 1 {
		t.Errorf("Error adding the firs %v\n%v", name, fo.String())
	}
	a2.Constructor("Netflix", 15.0, "2018-06-15", "", "work", 0)
	fo.Add(a2)
	if fo.Len() != 2 {
		t.Errorf("Error adding the second %v", name)
	}
	a3.Constructor("Condominio", 200.0, "2018-03-17", "", "utilities", 0)
	fo.Add(a3)
	if fo.Len() != 3 {
		t.Errorf("Error adding the third %v", name)
	}
	// Alphabetical order
	if fo.Len() != 3 || !(fo.GetElement(0)).EqualValue(a3) || !(fo.GetElement(1)).EqualValue(a1) || !(fo.GetElement(2)).EqualValue(a2) {
		t.Errorf("Error adding %v len = %v:\n%v", name, fo.Len(), fo.String())
	}
	// Checking if its possible to delete
	fo.Delete(fo.GetElement(2))
	if fo.Len() != 2 || !(fo.GetElement(0)).EqualValue(a3) || !(fo.GetElement(1)).EqualValue(a1) {
		t.Errorf("Error eransing %v len = %v:\n%v", name, fo.Len(), fo.String())
	}
}
