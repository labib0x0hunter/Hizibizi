package main

import (
	"database/sql"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"searchengine/db"
	"searchengine/handler"
	memorymapper "searchengine/memory_mapper"
	"searchengine/repositories"
	"searchengine/services"
	"searchengine/utils"
	"syscall"

	"github.com/gin-gonic/gin"
)

func main() {

	path, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	utils.Path = filepath.Join(path, "../../")

	newDb, err := db.NewDocumentMysqlDb()
	if err != nil {
		panic(err)
	}
	defer newDb.Close()

	newDict, err := memorymapper.NewDictionary()
	if err != nil {
		panic(err)
	}
	defer newDict.Close()

	newPost, err := memorymapper.NewPosting()
	if err != nil {
		panic(err)
	}
	defer newPost.Close()

	newHasher := utils.NewHash()

	// On shutdown CTRL + C
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func(newDb *sql.DB, newDict *memorymapper.Dictionary, newPost *memorymapper.Posting) {
		sig := <-sigChan
		fmt.Println("Received: ", sig)
		newDict.Close()
		newDb.Close()
		newPost.Close()
		os.Exit(0)
	}(newDb, newDict, newPost)

	docRepo := repositories.NewDocumentRepo(newDb)
	indexRepo := repositories.NewIndexRepo(newDict, newPost)
	engineService := services.NewEngineService(indexRepo, docRepo, newHasher)
	engineHandler := handler.NewEngineHandler(engineService)

	router := gin.Default()
	router.Static("/", filepath.Join(utils.Path, "static"))

	router.NoRoute(engineHandler.FrontPage)
	router.POST("/insert", engineHandler.Index)
	router.POST("/search", engineHandler.Search)

	router.Run(":8080")
}
