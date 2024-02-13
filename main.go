package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
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

func getAction(db *sql.DB) error {
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
					huh.NewOption("I want to exit the CLI", "close"),
				).
				Value(&action),
		),
	)
	err := firstForm.Run()

	if err != nil {
		log.Fatal(err)
	}

	switch action {
	case "view":
		handleView(db)
	case "add":
		handleAdd(db)
	case "update":
		handleUpdate(db)
	case "delete":
		handleDelete(db)
	default:
		fmt.Printf("You can restart the CLI by running %q\n", "go run .")
		os.Exit(0)
	}

	return err
}

func handleView(db *sql.DB) {
	tasks := showAllTasks(db)

	renderTasks(tasks)

	getAction(db)
}

func handleAdd(db *sql.DB) {
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title(
					"Enter task text: e.g Complete chapter 8 of Learning Go by Jon Bodner",
				).
				CharLimit(60).
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

					DueBy = fDueDate.AddDate(0, 0, -1)

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

	fmt.Println("Task saved!")
	fmt.Printf(
		"Description: %s, Due by: %s\n",
		savedTask.Description, savedTask.DueBy.Format("02-01-2006"),
	)

	if addMore == "yes" {
		text = ""
		handleAdd(db)
	} else {
		err = getAction(db)

		if err != nil {
			panic(err)
		}
	}
}

func handleUpdate(db *sql.DB) {
	tasks := showUncompletedTasks(db)

	if len(tasks) == 0 {
		fmt.Println("You have no uncompleted task, I commend your work rate!")
	} else {
		var taskOptions []huh.Option[string]

		for index, task := range tasks {
			taskOption := huh.NewOption(
				fmt.Sprintf(
					"#%d Description: %s | Due date: %s", index+1,
					task.Description, task.DueBy.Format("02-01-2006"),
				),
				strconv.Itoa(task.Id),
			)
			taskOptions = append(taskOptions, taskOption)
		}

		taskOptions = append(taskOptions, huh.NewOption("Go back", "back"))

		var action string

		taskForm := huh.NewForm(
			huh.NewGroup(
				huh.NewSelect[string]().
					Title("Select task to mark as completed.").
					Options(taskOptions...).
					Value(&action),
			),
		)
		err := taskForm.Run()

		if err != nil {
			log.Fatal(err)
		}

		if action != "back" {
			updatedTasks := markTaskAsCompleted(db, action)

			fmt.Println("Successful!")

			renderTasks(updatedTasks)
		}

		err = getAction(db)

		if err != nil {
			log.Fatal(err)
		}
	}
}

func handleDelete(db *sql.DB) {
	tasks := showAllTasks(db)

	if len(tasks) == 0 {
		fmt.Println(
			"You currently do not have any task. You can from the main menu",
		)
	} else {
		var taskOptions []huh.Option[string]

		for index, task := range tasks {
			taskOption := huh.NewOption(
				fmt.Sprintf(
					"#%d Description: %s | Due date: %s", index+1,
					task.Description, task.DueBy.Format("02-01-2006"),
				),
				strconv.Itoa(task.Id),
			)
			taskOptions = append(taskOptions, taskOption)
		}

		taskOptions = append(taskOptions, huh.NewOption("Go back", "back"))

		var action string

		taskForm := huh.NewForm(
			huh.NewGroup(
				huh.NewSelect[string]().
					Title("Select task you want to delete.").
					Options(taskOptions...).
					Value(&action),
			),
		)
		err := taskForm.Run()

		if err != nil {
			log.Fatal(err)
		}

		if action != "back" {
			remainingTasks := deleteTask(db, action)

			fmt.Println("Task deleted successfully!")

			renderTasks(remainingTasks)
		}

		err = getAction(db)

		if err != nil {
			log.Fatal(err)
		}
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

	err = getAction(db)

	if err != nil {
		panic(err)
	}
}
