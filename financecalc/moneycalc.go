package financecalc

import (
	"time"

	"github.com/Mot93/Personal-Finance/financedatabase"
)

// StartDate return the time.Time to start counting
func StartDate() time.Time {
	loc := time.FixedZone("CET", 0)
	return time.Date(2018, time.July, 1, 0, 0, 0, 0, loc)
}

// SumNonRecurring TODO:
func SumNonRecurring(fo financedatabase.Founds) float32 {
	fo.GetNonRecurrent()
	if fo.Len() > 0 {
		return fo.Sum()
	}
	return 0.0
}

// SingleRecurrency return the total sum of a recurrent saving/expence
// TODO:
func SingleRecurrency(fo financedatabase.Found) float32 {
	return 2
}

// SumRecurring returns the summ of all recurrent savings/expences
func SumRecurring(fo financedatabase.Founds) float32 {
	var t float32 = 0.0
	fo.GetRecurrent()
	for i := 0; i < fo.Len(); i++ {
		t += SingleRecurrency(fo.GetElement(i))
	}
	return t
}
