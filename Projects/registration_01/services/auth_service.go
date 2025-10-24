package services

import (
	"fmt"
	"registration/models"
	"registration/repositories"
	"registration/utils"

	"github.com/google/uuid"
)

// UserRepository -> interface -> abstract type -> mock
// UserRepository -> struct -> concrete type -> no mock

type AuthService struct {
	UserRepo repositories.UserRepository // interface
	// UserRepo *repositories.UserRepository // pointer-> Error -> interface is reference type
}

func NewAuthService(userRepo repositories.UserRepository) *AuthService {
	return &AuthService{
		UserRepo: userRepo,
	}
}

func (as *AuthService) RegisterUser(username, email, password string) error {
	passwordHash, err := utils.HashPassword(password)
	if err != nil {
		return err
	}
	verificationToken := uuid.New().String()

	// err -> nil, user -> nil
	existingUser, err := as.UserRepo.GetByEmail(email)
	if err != nil || existingUser != nil {
		return fmt.Errorf("user exists: %w", err)
	}

	user := models.User{
		Username:          username,
		Email:             email,
		PasswordHash:      passwordHash,
		IsVerified:        false,
		VerificationToken: verificationToken,
	}

	if err := as.UserRepo.Create(user); err != nil {
		return err
	}

	// send verification email

	return nil
}

func (as *AuthService) GetUserByEmail(email string) (*models.User, error) {
	return as.UserRepo.GetByEmail(email)
}
