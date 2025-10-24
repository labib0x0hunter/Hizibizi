package models

// CREATE TABLE users (
//     id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
//     username VARCHAR(50) NOT NULL UNIQUE,
//     email VARCHAR(100) NOT NULL UNIQUE,
//     password VARCHAR(255) NOT NULL,
//     is_verified BOOLEAN NOT NULL DEFAULT FALSE,
//     verification_token VARCHAR(100) DEFAULT NULL,
//     created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
//     updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
// );

type User struct {
	Id                int
	Username          string
	Email             string
	PasswordHash      string
	IsVerified        bool
	VerificationToken string
}
