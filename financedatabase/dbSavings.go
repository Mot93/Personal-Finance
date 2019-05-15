package financedatabase

import (
	"bytes"
	"fmt"
	"strings"
)

// SavingsNames is a set of names that assure consistency
// Make it esier to change names when the database is modified
var SavingsNames = map[string]string{
	"item":  "saving",
	"table": "savings",
}

// createTableSavings creates the table savings if it doensn't already exist
func createTableSavings() {
	createTableAmount(SavingsNames["table"])
}

// Saving are stored in the database
// as a row containing and id and a series of data rapresented by the struct Amount
type Saving struct {
	// id primary key in the database for the table savings
	id *int
	// a data contained in the row
	a *Amount
}

// Equals checks if the id and the values are the same
// The id is never set manually, is given by the database
func (s Saving) Equals(s2 Found) bool {
	if strings.Compare(s.String(), s2.String()) == 0 {
		return true
	}
	return false
}

// EqualValue checks if the have the same amount
// They might be different Expences
func (s Saving) EqualValue(s2 Amount) bool {
	return (*s.a).equals(s2)
}

// store adds rows to the expences tables
// Only used by the Expences struct
func (s Saving) store() {
	(*s.a).storeAmount(SavingsNames["table"], SavingsNames["item"])
}

// Update allows to modify the saving
// Can only be called on an existing Saving (already extracted from the database)
func (s Saving) Update(a Amount) {
	s.a.updateAmount(SavingsNames["table"], SavingsNames["item"], *s.id, a)
	s.store()
}

// delete an expence from the saving table
func (s Saving) delete() {
	s.a.deleteAmount(SavingsNames["table"], SavingsNames["item"], *s.id)
}

// GetAmount returns the struct Amount of the Saving
func (s Saving) GetAmount() Amount {
	return *s.a
}

// String return a string containing the id and the string of Amount
func (s Saving) String() string {
	return fmt.Sprintf("%v %v", *s.id, (*s.a).String())
}

// Savings is a collection (slice) of Saving
type Savings struct {
	savings []Saving
}

// NewSavings create a new Savings struct
// Slices need initialization
func NewSavings() Savings {
	sa := make([]Saving, 0, 0)
	return Savings{savings: sa}
}

// get return Savings from the database
// Wich Saving are returned are decided by the sqlSavings string
func (sa Savings) get(sqlSavings string) {
	sa.savings = make([]Saving, 0, 0)
	ids, am := retriveManyAmounts(sqlSavings, SavingsNames["name"])
	for i := 0; i < len(ids); i++ {
		s := Saving{id: &ids[i], a: &am[i]}
		sa.savings = append(sa.savings, s)
	}
}

// GetAll fill Expences with all the Saving contained in the database
// Ordered by most recent to eldest
func (sa Savings) GetAll() {
	sqlSavings := `
	SELECT * FROM savings
	ORDER BY name
	`
	sa.get(sqlSavings)
}

// GetRecurrent fill Savings with all the recurring Saving
func (sa Savings) GetRecurrent() {
	sqlSavings := `
	SELECT * FROM savings
	WHERE recurrency != 0
	ORDER BY name
	`
	sa.get(sqlSavings)
}

// GetNonRecurrent fill Savings with all the non recurring Saving
func (sa Savings) GetNonRecurrent() {
	sqlSavings := `
	SELECT * FROM savings
	WHERE recurrency = 0
	ORDER BY name
	`
	sa.get(sqlSavings)
}

// Len return the number of Saving in Savings
func (sa Savings) Len() int {
	return len(sa.savings)
}

// GetElement returns the Expence at the specified position
func (sa Savings) GetElement(i int) Found {
	// TODO: check that: -1 < i < Savings.Len()
	return sa.savings[i]
}

// Add adds a Saving to the database
func (sa Savings) Add(a Amount) {
	// TODO: check that the amount is not empty
	a.addAmount(SavingsNames["table"], SavingsNames["item"])
}

// Delete a Saving from the db and the struct
func (sa Savings) Delete(s Found) {
	for i := 0; i < sa.Len(); i++ {
		// Removing the Saving from the databse
		sa.savings[i].delete()
		// Removing the Expence from the struct
		sa.savings = append(sa.savings[0:i], sa.savings[i+1:]...)
	}
	// TODO: return error when the Saving is not present in the Savings
}

// Sum return the sum of all the Saving in Savings
func (sa Savings) Sum() float32 {
	var t float32
	if sa.Len() > 0 {
		for _, s := range sa.savings {
			t += (s.GetAmount()).sum
		}
	}
	return t
}

// String returns a string containing all the string of all the Saving contained in the struct
func (sa Savings) String() string {
	if sa.Len() == 0 {
		return "there are no savings"
	}
	var buffer bytes.Buffer
	for _, s := range sa.savings {
		buffer.WriteString(s.String())
		buffer.WriteString("\n")
	}
	return buffer.String()
}
