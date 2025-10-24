package repositories

import (
	"database/sql"
	"log/slog"
	"registration/models"
)

type UserRepository struct {
	db *sql.DB
}

// UserRepository is depended on sql.DB (dependency)
func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (u *UserRepository) Create(user models.User) error {
	query := "INSERT INTO users (username, email, password, is_verified, verification_token) VALUES (?, ?, ?, ?, ?)"
	_, err := u.db.Exec(query, user.Username, user.Email, user.PasswordHash, user.IsVerified, user.VerificationToken)
	if err != nil {
		slog.Error("at user creation")
	}

	return err
}

func (u *UserRepository) GetByEmail(email string) (*models.User, error) {
	query := "SELECT id, username, email, password, is_verified, verification_token FROM users WHERE email = ?"
	row := u.db.QueryRow(query, email)

	var user models.User
	err := row.Scan(&user.Id, &user.Username, &user.Email, &user.PasswordHash, &user.IsVerified, &user.VerificationToken)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil	// User not found
		}
		slog.Error("retriving user ", err.Error())
		return nil, err
	}
	return &user, nil
}