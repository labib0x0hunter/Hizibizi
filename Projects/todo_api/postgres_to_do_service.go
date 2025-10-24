package main

import "database/sql"

type PostgreSqlTaskManager struct {
	db *sql.DB
}

func NewPostgreSqlTaskManager(driverName string, dsn string) (*PostgreSqlTaskManager, error) {
	db, err := openDatabase(driverName, dsn)
	if err != nil {
		return nil, err
	}
	return &PostgreSqlTaskManager{db: db}, nil
}

// InComplete

func (p *PostgreSqlTaskManager) AddTask(task string) error {
	return nil
}
func (p *PostgreSqlTaskManager) ListTasks() []string {
	return []string{}
}
func (p *PostgreSqlTaskManager) DeleteTask(id int) error {
	return nil
}
