package handlers

import (
	"registration/services"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	service *services.AuthService
}

type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

func NewAuthHandler(serivce *services.AuthService) *AuthHandler {
	return &AuthHandler{
		service: serivce,
	}
}

func (ah *AuthHandler) Register(ctx *gin.Context) {
	var request RegisterRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(422, gin.H{
			"error": "validation error",
		})
		return
	}

	if err := ah.service.RegisterUser(request.Username, request.Email, request.Password); err != nil {
		ctx.JSON(500, gin.H{
			"error": "failed to register",
		})
		return
	}

	ctx.JSON(201, gin.H{
		"message": "successfully registered user",
	})
}

func (ah *AuthHandler) GetUserByEmail(ctx *gin.Context) {
	email := ctx.Param("email")
	// ctx.Query()
	user, err := ah.service.GetUserByEmail(email)
	if err != nil {
		ctx.JSON(404, gin.H{
			"error": "user not found",
		})
	}

	ctx.JSON(200, gin.H{
		"id":       user.Id,
		"email":    user.Email,
		"username": user.Username,
	})

}
