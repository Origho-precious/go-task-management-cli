package main

import (
	"log"
	"strconv"
	"strings"
	"time"
)

func main() {
	handlePrompt()

	dateSlice := strings.Split(dueDate, "-")

	year, yearErr := strconv.ParseInt(dateSlice[2], 10, 32)
	if yearErr != nil {
		log.Fatal(yearErr)
	}
	month, monthErr := strconv.ParseInt(dateSlice[1], 10, 32)
	if monthErr != nil || month < 1 || month > 12 {
		log.Fatal(monthErr)
	}
	day, dayErr := strconv.ParseInt(dateSlice[0], 10, 32)
	if dayErr != nil || day < 1 || !isValidDay(day, month) {
		log.Fatal(monthErr)
	}

	date := time.Date(
		int(year), time.Month(month), int(day), 24, 0, 0, 0, time.UTC,
	)

	todo := Task{
		Completed: false,
		Text:      text,
		CreatedAt: time.Now(),
		DueAt:     date,
	}

	todo.saveTask()
}
