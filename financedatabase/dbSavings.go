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

// Saving are stored in the databse
// as a row containing and id and a series of data rapresented by the struct amount
type Saving struct {
	// id primary key in the database for expences
	id *int
	// a data contained in the row
	a *Amount
}

// Equals checks if to saving are the same
func (s Saving) Equals(s2 Found) bool {
	if strings.Compare(s.String(), s2.String()) == 0 {
		return true
	}
	return false
}

// EqualValue TODO:
func (s Saving) EqualValue(s2 Amount) bool {
	return (*s.a).equals(s2)
}

// store a saving
func (s Saving) store() {
	(*s.a).storeAmount(SavingsNames["table"], SavingsNames["item"])
}

// Update updates a Saving
func (s Saving) Update(a Amount) {
	s.a.updateAmount(SavingsNames["table"], SavingsNames["item"], *s.id, a)
	s.store()
}

// Delete deltes a saving
func (s Saving) delete() {
	s.a.deleteAmount(SavingsNames["table"], SavingsNames["item"], *s.id)
}

// GetAmount returns the data of the saving
func (s Saving) GetAmount() Amount {
	return *s.a
}

// TODO:
func (s Saving) String() string {
	return fmt.Sprintf("%v %v", *s.id, (*s.a).String())
}

// Savings has multiple saving
type Savings struct {
	savings *[]Saving
}

// NewSavings TODO:
func NewSavings() Savings {
	sa := make([]Saving, 0, 0)
	return Savings{savings: &sa}
}

// Get TODO:
func (sa Savings) Get(sqlSavings string) {
	(*sa.savings) = make([]Saving, 0, 0)
	ids, am := retriveManyAmounts(sqlSavings, SavingsNames["name"])
	for i := 0; i < len(ids); i++ {
		s := Saving{id: &ids[i], a: &am[i]}
		(*sa.savings) = append(*sa.savings, s)
	}
}

// GetAll returns all the savings
func (sa Savings) GetAll() {
	sqlSavings := `
	SELECT * FROM savings
	ORDER BY name
	`
	sa.Get(sqlSavings)
}

// GetNonRecurrent returns all the savings
func (sa Savings) GetNonRecurrent() {
	sqlSavings := `
	SELECT * FROM savings
	WHERE recurrency = 0
	ORDER BY name
	`
	sa.Get(sqlSavings)
}

// Len TODO:
func (sa Savings) Len() int {
	return len(*sa.savings)
}

// ReturnElement TODO:
func (sa Savings) ReturnElement(i int) Found {
	if i >= sa.Len() {
		panic(fmt.Errorf("Too long: i = %v len = %v", i, sa.Len()))
	}
	return (*sa.savings)[i]
}

// Add a saving
func (sa Savings) Add(a Amount) {
	a.addAmount(SavingsNames["table"], SavingsNames["item"])
	sa.GetAll()
}

// Delete delete one saving
func (sa Savings) Delete(s Found) {
	for _, ss := range *sa.savings {
		if s.Equals(ss) {
			ss.delete()
			sa.GetAll()
		}
	}
}

// Sum TODO:
func (sa Savings) Sum() float32 {
	var t float32
	if sa.Len() > 0 {
		for _, s := range *sa.savings {
			t += (s.GetAmount()).sum
		}
	}
	return t
}

// String TODO:
func (sa Savings) String() string {
	if sa.Len() == 0 {
		return "there are no savings"
	}
	var buffer bytes.Buffer
	for _, s := range *sa.savings {
		buffer.WriteString(s.String())
		buffer.WriteString("\n")
	}
	return buffer.String()
}
