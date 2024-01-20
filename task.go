package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"regexp"
	"time"

	"github.com/charmbracelet/huh"
)

type Task struct {
	Completed bool      `json:"completed"`
	Text      string    `json:"text"`
	CreatedAt time.Time `json:"createdAt"`
	DueAt     time.Time `json:"dueAt"`
}

var (
	text    string
	dueDate string
)

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
	f, openErr := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

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

func handlePrompt() {
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title(
					"Enter task text: e.g Complete chapter 8 of Learning Go by Jon Bodner",
				).
				CharLimit(100).
				Value(&text).
				Validate(func(desc string) error {
					if desc == "" {
						return errors.New("sorry, you need to enter the task description")
					}
					return nil
				}),

			huh.NewInput().
				Title("When is it due? format: dd-mm-yyyy e.g 22-07-2024").
				CharLimit(10).
				Value(&dueDate).
				Validate(func(date string) error {
					dateFormat := regexp.MustCompile(`^\d{2}-\d{2}-\d{4}$`)
					if !dateFormat.MatchString(date) {
						return errors.New("sorry, format must be dd-mm-yyyy e.g 20-01-2024")
					}
					return nil
				}),
		),
	)

	err := form.Run()

	if err != nil {
		log.Fatal(err)
	}
}
