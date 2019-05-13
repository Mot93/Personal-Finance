package financecalc

import (
	"time"

	"github.com/Mot93/Personal-Finance/financedatabase"
)

var loc = time.FixedZone("CET", 0)
var start = time.Date(2018, time.July, 1, 0, 0, 0, 0, loc)

// GetDate return the time.Time to start counting
func GetDate(f financedatabase.Found) time.Time {
	// TODO: as of right now the date is handled manually
	// it should be retrived from the database
	a := f.GetAmount()
	y, m, d := a.GetStart()
	return time.Date(y, time.Month(m), d, 0, 0, 0, 0, loc)
}

// SumNonRecurring return the sum of all the non recurring Expence/Saving of a collection
func SumNonRecurring(fo financedatabase.Founds) float32 {
	fo.GetNonRecurrent()
	if fo.Len() > 0 {
		return fo.Sum()
	}
	return 0.0
}

// SingleRecurrency return the sum of money a recurrent Expence/Saving has accumulated from it's start to:
// 		1) Now if the end date is not set, or is set for after this instant in time
// 		2) The end date specified if it's in the past
// Since the start date is option, always check if there is any start date
func SingleRecurrency(fo financedatabase.Found) float32 {
	// TODO: implements the method
	return 2
}

// SumRecurring returns the summ of all recurrent savings/expences
func SumRecurring(fo financedatabase.Founds) float32 {
	var t float32
	fo.GetRecurrent()
	for i := 0; i < fo.Len(); i++ {
		t += SingleRecurrency(fo.GetElement(i))
	}
	return t
}
