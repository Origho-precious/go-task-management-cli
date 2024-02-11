package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/Origho-precious/go-task-management-cli/configs"
	"github.com/charmbracelet/huh"
)

var (
	text    string
	date    string
	addMore string
	DueBy   time.Time
)

func getAction() (string, error) {
	var action string

	firstForm := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Choose one").
				Options(
					huh.NewOption("I want to view all my tasks", "view"),
					huh.NewOption("I want to add new task(s)", "add"),
					huh.NewOption("I want to mark task(s) as completed", "update"),
					huh.NewOption("I want to delete task(s)", "delete"),
				).
				Value(&action),
		),
	)
	err := firstForm.Run()

	if err != nil {
		log.Fatal(err)
	}

	return action, err
}

func handleView(db *sql.DB) {
	tasks := showAllTasks(db)

	// Print table header
	fmt.Printf("%-5s%-11s%-20s%-12s%-12s\n",
		"ID", "Completed", "Description", "CreatedAt", "DueBy",
	)

	// Print table rows
	for _, task := range tasks {
		completed := "false"

		if task.Completed {
			completed = "true"
		}

		fmt.Printf(
			"%-5d%-11s%-20s%-12s%-12s\n",
			task.Id, completed, task.Description,
			task.CreatedAt.Format("02-01-2006"),
			task.DueBy.Format("02-01-2006"),
		)
	}
}

func handlePrompt(db *sql.DB) {
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

					DueBy = fDueDate

					return nil
				}),

			huh.NewSelect[string]().
				Title("Do you want to add more?").
				Options(
					huh.NewOption("yes", "yes"),
					huh.NewOption("no", "no"),
				).
				Value(&addMore),
		),
	)

	err := form.Run()

	if err != nil {
		log.Fatal(err)
	}

	todo := Task{
		Completed:   false,
		Description: text,
		CreatedAt:   time.Now(),
		DueBy:       DueBy,
	}

	savedTask := todo.saveTask(db)

	fmt.Printf("Taak(s) saved")
	fmt.Printf("Description: %s, Due by: %s\n",
		savedTask.Description, savedTask.DueBy,
	)

	if addMore == "yes" {
		text = ""
		handlePrompt(db)
	}
}

func main() {
	db, err := configs.ConnectDB()

	if err != nil {
		panic(err)
	}

	defer db.Close()

	// Ping the database
	err = db.Ping()

	if err != nil {
		panic(err)
	}

	fmt.Println("Database connection established!")

	// Creating Tables
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS tasks(
			id SERIAL PRIMARY KEY,
			description TEXT NOT NULL,
			created_at TIMESTAMP NOT NULL,
			due_by TIMESTAMP NOT NULL,
			completed BOOLEAN DEFAULT false
		);
	`)

	// CREATE TABLE IF NOT EXISTS users(
	// 		id SERIAL PRIMARY KEY,
	// 		fullName TEXT NOT NULL,
	// 		email TEXT UNIQUE NOT NULL,
	// 		password TEXT NOT NULL
	// 	);

	if err != nil {
		panic(err)
	}

	action, actionErr := getAction()

	if actionErr != nil {
		panic(err)
	}

	switch action {
	case "view":
		handleView(db)
	case "add":
		handlePrompt(db)
	case "update":
		// handlePrompt(db)
	case "delete":
		// handlePrompt(db)
	default:
		fmt.Println("Invalid action")
	}
}
