package financecalc

import (
	"errors"
	"time"
)

// Total returns the total sum of money that had to be saved
func Total() (int, error) {
	result, err := MonthSince(StartDate(), time.Now())
	//fmt.Printf("Month = %v\n", result)
	result *= MontlyCost()
	result -= PayedBills()
	result += FrozenFound()
	result += Savings()
	return result, err
}

// MonthSince returns an integer that rappresent the number of month the money had to be saved
// Returns 0 when an error occours
// The second (end) Time has to be set after the first one (begin)
func MonthSince(begin, end time.Time) (nmonth int, err error) {
	beginY, beginM, beginD := begin.Date()
	endY, endM, endD := end.Date()
	if !end.After(begin) { // Checking if the the second date comes after the first one
		if end.Equal(begin) {
			nmonth = 1
		} else {
			err = errors.New("The first date is set after the second one")
		}
	} else if beginY == endY { // Supposing the year is the same
		// Same month different days
		// the minimum is always 1
		if beginM == endM || beginM == endM+1 {
			nmonth = 1
		} else {
			nmonth = int(endM) - int(beginM)
		}
	} else { // different years
		nmonth = (12 - int(beginM)) + int(endM)
		if endY-beginY > 1 {
			nmonth += (endY - beginY - 1) * 12
		}
	}
	// A new month start when the the day start day is reached
	// if nmonth < 1 than I shouldn't modify it because:
	//		it's an error (second date is wrong)
	// 		it has to be at least 1
	// EX: In the event of a start date on a 31 vs a month that doesn't have 31 days, the new month will be counted on the 1 of the following month
	if beginD > endD && nmonth > 1 {
		nmonth--
	}
	return nmonth, err
}

// StartDate return the time.Time to start counting
func StartDate() time.Time {
	loc := time.FixedZone("CET", 0)
	return time.Date(2018, time.July, 1, 0, 0, 0, 0, loc)
}

// MontlyCost return how much should have be saved each month
func MontlyCost() int {
	return 180
}

// PayedBills returs the ammount of money that has already been spent in bills and more
func PayedBills() int {
	payed := 563 // IMU
	payed += 217 // Tasse condominio inizio 2019
	payed += 25  // Libro telecomunicazioni
	return payed
}

// Savings return the ammount of saving
func Savings() int {
	result := 2000 // Savings from mother inheritance
	return result
}

// FrozenFound return the ammount of money that musn't be spent
// For example the down payment
func FrozenFound() int {
	result := 1440 // Down payment
	return result
}
