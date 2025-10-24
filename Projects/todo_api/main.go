package main

import "log"

func main() {

	dsnMysql := "username:password@tcp(127.0.0.1:3306)/yourdbname"
	dsnPostgres := "postgres://username:password@localhost:5432/yourdbname?sslmode=disable"

	mysqlService, err := NewMySqlTaskManager("mysql", dsnMysql)
	if err != nil {
		log.Panic(err)
	}
	postgreService, err := NewPostgreSqlTaskManager("postgres", dsnPostgres)
	if err != nil {
		log.Panic(err)
	}

	AddTask(mysqlService, "task from mysql")
	AddTask(postgreService, "task from postgres")
}