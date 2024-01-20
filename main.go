package main

import (
	"errors"
	"log"
	"regexp"
	"time"

	"github.com/charmbracelet/huh"
)

var (
	text  string
	date  string
	dueAt time.Time
)

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
				Value(&date).
				Validate(func(date string) error {
					dateFormat := regexp.MustCompile(`^\d{2}-\d{2}-\d{4}$`)
					if !dateFormat.MatchString(date) {
						return errors.New("sorry, format must be dd-mm-yyyy e.g 20-01-2024")
					}

					fDueDate, err := formatDueDate(date)

					if err != nil {
						return err
					}

					dueAt = fDueDate

					return nil
				}),
		),
	)

	err := form.Run()

	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	handlePrompt()

	todo := Task{
		Completed: false,
		Text:      text,
		CreatedAt: time.Now(),
		DueAt:     dueAt,
	}

	todo.saveTask()
}
