package database

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/Jitesh117/snippet-manager-backend/models"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func CreateUser(user models.User) (uuid.UUID, error) {
	// hash the password before using it in the db
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("failed to hash password: %v", err)
	}

	query := `
        INSERT INTO users (user_id, username, email, password_hash, created_at)
        VALUES (gen_random_uuid(), $1, $2, $3, $4)
        RETURNING user_id
    `

	var userID uuid.UUID

	err = DB.QueryRow(query, user.UserName, user.Email, hashedPassword, time.Now().UTC()).
		Scan(&userID)
	if err != nil {
		return uuid.UUID{}, err
	}

	return userID, nil
}

func CheckUserCredentials(email, password string) (uuid.UUID, error) {
	query := `
  SELECT user_id, password_hash
  FROM users
  WHERE email = $1
  `

	var userID uuid.UUID
	var passwordHash string

	err := DB.QueryRow(query, email).Scan(&userID, &passwordHash)
	if err != nil {
		if err == sql.ErrNoRows {
			return uuid.UUID{}, fmt.Errorf("invalid credentials")
		}
		return uuid.UUID{}, fmt.Errorf("failed to retreive user: %v", err)
	}
	err = bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password))
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("invalid credentials")
	}
	return userID, nil
}

func DeleteUser(userID uuid.UUID) (uuid.UUID, error) {
	query := `
  DELETE FROM users
  where user_id = $1
  RETURNING user_id;
  `

	var deletedUserID uuid.UUID

	err := DB.QueryRow(query, userID).Scan(&deletedUserID)
	if err != nil {
		if err == sql.ErrNoRows {
			return uuid.Nil, fmt.Errorf("user with ID %s not found", userID)
		}
		return uuid.Nil, err
	}
	return deletedUserID, nil
}

func ChangePassword(userID uuid.UUID, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %v", err)
	}
	query := `UPDATE users
  SET password_hash = $1
  WHERE user_id = $2
  RETURNING user_id
  `
	var updatedUserId uuid.UUID

	err = DB.QueryRow(query, hashedPassword, userID).Scan(&updatedUserId)
	if err != nil {
		return err
	}
	return nil
}
