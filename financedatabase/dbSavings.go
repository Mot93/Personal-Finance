package financedatabase

import (
	"bytes"
	"fmt"
	"strings"
)

// SavingsNames avoid repetition
var SavingsNames = map[string]string{
	"item":  "saving",
	"table": "savings",
}

// createTableSavings creates the table savings if it doensn't already exist
// recurrency is an int rapresenting the day of recurrency
// if recurrency is set to -1 than the recurrency will be on a montly basis
func createTableSavings() {
	createTableAmount(SavingsNames["table"])
}

// Saving is one saving
type Saving struct {
	id *int
	a  *Amount
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

// ReturnAmount returns the data of the saving
func (s Saving) ReturnAmount() Amount {
	return *s.a
}

//
func (s Saving) String() string {
	return fmt.Sprintf("%v %v", *s.id, (*s.a).String())
}

// Savings has multiple saving
type Savings struct {
	savings *[]Saving
}

// TODO:
func NewSavings() Savings {
	sa := make([]Saving, 0, 0)
	return Savings{savings: &sa}
}

// GetAll returns all the savings
func (sa Savings) GetAll() {
	sqlSavings := `
	SELECT * FROM savings
	ORDER BY name
	`
	(*sa.savings) = make([]Saving, 0, 0)
	ids, am := retriveManyAmounts(sqlSavings, SavingsNames["name"])
	for i := 0; i < len(ids); i++ {
		s := Saving{id: &ids[i], a: &am[i]}
		(*sa.savings) = append(*sa.savings, s)
	}
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
