package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

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

func formatDueDate(dueDate string) (time.Time, error) {
	dateSlice := strings.Split(dueDate, "-")

	year, yearErr := strconv.ParseInt(dateSlice[2], 10, 32)
	if yearErr != nil {
		return time.Time{}, yearErr
	}
	month, monthErr := strconv.ParseInt(dateSlice[1], 10, 32)
	if monthErr != nil || month < 1 || month > 12 {
		if monthErr != nil {
			return time.Time{}, monthErr
		}
		return time.Time{}, errors.New("invalid month, must be within 01 and 12")
	}
	day, dayErr := strconv.ParseInt(dateSlice[0], 10, 32)
	if dayErr != nil || day < 1 || day > 31 || !isValidDay(day, month) {
		if dayErr != nil {
			return time.Time{}, dayErr
		}
		return time.Time{}, errors.New("invalid day of the month")
	}

	date := time.Date(
		int(year), time.Month(month), int(day), 24, 0, 0, 0, time.UTC,
	)

	return date, nil
}

func renderTasks(tasks []Task) {
	// Print table header
	fmt.Printf("\n%-11s%-62s%-15s%-15s\n",
		"Completed", "Description", "Due date", "Created At",
	)

	// Print table body
	for _, task := range tasks {
		completed := "false"

		if task.Completed {
			completed = "true"
		}

		fmt.Printf(
			"%-11s%-62s%-15s%-15s\n",
			completed, task.Description,
			task.DueBy.Format("02-01-2006"),
			task.CreatedAt.Format("02-01-2006"),
		)
	}
}
