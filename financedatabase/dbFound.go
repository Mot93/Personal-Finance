package financedatabase

/* Since Expence and Saving are the same data with just different name, Found and Founds are interfaces created to make a single tool (Amount) to reduce the code write
 * */

// Found define the intefaces for both Expence and Saving
type Found interface {
	// Equals two Expence/Saving by id and value
	Equals(Found) bool
	// Equals the Amount contained in two Expence/Saving
	EqualValue(Amount) bool
	// store a Found in the database
	store()
	// Update the values of the Amount of the Found and the database
	Update(Amount)
	// delete an Expence/Saving from the database
	delete()
	// GetAmount returns the Amount of the Expence/Saving
	GetAmount() Amount
	// String returns a sctring containing id and values
	String() string
}

// Founds interfaces to manage group of Expences and Savings
type Founds interface {
	// Get returns
	get(string)
	// GetAll loads all the Expence/Saving in the DB
	GetAll()
	// GetRecurrent loads all the Expence/Saving in the DB, that are recurring
	GetRecurrent()
	// GetNonRecurrent loads all the Expence/Saving in the DB, that are non recurring
	GetNonRecurrent()
	// Len returns how many Found are in Founds
	Len() int
	// GetElement return a single Found a the specified position
	GetElement(int) Found
	// Add an Expence/Saving to the databade but not to the slice
	// The slice is not update since it's uncertain if the Expence/Saving belongs to that collection of Saving/Expence
	Add(Amount)
	// Delete and Expence/Saving from the database and the struct
	Delete(Found)
	// Sum all the Expence/saving values contained in the struct
	Sum() float32
	// String returns a string with all the strings of each Expence/Saving contained in the struct
	String() string
}
