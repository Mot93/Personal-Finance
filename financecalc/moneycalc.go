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

// NonRecurring TODO:
func NonRecurring(fo financedatabase.Founds) float32 {
	fo.GetNonRecurrent()
	if fo.Len() > 0 {
		return fo.Sum()
	}
	return 0.0
}

// Recurring TODO:
/*func Recurring(fo financedatabase.Founds) float32 {

}*/
