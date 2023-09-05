package timeutil

import (
	"errors"
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	Timezone string = "Asia/Kathmandu"
)

// Reference date for conversion is 2000/01/01 BS and 1943/4/14 AD
var npInitialYear int16 = 2000
var referenceEnDate = [3]int16{1943, 4, 14}

type NepaliMonthData struct {
	monthData [12]int8
	yearDays  int16
}

var enMonths = [12]int8{31, 28, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31}
var enLeapMonths = [12]int8{31, 29, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31}

var npMonthData = [100]NepaliMonthData{
	{[12]int8{30, 32, 31, 32, 31, 30, 30, 30, 29, 30, 29, 31}, 365}, // 2000 BS - 1943/1944 AD
	{[12]int8{31, 31, 32, 31, 31, 31, 30, 29, 30, 29, 30, 30}, 365}, // 2001 BS
	{[12]int8{31, 31, 32, 32, 31, 30, 30, 29, 30, 29, 30, 30}, 365},
	{[12]int8{31, 32, 31, 32, 31, 30, 30, 30, 29, 29, 30, 31}, 366},
	{[12]int8{30, 32, 31, 32, 31, 30, 30, 30, 29, 30, 29, 31}, 365},
	{[12]int8{31, 31, 32, 31, 31, 31, 30, 29, 30, 29, 30, 30}, 365},
	{[12]int8{31, 31, 32, 32, 31, 30, 30, 29, 30, 29, 30, 30}, 365},
	{[12]int8{31, 32, 31, 32, 31, 30, 30, 30, 29, 29, 30, 31}, 366},
	{[12]int8{31, 31, 31, 32, 31, 31, 29, 30, 30, 29, 29, 31}, 365},
	{[12]int8{31, 31, 32, 31, 31, 31, 30, 29, 30, 29, 30, 30}, 365},
	{[12]int8{31, 31, 32, 32, 31, 30, 30, 29, 30, 29, 30, 30}, 365},
	{[12]int8{31, 32, 31, 32, 31, 30, 30, 30, 29, 29, 30, 31}, 366},
	{[12]int8{31, 31, 31, 32, 31, 31, 29, 30, 30, 29, 30, 30}, 365},
	{[12]int8{31, 31, 32, 31, 31, 31, 30, 29, 30, 29, 30, 30}, 365},
	{[12]int8{31, 31, 32, 32, 31, 30, 30, 29, 30, 29, 30, 30}, 365},
	{[12]int8{31, 32, 31, 32, 31, 30, 30, 30, 29, 29, 30, 31}, 366},
	{[12]int8{31, 31, 31, 32, 31, 31, 29, 30, 30, 29, 30, 30}, 365},
	{[12]int8{31, 31, 32, 31, 31, 31, 30, 29, 30, 29, 30, 30}, 365},
	{[12]int8{31, 32, 31, 32, 31, 30, 30, 29, 30, 29, 30, 30}, 365},
	{[12]int8{31, 32, 31, 32, 31, 30, 30, 30, 29, 30, 29, 31}, 366},
	{[12]int8{31, 31, 31, 32, 31, 31, 30, 29, 30, 29, 30, 30}, 365},
	{[12]int8{31, 31, 32, 31, 31, 31, 30, 29, 30, 29, 30, 30}, 365},
	{[12]int8{31, 32, 31, 32, 31, 30, 30, 30, 29, 29, 30, 30}, 365},
	{[12]int8{31, 32, 31, 32, 31, 30, 30, 30, 29, 30, 29, 31}, 366},
	{[12]int8{31, 31, 31, 32, 31, 31, 30, 29, 30, 29, 30, 30}, 365},
	{[12]int8{31, 31, 32, 31, 31, 31, 30, 29, 30, 29, 30, 30}, 365},
	{[12]int8{31, 32, 31, 32, 31, 30, 30, 30, 29, 29, 30, 31}, 366},
	{[12]int8{30, 32, 31, 32, 31, 30, 30, 30, 29, 30, 29, 31}, 365},
	{[12]int8{31, 31, 32, 31, 31, 31, 30, 29, 30, 29, 30, 30}, 365},
	{[12]int8{31, 31, 32, 31, 32, 30, 30, 29, 30, 29, 30, 30}, 365},
	{[12]int8{31, 32, 31, 32, 31, 30, 30, 30, 29, 29, 30, 31}, 366},
	{[12]int8{30, 32, 31, 32, 31, 30, 30, 30, 29, 30, 29, 31}, 365},
	{[12]int8{31, 31, 32, 31, 31, 31, 30, 29, 30, 29, 30, 30}, 365},
	{[12]int8{31, 31, 32, 32, 31, 30, 30, 29, 30, 29, 30, 30}, 365},
	{[12]int8{31, 32, 31, 32, 31, 30, 30, 30, 29, 29, 30, 31}, 366},
	{[12]int8{30, 32, 31, 32, 31, 31, 29, 30, 30, 29, 29, 31}, 365},
	{[12]int8{31, 31, 32, 31, 31, 31, 30, 29, 30, 29, 30, 30}, 365},
	{[12]int8{31, 31, 32, 32, 31, 30, 30, 29, 30, 29, 30, 30}, 365},
	{[12]int8{31, 32, 31, 32, 31, 30, 30, 30, 29, 29, 30, 31}, 366},
	{[12]int8{31, 31, 31, 32, 31, 31, 29, 30, 30, 29, 30, 30}, 365},
	{[12]int8{31, 31, 32, 31, 31, 31, 30, 29, 30, 29, 30, 30}, 365},
	{[12]int8{31, 31, 32, 32, 31, 30, 30, 29, 30, 29, 30, 30}, 365},
	{[12]int8{31, 32, 31, 32, 31, 30, 30, 30, 29, 29, 30, 31}, 366},
	{[12]int8{31, 31, 31, 32, 31, 31, 29, 30, 30, 29, 30, 30}, 365},
	{[12]int8{31, 31, 32, 31, 31, 31, 30, 29, 30, 29, 30, 30}, 365},
	{[12]int8{31, 32, 31, 32, 31, 30, 30, 29, 30, 29, 30, 30}, 365},
	{[12]int8{31, 32, 31, 32, 31, 30, 30, 30, 29, 29, 30, 31}, 366},
	{[12]int8{31, 31, 31, 32, 31, 31, 30, 29, 30, 29, 30, 30}, 365},
	{[12]int8{31, 31, 32, 31, 31, 31, 30, 29, 30, 29, 30, 30}, 365},
	{[12]int8{31, 32, 31, 32, 31, 30, 30, 30, 29, 29, 30, 30}, 365},
	{[12]int8{31, 32, 31, 32, 31, 30, 30, 30, 29, 30, 29, 31}, 366},
	{[12]int8{31, 31, 31, 32, 31, 31, 30, 29, 30, 29, 30, 30}, 365},
	{[12]int8{31, 31, 32, 31, 31, 31, 30, 29, 30, 29, 30, 30}, 365},
	{[12]int8{31, 32, 31, 32, 31, 30, 30, 30, 29, 29, 30, 30}, 365},
	{[12]int8{31, 32, 31, 32, 31, 30, 30, 30, 29, 30, 29, 31}, 366},
	{[12]int8{31, 31, 32, 31, 31, 31, 30, 29, 30, 29, 30, 30}, 365},
	{[12]int8{31, 31, 32, 31, 32, 30, 30, 29, 30, 29, 30, 30}, 365},
	{[12]int8{31, 32, 31, 32, 31, 30, 30, 30, 29, 29, 30, 31}, 366},
	{[12]int8{30, 32, 31, 32, 31, 30, 30, 30, 29, 30, 29, 31}, 365},
	{[12]int8{31, 31, 32, 31, 31, 31, 30, 29, 30, 29, 30, 30}, 365},
	{[12]int8{31, 31, 32, 32, 31, 30, 30, 29, 30, 29, 30, 30}, 365},
	{[12]int8{31, 32, 31, 32, 31, 30, 30, 30, 29, 29, 30, 31}, 366},
	{[12]int8{30, 32, 31, 32, 31, 31, 29, 30, 29, 30, 29, 31}, 365},
	{[12]int8{31, 31, 32, 31, 31, 31, 30, 29, 30, 29, 30, 30}, 365},
	{[12]int8{31, 31, 32, 32, 31, 30, 30, 29, 30, 29, 30, 30}, 365},
	{[12]int8{31, 32, 31, 32, 31, 30, 30, 30, 29, 29, 30, 31}, 366},
	{[12]int8{31, 31, 31, 32, 31, 31, 29, 30, 30, 29, 29, 31}, 365},
	{[12]int8{31, 31, 32, 31, 31, 31, 30, 29, 30, 29, 30, 30}, 365},
	{[12]int8{31, 31, 32, 32, 31, 30, 30, 29, 30, 29, 30, 30}, 365},
	{[12]int8{31, 32, 31, 32, 31, 30, 30, 30, 29, 29, 30, 31}, 366},
	{[12]int8{31, 31, 31, 32, 31, 31, 29, 30, 30, 29, 30, 30}, 365},
	{[12]int8{31, 31, 32, 31, 31, 31, 30, 29, 30, 29, 30, 30}, 365},
	{[12]int8{31, 32, 31, 32, 31, 30, 30, 29, 30, 29, 30, 30}, 365},
	{[12]int8{31, 32, 31, 32, 31, 30, 30, 30, 29, 29, 30, 31}, 366},
	{[12]int8{31, 31, 31, 32, 31, 31, 30, 29, 30, 29, 30, 30}, 365},
	{[12]int8{31, 31, 32, 31, 31, 31, 30, 29, 30, 29, 30, 30}, 365},
	{[12]int8{31, 32, 31, 32, 31, 30, 30, 30, 29, 29, 30, 30}, 365},
	{[12]int8{31, 32, 31, 32, 31, 30, 30, 30, 29, 30, 29, 31}, 366},
	{[12]int8{31, 31, 31, 32, 31, 31, 30, 29, 30, 29, 30, 30}, 365},
	{[12]int8{31, 31, 32, 31, 31, 31, 30, 29, 30, 29, 30, 30}, 365},
	{[12]int8{31, 32, 31, 32, 31, 30, 30, 30, 29, 29, 30, 30}, 365},
	{[12]int8{31, 31, 32, 32, 31, 30, 30, 30, 29, 30, 30, 30}, 366},
	{[12]int8{30, 32, 31, 32, 31, 30, 30, 30, 29, 30, 30, 30}, 365},
	{[12]int8{31, 31, 32, 31, 31, 30, 30, 30, 29, 30, 30, 30}, 365},
	{[12]int8{31, 31, 32, 31, 31, 30, 30, 30, 29, 30, 30, 30}, 365},
	{[12]int8{31, 32, 31, 32, 30, 31, 30, 30, 29, 30, 30, 30}, 366},
	{[12]int8{30, 32, 31, 32, 31, 30, 30, 30, 29, 30, 30, 30}, 365},
	{[12]int8{31, 31, 32, 31, 31, 31, 30, 30, 29, 30, 30, 30}, 366},
	{[12]int8{30, 31, 32, 32, 30, 31, 30, 30, 29, 30, 30, 30}, 365},
	{[12]int8{30, 32, 31, 32, 31, 30, 30, 30, 29, 30, 30, 30}, 365},
	{[12]int8{30, 32, 31, 32, 31, 30, 30, 30, 29, 30, 30, 30}, 365},
	{[12]int8{31, 31, 32, 31, 31, 31, 30, 30, 29, 30, 30, 30}, 366},
	{[12]int8{30, 31, 32, 32, 31, 30, 30, 30, 29, 30, 30, 30}, 365},
	{[12]int8{30, 32, 31, 32, 31, 30, 30, 30, 29, 30, 30, 30}, 365},
	{[12]int8{31, 31, 32, 31, 31, 30, 30, 30, 29, 30, 30, 30}, 365},
	{[12]int8{31, 31, 32, 31, 31, 31, 30, 29, 30, 30, 30, 30}, 366},
	{[12]int8{30, 31, 32, 32, 31, 30, 30, 29, 30, 29, 30, 30}, 364},
	{[12]int8{31, 32, 31, 32, 31, 30, 30, 30, 29, 30, 30, 30}, 366},
	{[12]int8{31, 31, 32, 31, 31, 31, 29, 30, 29, 30, 29, 31}, 365},
	{[12]int8{31, 31, 32, 31, 31, 31, 30, 29, 29, 30, 30, 30}, 365}, // 2099 BS - 2042 AD
}

func enMinYear() int {
	return int(referenceEnDate[0]) + 1
}

func enMaxYear() int {
	return int(referenceEnDate[0]) + len(npMonthData) - 1
}

func npMinYear() int {
	return int(npInitialYear)
}

func npMaxYear() int {
	return int(npInitialYear) + len(npMonthData) - 1
}

/* Checks if the english year is leap year or not */
func isLeapYear(year int) bool {
	if year%4 == 0 {
		if year%100 == 0 {
			return (year%400 == 0)
		}
		return true
	}
	return false
}

func getEnMonths(year int) *[12]int8 {
	if isLeapYear(year) {
		return &enLeapMonths
	}
	return &enMonths
}

/*
Returns diff from the reference with its absolute reference.
Used in function `NepaliToEnglish`.

Eg. ref: 1943/4/14 - 1943/01/01
*/
func getDiffFromEnAbsoluteReference() int {
	var diff int = 0

	// adding sum of month of year till the reference month
	months := getEnMonths(int(referenceEnDate[0]))
	for i := 0; i < int(referenceEnDate[1])-1; i++ {
		diff += int(months[i])
	}

	return diff + int(referenceEnDate[2]) - 1 // added day too
}

// ENGLISH DATE CONVERSION

// checks if english date in within range 1944 - 2042
func checkEnglishDate(year int, month int, day int) bool {
	if year < enMinYear() || year > enMaxYear() {
		return false
	}
	if month < 1 || month > 12 {
		return false
	}
	if day < 1 || day > 31 {
		return false
	}
	return true
}

// counts and returns total days from the date 0000-01-01
func getTotalDaysFromEnglishDate(year int, month int, day int) int {
	var totalDays int = year*365 + day
	for i := 0; i < month-1; i++ {
		totalDays = totalDays + int(enMonths[i])
	}

	// adding leap days (ie. leap year count)
	if month <= 2 { // checking February month (where leap exists)
		year -= 1
	}
	totalDays += int(year/4) - int(year/100) + int(year/400)

	return totalDays
}

// NEPALI DATE CONVERSION

// checks if nepali date is in range
func checkNepaliDate(year int, month int, day int) bool {
	if year < npMinYear() || year > npMaxYear() {
		return false
	}
	if month < 1 || month > 12 {
		return false
	}

	if day < 1 || day > int(npMonthData[year-int(npInitialYear)].monthData[month-1]) {
		return false
	}
	return true
}

// counts and returns total days from the nepali initial date
func getTotalDaysFromNepaliDate(year int, month int, day int) int {
	var totalDays int = day - 1

	// adding days of months of initial year
	var yearIndex int = year - int(npInitialYear)
	for i := 0; i < month-1; i++ {
		totalDays = totalDays + int(npMonthData[yearIndex].monthData[i])
	}

	// adding days of year
	for i := 0; i < yearIndex; i++ {
		totalDays = totalDays + int(npMonthData[i].yearDays)
	}
	return totalDays
}

// Public methods

// Converts english date to nepali.
// Accepts the input parameters year, month, day.
// Returns dates in array and error.
func EnglishToNepali(year int, month int, day int) (*[3]int, error) {
	// VALIDATION
	// checking if date is in range
	if !checkEnglishDate(year, month, day) {
		return nil, errors.New("date is out of range")
	}

	// REFERENCE
	npYear, npMonth, npDay := int(npInitialYear), 1, 1

	// DIFFERENCE
	// calculating days count from the reference date
	var difference int = int(
		math.Abs(float64(
			getTotalDaysFromEnglishDate(
				year, month, day,
			) - getTotalDaysFromEnglishDate(
				int(referenceEnDate[0]), int(referenceEnDate[1]), int(referenceEnDate[2]),
			),
		)),
	)

	// YEAR
	// Incrementing year until the difference remains less than 365
	var yearDataIndex int = 0
	for difference >= int(npMonthData[yearDataIndex].yearDays) {
		difference -= int(npMonthData[yearDataIndex].yearDays)
		npYear += 1
		yearDataIndex += 1
	}

	// MONTH
	// Incrementing month until the difference remains less than next nepali month days (mostly 31)
	var i int = 0
	for difference >= int(npMonthData[yearDataIndex].monthData[i]) {
		difference -= int(npMonthData[yearDataIndex].monthData[i])
		npMonth += 1
		i += 1
	}

	// DAY
	// Remaining difference is the day
	npDay += difference

	return &[3]int{npYear, npMonth, npDay}, nil
}

// Converts nepali date to english.
// Accepts the input parameters year, month, day.
// Returns dates in array and error.
func NepaliToEnglish(year int, month int, day int) (*[3]int, error) {
	// VALIDATION
	// checking if date is in range
	if !checkNepaliDate(year, month, day) {
		return nil, errors.New("date is out of range")
	}

	// REFERENCE
	// For absolute reference, moving date to Jan 1
	// Eg. ref: 1943/4/14 => 1943/01/01
	enYear, enMonth, enDay := int(referenceEnDate[0]), 1, 1
	// calculating difference from the adjusted reference (eg. 1943/4/14 - 1943/01/01)
	referenceDiff := getDiffFromEnAbsoluteReference()

	// DIFFERENCE
	// calculating days count from the reference date
	var difference int = getTotalDaysFromNepaliDate(year, month, day) + referenceDiff

	// YEAR
	// Incrementing year until the difference remains less than 365 (or 365)
	for (difference >= 366 && isLeapYear(enYear)) || (difference >= 365 && !(isLeapYear(enYear))) {
		if isLeapYear(enYear) {
			difference -= 366
		} else {
			difference -= 365
		}
		enYear += 1
	}

	// MONTH
	// Incrementing month until the difference remains less than next english month (mostly 31)
	monthDays := getEnMonths(enYear)

	var i int = 0
	for difference >= int(monthDays[i]) {
		difference -= int(monthDays[i])
		enMonth += 1
		i += 1
	}

	// DAY
	// Remaining difference is the day
	enDay += difference

	return &[3]int{enYear, enMonth, enDay}, nil
}

func NewFormatter(nepaliTime *NepaliTime) *NepaliFormatter {
	return &NepaliFormatter{nepaliTime: nepaliTime}
}

type NepaliFormatter struct {
	nepaliTime *NepaliTime
}

func (obj *NepaliFormatter) Format(format string) (string, error) {
	index, num, timeStr := 0, len(format), ""

	for index < num {
		char := string(format[index])
		index++

		if char == "%" && index < num {
			char = string(format[index])

			if char == "%" {
				timeStr += char
			} else if char == "-" {
				specialChar := char

				if (index + 1) < num {
					index++
					char = string(format[index])
					res, err := obj.getFormatString(specialChar + char)
					if err != nil {
						return "", errors.New("error while formatting NepaliTime with given format")
					}
					timeStr += res
				}
			} else {
				res, err := obj.getFormatString(char)
				if err != nil {
					return "", errors.New("error while formatting NepaliTime with given format")
				}
				timeStr += res
			}
			index++
		} else {
			timeStr += char
		}
	}

	return timeStr, nil
}

// utility method that operates based on the type of directive
func (obj *NepaliFormatter) getFormatString(directive string) (string, error) {
	switch directive {
	case "d":
		return obj.day_(), nil
	case "-d":
		return obj.dayNonzero(), nil
	case "m":
		return obj.monthNumber(), nil
	case "-m":
		return obj.monthNumberNonzero(), nil
	case "y":
		return obj.yearHalf(), nil
	case "Y":
		return obj.yearFull(), nil
	case "H":
		return obj.hour24(), nil
	case "-H":
		return obj.hour24Nonzero(), nil
	case "I":
		return obj.hour12(), nil
	case "-I":
		return obj.hour12Nonzero(), nil
	case "p":
		return obj.ampm(), nil
	case "M":
		return obj.minute(), nil
	case "-M":
		return obj.minuteNonzero(), nil
	case "S":
		return obj.second(), nil
	case "-S":
		return obj.secondNonzero(), nil
	case "f":
		return obj.nanosecond_(), nil
	case "-f":
		return obj.nanosecondNonZero(), nil
	default:
		return "", errors.New("error while getting format string for passed directive")
	}
}

// %d
func (obj *NepaliFormatter) day_() string {
	day := strconv.Itoa(obj.nepaliTime.day)

	if len(day) < 2 {
		day = "0" + day
	}

	return day
}

// -d
func (obj *NepaliFormatter) dayNonzero() string {
	day := strconv.Itoa(obj.nepaliTime.day)

	return day
}

// %m
func (obj *NepaliFormatter) monthNumber() string {
	month := strconv.Itoa(obj.nepaliTime.month)

	if len(month) < 2 {
		month = "0" + month
	}

	return month
}

// %-m
func (obj *NepaliFormatter) monthNumberNonzero() string {
	month := strconv.Itoa(obj.nepaliTime.month)

	return month
}

// %y
func (obj *NepaliFormatter) yearHalf() string {
	year := strconv.Itoa(obj.nepaliTime.year)

	return year[2:]
}

// %Y
func (obj *NepaliFormatter) yearFull() string {
	return strconv.Itoa(obj.nepaliTime.year)
}

// %H
func (obj *NepaliFormatter) hour24() string {
	hour := strconv.Itoa(obj.nepaliTime.Hour())

	if len(hour) < 2 {
		hour = "0" + hour
	}

	return hour
}

// %-H
func (obj *NepaliFormatter) hour24Nonzero() string {
	return strconv.Itoa(obj.nepaliTime.Hour())
}

// %I
func (obj *NepaliFormatter) hour12() string {
	hour := obj.nepaliTime.Hour()

	if hour > 12 {
		hour -= 12
	}
	if hour == 0 {
		hour = 12
	}

	hourStr := strconv.Itoa(hour)
	if len(hourStr) < 2 {
		hourStr = "0" + hourStr
	}

	return hourStr
}

// %-I
func (obj *NepaliFormatter) hour12Nonzero() string {
	hour := obj.nepaliTime.Hour()

	if hour > 12 {
		hour -= 12
	}
	if hour == 0 {
		hour = 12
	}

	return strconv.Itoa(hour)
}

// %p
func (obj *NepaliFormatter) ampm() string {
	ampm := "AM"

	if obj.nepaliTime.Hour() > 12 {
		ampm = "PM"
	}

	return ampm
}

// %M
func (obj *NepaliFormatter) minute() string {
	minute := strconv.Itoa(obj.nepaliTime.Minute())

	if len(minute) < 2 {
		minute = "0" + minute
	}

	return minute
}

// %-M
func (obj *NepaliFormatter) minuteNonzero() string {
	return strconv.Itoa(obj.nepaliTime.Minute())
}

// %s
func (obj *NepaliFormatter) second() string {
	second := strconv.Itoa(obj.nepaliTime.Second())

	if len(second) < 2 {
		second = "0" + second
	}

	return second
}

// %-s
func (obj *NepaliFormatter) secondNonzero() string {
	return strconv.Itoa(obj.nepaliTime.Second())
}

// %f
func (obj *NepaliFormatter) nanosecond_() string {
	nsec := strconv.Itoa(obj.nepaliTime.Nanosecond())

	if len(nsec) < 6 {
		nsec = strings.Repeat("0", 6-len(nsec)) + nsec
	}

	return nsec
}

// %-f
func (obj *NepaliFormatter) nanosecondNonZero() string {
	return strconv.Itoa(obj.nepaliTime.Nanosecond())
}

// NOTE:
// Only Date() and FromEnglishTime() from utils are supposed to create NepaliTime object

type NepaliTime struct {
	// note that these members represent nepali values
	year        int
	month       int
	day         int
	englishTime *time.Time
}

// String returns a string representing the duration in the form "2079-10-06 01:00:05"
func (obj *NepaliTime) String() string {
	h, m, s := obj.Clock()
	return fmt.Sprintf(
		"%d-%s-%s %s:%s:%s",
		obj.year,
		twoDigitNumber(obj.month),
		twoDigitNumber(obj.day),
		twoDigitNumber(h),
		twoDigitNumber(m),
		twoDigitNumber(s),
	)
}

// Get's the corresponding english time
func (obj *NepaliTime) GetEnglishTime() time.Time {
	return *obj.englishTime
}

// Date returns the year, month, and day
func (obj *NepaliTime) Date() (year, month, day int) {
	return obj.year, obj.month, obj.day
}

// Year returns the year in which nepalitime occurs.
func (obj *NepaliTime) Year() int {
	return obj.year
}

// Month returns the month of the year.
func (obj *NepaliTime) Month() int {
	return obj.month
}

// Day returns the day of the month.
func (obj *NepaliTime) Day() int {
	return obj.day
}

// Weekday returns the day of the week.
// Sunday = 0,
// Monday = 1,
// Saturday = 6
func (obj *NepaliTime) Weekday() time.Weekday {
	return obj.englishTime.Weekday()
}

// Clock returns the hour, minute, and second of the day.
func (obj *NepaliTime) Clock() (hour, min, sec int) {
	return obj.englishTime.Clock()
}

// Hour returns the hour within the day, in the range [0, 23].
func (obj *NepaliTime) Hour() int {
	return obj.englishTime.Hour()
}

// Minute returns the minute offset within the hour, in the range [0, 59].
func (obj *NepaliTime) Minute() int {
	return obj.englishTime.Minute()
}

// Second returns the second offset within the minute, in the range [0, 59].
func (obj *NepaliTime) Second() int {
	return obj.englishTime.Second()
}

// Nanosecond returns the nanosecond offset within the second,
// in the range [0, 999999999].
func (obj *NepaliTime) Nanosecond() int {
	return obj.englishTime.Nanosecond()
}

// formats the nepalitime object into the passed format
func (obj *NepaliTime) Format(format string) (string, error) {
	formatter := NewFormatter(obj)

	formattedNepaliTime, err := formatter.Format(format)
	if err != nil {
		return "", err
	}

	return formattedNepaliTime, nil
}

type nepaliTimeRegex struct {
	PatternMap map[string]string
}

// nepaliTimeRe constructor
func newNepaliTimeRegex() *nepaliTimeRegex {
	obj := new(nepaliTimeRegex)
	obj.PatternMap = map[string]string{
		"d":  `(?P<d>3[0-2]|[1-2]\d|0[1-9]|[1-9]| [1-9])`,
		"-d": `(?P<d>3[0-2]|[1-2]\d|0[1-9]|[1-9]| [1-9])`, // same as "d"
		"f":  `(?P<f>[0-9]{1,6})`,
		"H":  `(?P<H>2[0-3]|[0-1]\d|\d)`,
		"-H": `(?P<H>2[0-3]|[0-1]\d|\d)`,
		"I":  `(?P<I>1[0-2]|0[1-9]|[1-9])`,
		"-I": `(?P<I>1[0-2]|0[1-9]|[1-9])`,
		"G":  `(?P<G>\d\d\d\d)`,
		"j":  `(?P<j>36[0-6]|3[0-5]\d|[1-2]\d\d|0[1-9]\d|00[1-9]|[1-9]\d|0[1-9]|[1-9])`,
		"m":  `(?P<m>1[0-2]|0[1-9]|[1-9])`,
		"-m": `(?P<m>1[0-2]|0[1-9]|[1-9])`, // same as "m"
		"M":  `(?P<M>[0-5]\d|\d)`,
		"-M": `(?P<M>[0-5]\d|\d)`, // same as "M"
		"S":  `(?P<S>6[0-1]|[0-5]\d|\d)`,
		"-S": `(?P<S>6[0-1]|[0-5]\d|\d)`, // same as "S"
		"w":  `(?P<w>[0-6])`,

		"y": `(?P<y>\d\d)`,
		"Y": `(?P<Y>\d\d\d\d)`,
		"z": `(?P<z>[+-]\d\d:?[0-5]\d(:?[0-5]\d(\.\d{1,6})?)?|(?-i:Z))`,

		// "A": obj.__seqToRE(EnglishChar.days, "A"),
		// "a": obj.__seqToRE(EnglishChar.days_half, "a"),
		// "B": obj.__seqToRE(EnglishChar.months, "B"),
		// "b": obj.__seqToRE(EnglishChar.months, "b"),
		// "p": obj.__seqToRE(("AM", "PM",), "p"),
		// TODO: implement for the above commented directives
		"p": "(?P<p>AM|PM)",

		"%": "%",
	}

	return obj
}

// Handles conversion from format directives to regexes
func (obj *nepaliTimeRegex) pattern(format string) (string, error) {
	processedFormat := ""
	regexChars := regexp.MustCompile(`([\.^$*+?\(\){}\[\]|])`)
	format = regexChars.ReplaceAllString(format, `\$1`)
	whitespaceReplacement := regexp.MustCompile(`\s+`)
	format = whitespaceReplacement.ReplaceAllString(format, `\s+`)

	for {
		index := strings.Index(format, "%")
		// index = -1 means the sub string does not exist
		// in the string
		if index == -1 {
			break
		}

		directiveIndex := index + 1
		indexIncrement := 1

		if string(format[directiveIndex]) == "-" {
			indexIncrement = 2
		}

		directiveToCheck := string(format[directiveIndex : directiveIndex+indexIncrement])

		if val, ok := obj.PatternMap[directiveToCheck]; ok {
			processedFormat = fmt.Sprintf("%s%s%s", processedFormat, string(format[:directiveIndex-1]), val)
			format = string(format[directiveIndex+indexIncrement:])
		} else {
			return "", errors.New("no pattern matched")
		}
	}

	tester := fmt.Sprintf("^%s%s$", processedFormat, format)

	return tester, nil
}

// handles regex compilation for format string
func (obj *nepaliTimeRegex) compile(format string) (*regexp.Regexp, error) {
	processedFormat, err := obj.pattern(format)

	if err != nil {
		return nil, err
	}

	// (?i) is for ignoring the case
	reg, err := regexp.Compile("(?i)" + processedFormat)

	if err != nil {
		return nil, err
	}

	return reg, nil
}

var nepaliTimeReCache *nepaliTimeRegex

// ParseNP equivalent to time.Parse()
func ParseNP(datetimeStr string, format string) (*NepaliTime, error) {
	nepalitime, err := validate(datetimeStr, format)

	if err != nil {
		return nil, errors.New("datetime string did not match with given format")
	}

	return nepalitime, nil
}

func getNepaliTimeReObject() *nepaliTimeRegex {
	if nepaliTimeReCache == nil {
		nepaliTimeReCache = newNepaliTimeRegex()
	}

	return nepaliTimeReCache
}

// validates datetimeStr with the format
func validate(datetimeStr string, format string) (*NepaliTime, error) {
	// validate if parse result is not empty
	parsedResult, err := extract(datetimeStr, format)
	if err != nil {
		return nil, err
	} else {
		_, ok1 := parsedResult["Y"]
		_, ok2 := parsedResult["y"]

		if !ok1 && !ok2 {
			return nil, errors.New("unable to parse year")
		}
	}

	// validate the transformation
	transformedData, err := transform(parsedResult)
	if err != nil {
		return nil, err
	}

	nepaliDate, err := Date(
		transformedData["year"],
		transformedData["month"],
		transformedData["day"],
		transformedData["hour"],
		transformedData["minute"],
		transformedData["second"],
		transformedData["nanosecond"],
	)
	if err != nil {
		return nil, err
	}

	return nepaliDate, nil
}

// extracts year, month, day, hour, minute, etc from the given format
// eg.
// USAGE: extract("2078-01-12", "%Y-%m-%d")
// INPUT:
// datetime_str="2078-01-12"
// format="%Y-%m-%d"
// OUTPUT:
//
//	{
//		"Y": 2078,
//		"m": 1,
//		"d": 12,
//	}
func extract(datetimeStr string, format string) (map[string]string, error) {
	reCompiledFormat, err := getNepaliTimeReObject().compile(format)

	if err != nil {
		return nil, err
	}

	match := reCompiledFormat.FindStringSubmatch(datetimeStr)

	if len(match) < 1 {
		return nil, errors.New("no pattern matched")
	}

	result := make(map[string]string)

	for index, name := range reCompiledFormat.SubexpNames() {
		if index != 0 && name != "" {
			result[name] = match[index]
		}
	}

	return result, nil
}

// transforms different format data to uniform data
// eg.
// INPUT:
//
//	data = {
//	    "Y": 2078,
//	    "b": "Mangsir",
//	    "d": 12,
//	    ...
//	}
//
// OUTPUT:
//
//	{
//	    "year": 2078,
//	    "month": 8,
//	    "day": 12,
//	    ...
//	}
func transform(data map[string]string) (map[string]int, error) {
	var (
		year                           int
		month, day                     int = 1, 1
		hour, minute, second, fraction int = 0, 0, 0, 0
	)

	for key, val := range data {
		if key == "y" {
			intVal, err := strconv.Atoi(val)
			if err != nil {
				return nil, errors.New("invalid value in %y")
			}

			year = intVal
			year += 2000
		} else if key == "Y" {
			intVal, err := strconv.Atoi(val)
			if err != nil {
				return nil, errors.New("invalid value in %Y")
			}

			year = intVal
		} else if key == "m" {
			intVal, err := strconv.Atoi(val)
			if err != nil {
				return nil, errors.New("invalid value in %Y")
			}

			month = intVal
		} else if key == "d" {
			intVal, err := strconv.Atoi(val)
			if err != nil {
				return nil, errors.New("invalid value in %d")
			}

			day = intVal
		} else if key == "H" {
			intVal, err := strconv.Atoi(val)
			if err != nil {
				return nil, errors.New("invalid value in %H")
			}

			hour = intVal
		} else if key == "I" {
			intVal, err := strconv.Atoi(val)
			if err != nil {
				return nil, errors.New("invalid value in %I")
			}

			hour = intVal

			ampm, ok := data["p"]
			if !ok {
				ampm = ""
			}

			ampm = strings.ToLower(ampm)
			// if there is no AM/PM indicator, we'll treat it as
			if ampm == "" || ampm == "am" {
				if hour == 12 {
					hour = 0
				}
			} else if ampm == "pm" {
				if hour != 12 {
					hour += 12
				}
			}
		} else if key == "M" {
			intVal, err := strconv.Atoi(val)
			if err != nil {
				return nil, errors.New("invalid value in %M")
			}

			minute = intVal
		} else if key == "S" {
			intVal, err := strconv.Atoi(val)
			if err != nil {
				return nil, errors.New("invalid value in %S")
			}

			second = intVal
		} else if key == "f" {
			var err error

			s := val
			s += strings.Repeat("0", 6-len(val))

			fraction, err = strconv.Atoi(s)

			if err != nil {
				return nil, errors.New("error while getting nanoseconds data")
			}
		}
	}

	return map[string]int{
		"year":       year,
		"month":      month,
		"day":        day,
		"hour":       hour,
		"minute":     minute,
		"second":     second,
		"nanosecond": fraction,
	}, nil
}

// Date returns the Time corresponding to
//
// yyyy-mm-dd hh:mm:ss + nsec nanoseconds
func Date(year, month, day, hour, min, sec, nsec int) (*NepaliTime, error) {
	englishDate, err := NepaliToEnglish(year, month, day)
	if err != nil {
		return nil, err
	}

	englishTime := time.Date(englishDate[0], time.Month(englishDate[1]), englishDate[2],
		hour, min, sec, nsec, GetNepaliLocation())
	return &NepaliTime{year, month, day, &englishTime}, nil
}

// Converts Time object to NepaliTime
func FromEnglishTime(englishTime time.Time) (*NepaliTime, error) {
	englishTime = englishTime.In(GetNepaliLocation())
	enYear, enMonth, enDay := englishTime.Date()
	englishDate, err := EnglishToNepali(enYear, int(enMonth), enDay)
	if err != nil {
		return nil, err
	}

	return &NepaliTime{englishDate[0], englishDate[1], englishDate[2], &englishTime}, nil
}

// Now returns the current nepali time.
// this function should always work
// and should not return nil in normal circumstances
func Now() *NepaliTime {
	now, _ := FromEnglishTime(GetCurrentEnglishTime())

	return now
}

// adds zero on the number if the number is less than 10
// Converts single digit number into two digits.
// Adds zero on the number if the number is less than 10.
//
// eg. 19 => 19 and 8 => 08
//
// Note: if number is 144 it will return 144
func twoDigitNumber(number int) string {
	if number < 10 && number >= 0 {
		return fmt.Sprintf("0%d", number)
	}
	return fmt.Sprint(number)
}

// Gets current English date along with time level precision.
// Current Time of Asia/Kathmandu
func GetCurrentEnglishTime() time.Time {
	return time.Now().In(GetNepaliLocation())
}

// Returns location for Asia/Kathmandu (constants.Timezone)
func GetNepaliLocation() *time.Location {
	location, _ := time.LoadLocation(Timezone)
	// ignoring error since location with Asia/Kathmandu will not fail.
	return location
}
