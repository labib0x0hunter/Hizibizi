package main

import (
	"fmt"
	"registration/db"
	"registration/handlers"
	"registration/repositories"
	"registration/services"

	"github.com/gin-gonic/gin"
)

func main() {
	db, err := db.NewMysqlDB()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	fmt.Println("Connected")

	userRepo := repositories.NewUserRepository(db)
	authService := services.NewAuthService(userRepo)
	authHandler := handlers.NewAuthHandler(authService)

	router := gin.Default()

	router.POST("/register", authHandler.Register)

	router.Run(":8080")
}