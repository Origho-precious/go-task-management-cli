package main

import (
	"database/sql"
	"time"
)

type Task struct {
	Completed   bool
	Description string
	CreatedAt   time.Time
	DueBy       time.Time
}

type TaskRow struct {
	Task
	Id int
}

func (t Task) saveTask(db *sql.DB) TaskRow {
	var taskRow TaskRow
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

func showAllTasks(db *sql.DB) []TaskRow {
	var taskRows []TaskRow

	rows, err := db.Query(`SELECT * FROM tasks`)

	if err != nil {
		panic(err)
	}

	defer db.Close()

	for rows.Next() {
		var taskRow TaskRow

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
