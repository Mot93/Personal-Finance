package financetest

import (
	"testing"
	"time"

	"github.com/Mot93/personalfinance/financecalc"
)

func TestMonthSince(t *testing.T) {
	loc := time.FixedZone("CET", 0)
	start := time.Date(2000, time.June, 15, 0, 0, 0, 0, loc)
	tables := []struct {
		start    time.Time
		end      time.Time
		expected int
	}{
		{ // Same days
			start,
			start,
			1,
		},
		{ // Same year & month, succesive day
			start,
			time.Date(2000, time.July, 20, 0, 0, 0, 0, loc),
			1,
		},
		{ // Same year, month after, day before renew of the month
			start,
			time.Date(2000, time.July, 5, 0, 0, 0, 0, loc),
			1,
		},
		{ // Same year, month after, day after renew of the month
			start,
			time.Date(2000, time.July, 20, 0, 0, 0, 0, loc),
			1,
		},
		{ // Same year, 2 month after, day before renew of the month
			start,
			time.Date(2000, time.July, 10, 0, 0, 0, 0, loc),
			1,
		},
		{ // Few year later
			start,
			time.Date(2002, time.February, 1, 0, 0, 0, 0, loc),
			12 + 6 + 1,
		},
		{ // Year later
			start,
			time.Date(2001, time.March, 1, 0, 0, 0, 0, loc),
			6 + 2,
		},
		{ // The second date is set before the first one, I expect an error and nmonth to be 0
			start,
			time.Date(1999, time.February, 1, 0, 0, 0, 0, loc),
			0,
		},
		{ // 31 of febraruarhy becomes the 3 of march
			start,
			time.Date(2017, time.February, 31, 0, 0, 0, 0, loc),
			15*12 + 2 + 6,
		},
	}
	for _, table := range tables {
		out, err := financecalc.MonthSince(table.start, table.end)
		if err != nil { // Checking for erorrs
			if out != 0 { // when an error occours, the function should returns 0
				t.Errorf("nmonth is not 0, error: %v", err)
			}
		} else if out != table.expected {
			t.Errorf("Failed:\nstart %v\nend %v\n result %v ecspected %v", table.start, table.end, out, table.expected)
		}
	}
} // TestMonthSince
