package db

import (
	"database/sql"
	"log/slog"

	_ "github.com/go-sql-driver/mysql"
)

// CREATE DATABASE testDb;
// SHOW DATABASES;
// USE testDb;

/**
CREATE TABLE items (
    id INT AUTO_INCREMENT PRIMARY KEY,
    content TEXT NOT NULL
);
**/

func NewDocumentMysqlDb() (*sql.DB, error) {
	dsn := user + ":" + password + "@tcp(" + ip + ":" + port + ")/" + dbName
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		slog.Error(" [db.go] [NewDbService()] database open error ", err)
		return nil, err
	}
	if err := db.Ping(); err != nil {
		slog.Error(" [db.go] [NewDbService()] database ping error ", err)
		return nil, err
	}

	// DELETE database table
	if _, err := db.Exec(deleteTable + tableName); err != nil {
		slog.Error(" [db.go] [NewDbService()] database table delete error ", err)
		return nil, err
	}

	// RESET AUTO_INCREMENT
	if _, err := db.Exec(resetTable); err != nil {
		slog.Error(" [db.go] [NewDbService()] database reset AUTO_INCREMENT error ", err)
		return nil, err
	}

	return db, nil
}
