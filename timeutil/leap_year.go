package timeutil

import "math"

// IsLeapYear returns true only if the given year contains a leap day,
// meaning that the year is a leap year.
func IsLeapYear(givenYear int) bool {
	if givenYear%400 == 0 {
		return true
	} else if givenYear%100 == 0 {
		return false
	} else if givenYear%4 == 0 {
		return true
	}

	return false
}

// PrevLeapYear returns the previous leap year before a given year.
// It also returns a boolean value that is set to false if no
// such year can be found.
func PrevLeapYear(givenYear int) (foundYear int, found bool) {
	for foundYear = givenYear; foundYear > math.MinInt; {
		foundYear--

		if found = IsLeapYear(foundYear); found {
			return
		}
	}

	return
}

// NextLeapYear returns the next leap year after a given year.
// It also returns a boolean value that is set to false if no
// such year can be found.
func NextLeapYear(givenYear int) (foundYear int, found bool) {
	for foundYear = givenYear; foundYear < math.MaxInt; {
		foundYear++

		if found = IsLeapYear(foundYear); found {
			return
		}
	}

	return
}
