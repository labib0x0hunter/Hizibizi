package handler

import (
	"path/filepath"
	"searchengine/models"
	"searchengine/services"
	"searchengine/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

type EngineHandler struct {
	engine *services.EngineService
}

type DocumentRequest struct {
	Document string `json:"document" binding:"required"`
}

func NewEngineHandler(engine *services.EngineService) *EngineHandler {
	return &EngineHandler{
		engine: engine,
	}
}

func (e *EngineHandler) Index(ctx *gin.Context) {
	var request DocumentRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(422, gin.H{
			"error" : "validation error",
		})
		return
	}

	if err := e.engine.IndexDocument(request.Document); err != nil {
		ctx.JSON(500, gin.H{
			"error" : "failed to store document",
		})
		return
	}

	ctx.JSON(200, gin.H{
		"msg" : "document inserted",
	})
}

func (e *EngineHandler) Search(ctx *gin.Context) {
	var request DocumentRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(422, gin.H{
			"error" : "validation error",
		})
		return
	}
	
	var documents []models.Document
	for index, doc := range e.engine.SearchDocument(request.Document) {
		documents = append(documents, models.Document{
			DocId: "Doc_" + strconv.Itoa(index),
			Document: doc,
		})
	}

	ctx.JSON(200, documents)
}

func (e *EngineHandler) FrontPage(ctx *gin.Context) {
	ctx.File(filepath.Join(utils.Path, "static", "index.html"))
}