package db

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

// SHOW tables;

func NewMysqlDB() (*sql.DB, error) {
	databaseName := "registration"
	dsn := "root:@tcp(127.0.0.1:3306)/" + databaseName
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
