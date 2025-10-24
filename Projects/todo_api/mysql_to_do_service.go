package main

import "database/sql"

type MySqlTaskManager struct {
	db *sql.DB
}

func NewMySqlTaskManager(driverName, dsn string) (*MySqlTaskManager, error) {
	db, err := openDatabase(driverName, dsn)
	if err != nil {
		return nil, err
	}
	return &MySqlTaskManager{db: db}, nil
}

// Incomplete

func (m *MySqlTaskManager) AddTask(task string) error {
	return nil
}
func (m *MySqlTaskManager) ListTasks() []string {
	return []string{}
}
func (m *MySqlTaskManager) DeleteTask(id int) error {
	return nil
}
