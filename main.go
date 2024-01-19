package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
	// "github.com/charmbracelet/huh"
)

type Task struct {
	Completed   bool      `json:"isCompleted"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"createdAt"`
	DueAt       time.Time `json:"dueAt"`
}

func (t Task) saveTask() {
	// Creating file name from date created
	formattedDate := t.CreatedAt.Format("02-01-2006 15:04:05")
	path := fmt.Sprintf("./tasks/%v.json", formattedDate)

	// Creating "tasks" folder is it isn't created yet
	err := os.MkdirAll("tasks", 0755)

	if err != nil {
		panic(err)
	}

	// Open file to write
	f, openErr := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, 0644)

	if openErr != nil {
		fmt.Println("Path: ", path)
		panic(openErr)
	}

	defer f.Close()

	// Create json
	b, jsonErr := json.Marshal(t)

	if jsonErr != nil {
		panic(jsonErr)
	}

	// Write json to file
	_, writeErr := f.Write(b)

	if writeErr != nil {
		f.Close()
		panic(writeErr)
	}

	fmt.Println("Task(s) saved to file:", path)
}

func main() {
	date := time.Date(2023, time.January, 26, 12, 0, 0, 0, time.UTC)

	todo := Task{
		Completed:   false,
		Description: "Complete Chapter 8 of Go book",
		CreatedAt:   time.Now(),
		DueAt:       date,
	}

	todo.saveTask()
}
