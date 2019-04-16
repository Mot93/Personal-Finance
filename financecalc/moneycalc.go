package financecalc

import (
	"time"
)

// Total returns the total sum of money that had to be saved
func Total() float32 {
	result := Expences(StartDate(), time.Now())
	return result
}

// StartDate return the time.Time to start counting
func StartDate() time.Time {
	loc := time.FixedZone("CET", 0)
	return time.Date(2018, time.July, 1, 0, 0, 0, 0, loc)
}

// Return the total of money that has to be spent in expences from 
func Expences(start, end time.Time) float32 {
	financedatabase.
}
