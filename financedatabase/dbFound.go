package financedatabase

// Found define the intefaces for both Expence and Saving
type Found interface {
	Equals(Found) bool
	EqualValue(Amount) bool
	store()
	Update(Amount)
	delete()
	ReturnAmount() Amount
	String() string
}

// Founds interfaces to manage group of Expences and Savings
type Founds interface {
	// GetAll loads all the Expences/Savings in the DB
	GetAll()
	Len() int
	ReturnElement(int) Found
	// Add an Expences/Savings and then call GetAll
	Add(Amount)
	Delete(Found)
	String() string
}
