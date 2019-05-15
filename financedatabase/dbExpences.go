package financedatabase

import (
	"bytes"
	"fmt"
	"strings"
)

// ExpencesNames is a set of names that assure consistency
// Make it esier to change names when the database is modified
var ExpencesNames = map[string]string{
	"item":  "expence",
	"table": "expences",
}

// createTableExpences creates the table exences if it doensn't already exist
func createTableExpences() {
	createTableAmount(ExpencesNames["table"])
}

// Expence are stored in the database
// as a row containing and id and a series of data rapresented by the struct Amount
type Expence struct {
	// id primary key in the database for the table expences
	id *int
	// a data contained in the row
	a *Amount
}

// Equals checks if the id and the values are the same
// The id is never set manually, is given by the database
func (e Expence) Equals(e2 Found) bool {
	if strings.Compare(e.String(), e2.String()) == 0 {
		return true
	}
	return false
}

// EqualValue checks if the have the same amount
// They might be different Expences
func (e Expence) EqualValue(a Amount) bool {
	return (*e.a).equals(a)
}

// store adds rows to the expences tables
// Only used by the Expences struct
func (e Expence) store() {
	(*e.a).storeAmount(ExpencesNames["table"], ExpencesNames["name"])
}

// Update allows to modify the expence
// Can only be called on an existing Exence (already extracted from the database)
func (e Expence) Update(a Amount) {
	// TODO: checks if the expence is in the db
	// only existing expences can be updated
	(*e.a).updateAmount(ExpencesNames["table"], ExpencesNames["name"], *e.id, a)
	e.store()
}

// delete an expence from the expence table
func (e Expence) delete() {
	e.a.deleteAmount(ExpencesNames["table"], ExpencesNames["name"], *e.id)
}

// GetAmount returns the struct Amount of the Expence
func (e Expence) GetAmount() Amount {
	return *e.a
}

// String return a string containing the id and the string of Amount
func (e Expence) String() string {
	if e.id == nil || e.a == nil {
		return "empty"
	}
	return fmt.Sprintf("%v %v", *e.id, (*e.a).String())
}

// Expences is a collection (slice) of Expence
type Expences struct {
	// expences is a slice, it already use a pointer
	expences []Expence
}

// NewExpences create a new Expences struct
// Slices need initialization
func NewExpences() Expences {
	ex := make([]Expence, 0, 0)
	return Expences{expences: ex}
}

// get return Expences from the database
// Wich Expences are returned are decided by the sqlExpences string
func (ex Expences) get(sqlExpences string) {
	ex.expences = make([]Expence, 0, 0)
	ids, am := retriveManyAmounts(sqlExpences, ExpencesNames["name"])
	for i := 0; i < len(ids); i++ {
		e := Expence{id: &(ids[i]), a: &(am[i])}
		//e.constructor(ids[i], am[i])
		ex.expences = append(ex.expences, e)
	}
}

// GetAll fill Expences with all the Expence contained in the database
// Ordered by most recent to eldest
func (ex Expences) GetAll() {
	sqlExpences := `
	SELECT * FROM expences
	ORDER BY name
	`
	ex.get(sqlExpences)
}

// GetRecurrent fill Expences with all the recurring Expence
func (ex Expences) GetRecurrent() {
	sqlExpences := `
	SELECT * FROM expences
	WHERE recurrency != 0
	ORDER BY name
	`
	ex.get(sqlExpences)
}

// GetNonRecurrent fill Expences with all the non recurring Expence
func (ex Expences) GetNonRecurrent() {
	sqlExpences := `
	SELECT * FROM expences
	WHERE recurrency = 0
	ORDER BY name
	`
	ex.get(sqlExpences)
}

// Len return the number of Expence in Expences
func (ex Expences) Len() int {
	return len(ex.expences)
}

// GetElement returns the Expence at the specified position
func (ex Expences) GetElement(i int) Found {
	// TODO: check that: -1 < i < Savings.Len()
	return ex.expences[i]
}

// Add adds an Expence to the database
func (ex Expences) Add(a Amount) {
	// TODO: check that the amount is not empty
	a.addAmount(ExpencesNames["table"], ExpencesNames["name"])
}

// Delete an Expence from the db and the struct
func (ex Expences) Delete(e Found) {
	for i := 0; i < ex.Len(); i++ {
		if ex.expences[i].Equals(e) {
			// Removing the Expence from the databse
			ex.expences[i].delete()
			// Removing the Expence from the struct
			ex.expences = append(ex.expences[:i], ex.expences[i+1:]...)
		}
	}
	// TODO: return error when the Expence is not present in the Expences
}

// Sum return the sum of all the Expence in Expences
func (ex Expences) Sum() float32 {
	var t float32
	if ex.Len() > 0 {
		for _, e := range ex.expences {
			t += (e.GetAmount()).sum
		}
	}
	return t
}

// String returns a string containing all the string of all the Expence contained in the struct
func (ex Expences) String() string {
	if ex.Len() == 0 {
		return " there are no expences "
	}
	var buffer bytes.Buffer
	for _, e := range ex.expences {
		buffer.WriteString(e.String())
		buffer.WriteString("\n")
	}
	return buffer.String()
}
