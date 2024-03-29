package main

import (
	"database/sql"
	"time"
)

type Task struct {
	Id          int
	Completed   bool
	Description string
	CreatedAt   time.Time
	DueBy       time.Time
}

func (t Task) saveTask(db *sql.DB) Task {
	var taskRow Task
	row := db.QueryRow(`
		INSERT INTO tasks(description, created_at, due_by, completed)
		VALUES ($1, $2, $3, $4) 
		RETURNING id, description, created_at, due_by, completed`,
		t.Description, t.CreatedAt, t.DueBy, false,
	)

	if row.Err() != nil {
		panic(row.Err())
	}

	err := row.Scan(
		&taskRow.Id, &taskRow.Description, &taskRow.CreatedAt,
		&taskRow.DueBy, &taskRow.Completed,
	)

	if err != nil {
		panic(err)
	}

	return taskRow
}

func showAllTasks(db *sql.DB) []Task {
	var taskRows []Task

	rows, err := db.Query("SELECT * FROM tasks")

	if err != nil {
		panic(err)
	}

	for rows.Next() {
		var taskRow Task

		err = rows.Scan(
			&taskRow.Id, &taskRow.Description,
			&taskRow.CreatedAt, &taskRow.DueBy, &taskRow.Completed,
		)

		if err != nil {
			panic(err)
		}

		taskRows = append(taskRows, taskRow)
	}

	if rows.Err() != nil {
		panic(rows.Err())
	}

	return taskRows
}

func markTaskAsCompleted(db *sql.DB, id string) []Task {
	_, err := db.Exec(`
		UPDATE tasks
		SET completed = true
		WHERE id = $1
	`, id)

	if err != nil {
		panic(err)
	}

	updatedTasks := showAllTasks(db)

	return updatedTasks
}

func showUncompletedTasks(db *sql.DB) []Task {
	var taskRows []Task

	rows, err := db.Query(`
		SELECT * FROM tasks
		WHERE completed = false
	`)

	if err != nil {
		panic(err)
	}

	for rows.Next() {
		var taskRow Task

		err = rows.Scan(
			&taskRow.Id, &taskRow.Description,
			&taskRow.CreatedAt, &taskRow.DueBy, &taskRow.Completed,
		)

		if err != nil {
			panic(err)
		}

		taskRows = append(taskRows, taskRow)
	}

	if rows.Err() != nil {
		panic(rows.Err())
	}

	return taskRows
}

func deleteTask(db *sql.DB, id string) []Task {
	_, err := db.Exec(`
		DELETE FROM tasks
		WHERE id = $1
	`, id)

	if err != nil {
		panic(err)
	}

	remainingTasks := showAllTasks(db)

	return remainingTasks
}
