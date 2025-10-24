package main

import (
	"fmt"
	"registration/db"
	"registration/handlers"
	"registration/redis"
	"registration/repositories"
	"registration/services"

	"github.com/gin-gonic/gin"
)

func main() {
	db, err := db.NewMysqlDB()
	if err != nil {
		panic("Failed to connect to MySQL: " + err.Error())
	}
	defer db.Close()

	redis, err := redis.NewRedisClient()
	if err != nil {
		panic("Failed to connect to Redis: " + err.Error())
	}
	defer redis.Close()

	fmt.Println("Connected")

	mysqlUserRepo := repositories.NewMysqlUserRepository(db)
	redisMysqlUserRepo := repositories.NewRedisMysqlUserRepository(mysqlUserRepo, redis)
	authService := services.NewAuthService(redisMysqlUserRepo)
	authHandler := handlers.NewAuthHandler(authService)

	router := gin.Default()

	router.POST("/register", authHandler.Register)
	router.GET("/user/:email", authHandler.GetUserByEmail)

	router.Run(":8080")
}
