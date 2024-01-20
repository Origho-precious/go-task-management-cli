package main

import "time"

func isValidDay(day, month int64) bool {
	if day < 1 {
		return false
	}

	switch time.Month(month) {
	case time.September:
	case time.April:
	case time.June:
	case time.November:
		if day > 30 {
			return false
		}

		return true
	case time.January:
	case time.March:
	case time.May:
	case time.July:
	case time.August:
	case time.October:
	case time.December:
		if day > 31 {
			return false
		}

		return true
	default:
		if day > 29 {
			return false
		}

		return true
	}

	return true
}
