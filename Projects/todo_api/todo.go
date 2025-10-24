package main

import (
	"database/sql"
	"fmt"
)

type TodoService interface {
	AddTask(task string) error
	ListTasks() []string
	DeleteTask(id int) error
}

func openDatabase(driverName, dsn string) (*sql.DB, error) {
	db, err := sql.Open(driverName, dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open db: %w", err)
	}
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to connect: %w", err)
	}
	return db, nil
}

func AddTask(ts TodoService, task string) error {
	return ts.AddTask(task)
}

func ListTasks(ts TodoService) []string {
	return ts.ListTasks()
}

func DeleteTask(ts TodoService, id int) error {
	return ts.DeleteTask(id)
}
