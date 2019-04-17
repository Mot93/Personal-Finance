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
	Get(string)
	// GetAll loads all the Expences/Savings in the DB
	GetAll()
	GetNonRecurrent()
	//GetRecurrent()
	Len() int
	ReturnElement(int) Found
	// Add an Expences/Savings and then call GetAll
	Add(Amount)
	Delete(Found)
	Sum() float32
	String() string
}
