package financedatabase

import (
	"bytes"
	"fmt"
	"strings"
)

// ExpencesNames avoid repetition
var ExpencesNames = map[string]string{
	"item":  "expence",
	"table": "expences",
}

// createTableExpences creates the table exences if it doensn't already exist
func createTableExpences() {
	createTableAmount(ExpencesNames["table"])
}

// Expence is an amount
type Expence struct {
	id *int
	a  *Amount
}

// Equals check id they are the same
func (e Expence) Equals(e2 Found) bool {
	if strings.Compare(e.String(), e2.String()) == 0 {
		return true
	}
	return false
}

// EqualValue check if the have the same values
// They might be different Expences
func (e Expence) EqualValue(a Amount) bool {
	return (*e.a).equals(a)
}

// store adds rows to the expences tables
func (e Expence) store() {
	(*e.a).storeAmount(ExpencesNames["table"], ExpencesNames["name"])
}

// Update allows to modify the expence
func (e Expence) Update(a Amount) {
	(*e.a).updateAmount(ExpencesNames["table"], ExpencesNames["name"], *e.id, a)
	e.store()
}

// delete delete an expence from expences
func (e Expence) delete() {
	e.a.deleteAmount(ExpencesNames["table"], ExpencesNames["name"], *e.id)
}

// ReturnAmount returns the amount struct containing the expences data
func (e Expence) ReturnAmount() Amount {
	return *e.a
}

// String print a string with all
func (e Expence) String() string {
	if e.id == nil || e.a == nil {
		return "empty"
	}
	return fmt.Sprintf("%v %v", *e.id, (*e.a).String())
}

// Expences conteins multiple Expence
type Expences struct {
	expences *[]Expence
}

// NewExpences create a new Expences struct
func NewExpences() Expences {
	ex := make([]Expence, 0, 0)
	return Expences{expences: &ex}
}

// Get TODO:
func (ex Expences) Get(sqlExpences string) {
	*ex.expences = make([]Expence, 0, 0)
	ids, am := retriveManyAmounts(sqlExpences, ExpencesNames["name"])
	for i := 0; i < len(ids); i++ {
		e := Expence{id: &(ids[i]), a: &(am[i])}
		//e.constructor(ids[i], am[i])
		*ex.expences = append(*ex.expences, e)
	}
}

// GetAll return all the expences contained in the db
// Ordered by start from the eldest to the most recent
func (ex Expences) GetAll() {
	sqlExpences := `
	SELECT * FROM expences
	ORDER BY name
	`
	ex.Get(sqlExpences)
}

// GetNonRecurrent TODO:
func (ex Expences) GetNonRecurrent() {
	sqlExpences := `
	SELECT * FROM expences
	WHERE recurrency = 0
	ORDER BY name
	`
	ex.Get(sqlExpences)
}

// Len return the number of Expence
func (ex Expences) Len() int {
	return len(*ex.expences)
}

// ReturnElement returns an expence at the specified position
func (ex Expences) ReturnElement(i int) Found {
	return (*ex.expences)[i]
}

// Add adds an Expence to the database
func (ex Expences) Add(a Amount) {
	a.addAmount(ExpencesNames["table"], ExpencesNames["name"])
	ex.GetAll()
}

// Delete an expence
func (ex Expences) Delete(e Found) {
	for _, ee := range *ex.expences {
		if e.Equals(ee) {
			ee.delete()
			ex.GetAll()
		}
	}
}

// Sum TODO:
func (ex Expences) Sum() float32 {
	var t float32 = 0.0
	if ex.Len() > 0 {
		for _, e := range *ex.expences {
			t += (e.ReturnAmount()).sum
		}
	}
	return t
}

// String TODO: all the documentation
func (ex Expences) String() string {
	if ex.Len() == 0 {
		return " there are no expences "
	}
	var buffer bytes.Buffer
	for _, e := range *ex.expences {
		buffer.WriteString(e.String())
		buffer.WriteString("\n")
	}
	return buffer.String()
}
