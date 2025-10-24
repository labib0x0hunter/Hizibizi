package repositories

import "registration/models"

type UserRepository interface {
	Create(user models.User) error
	GetByEmail(email string) (*models.User, error)
}
