package database

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/Jitesh117/snippet-manager-backend/models"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

var DB *sql.DB

func InitDB() {
	connStr := "host=localhost port=5432 user=postgres password=mysecretpassword dbname=snippet_manager sslmode=disable"
	var err error
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}

	initQuery := `
    CREATE EXTENSION IF NOT EXISTS "pgcrypto";

    -- Create the users table
    CREATE TABLE IF NOT EXISTS users (
        user_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
        username TEXT NOT NULL UNIQUE,
        email TEXT NOT NULL UNIQUE,
        password_hash TEXT NOT NULL,
        created_at TIMESTAMP WITH TIME ZONE DEFAULT (CURRENT_TIMESTAMP AT TIME ZONE 'UTC')
    );

    -- Create the snippets table with a foreign key referencing users
    CREATE TABLE IF NOT EXISTS snippets (
        snippet_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
        title TEXT NOT NULL,
        language TEXT NOT NULL,
        content TEXT NOT NULL,
        user_id UUID NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
        created_at TIMESTAMP WITH TIME ZONE DEFAULT (CURRENT_TIMESTAMP AT TIME ZONE 'UTC'),
        updated_at TIMESTAMP WITH TIME ZONE DEFAULT (CURRENT_TIMESTAMP AT TIME ZONE 'UTC')
    );
  `
	_, err = DB.Exec(initQuery)
	if err != nil {
		log.Fatal("Failed to create tables: ", err)
	}

	err = DB.Ping()
	if err != nil {
		log.Fatal("Failed to ping database: ", err)
	}
	log.Println("Connected to database!")
}

func CloseDB() {
	if DB != nil {
		err := DB.Close()
		if err != nil {
			log.Println("Failed to close the database connection: ", err)
		}
	}
}

func GetAllSnippets() ([]models.Snippet, error) {
	query := "SELECT snippet_id, title, language, content, created_at, updated_at FROM snippets"
	rows, err := DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	snippets := []models.Snippet{}
	for rows.Next() {
		var tempSnippet models.Snippet
		if err = rows.Scan(&tempSnippet.SnippetId, &tempSnippet.Title, &tempSnippet.Language, &tempSnippet.Content, &tempSnippet.CreatedAt, &tempSnippet.UpdatedAt); err != nil {
			return nil, err
		}
		snippets = append(snippets, tempSnippet)
	}
	return snippets, nil
}

func CreateSnippet(
	title string,
	language string,
	content string,
	userID uuid.UUID,
) (models.Snippet, error) {
	var snippet models.Snippet
	query := `
		INSERT INTO snippets (title, language, content, user_id) 
		VALUES ($1, $2, $3, $4) 
		RETURNING snippet_id, title, language, content, created_at, updated_at
	`
	err := DB.QueryRow(query, title, language, content, userID).Scan(
		&snippet.SnippetId,
		&snippet.Title,
		&snippet.Language,
		&snippet.Content,
		&snippet.CreatedAt,
		&snippet.UpdatedAt,
	)
	if err != nil {
		return models.Snippet{}, err
	}
	return snippet, nil
}

func UpdateSnippet(
	title string,
	language string,
	content string,
	snippetID uuid.UUID,
	userID uuid.UUID,
) (models.Snippet, error) {
	var realUserID uuid.UUID
	verifyQuery := "SELECT user_id FROM snippets WHERE snippet_id = $1"
	err := DB.QueryRow(verifyQuery, snippetID).Scan(&realUserID)
	if err != nil {
		return models.Snippet{}, err
	}
	if realUserID != userID {
		return models.Snippet{}, fmt.Errorf("access denied")
	}
	var snippet models.Snippet
	query := `
		UPDATE snippets 
		SET title = $1, language = $2, content = $3, updated_at = NOW() AT TIME ZONE 'UTC'
		WHERE snippet_id = $4 
		RETURNING snippet_id, title, language, content, created_at, updated_at
	`

	err = DB.QueryRow(query, title, language, content, snippetID).Scan(
		&snippet.SnippetId,
		&snippet.Title,
		&snippet.Language,
		&snippet.Content,
		&snippet.CreatedAt,
		&snippet.UpdatedAt,
	)
	if err != nil {
		return models.Snippet{}, err
	}

	return snippet, nil
}

func GetSnippetByID(snippetID uuid.UUID, userID uuid.UUID) (models.Snippet, error) {
	var realUserID uuid.UUID
	verifyQuery := "SELECT user_id FROM snippets WHERE snippet_id = $1"
	err := DB.QueryRow(verifyQuery, snippetID).Scan(&realUserID)
	if err != nil {
		return models.Snippet{}, err
	}
	if realUserID != userID {
		return models.Snippet{}, fmt.Errorf("access denied")
	}
	var snippet models.Snippet
	query := "SELECT snippet_id, title, language, content, created_at, updated_at FROM snippets WHERE snippet_id = $1"

	err = DB.QueryRow(query, snippetID).Scan(
		&snippet.SnippetId,
		&snippet.Title,
		&snippet.Language,
		&snippet.Content,
		&snippet.CreatedAt,
		&snippet.UpdatedAt,
	)
	if err != nil {
		return models.Snippet{}, err
	}

	return snippet, nil
}

func DeleteSnippetByID(snippetID uuid.UUID, userID uuid.UUID) (models.Snippet, error) {
	var realUserID uuid.UUID
	verifyQuery := "SELECT user_id FROM snippets WHERE snippet_id = $1"
	err := DB.QueryRow(verifyQuery, snippetID).Scan(&realUserID)
	if err != nil {
		return models.Snippet{}, err
	}
	if realUserID != userID {
		return models.Snippet{}, fmt.Errorf("access denied")
	}
	var snippet models.Snippet
	selectQuery := `SELECT snippet_id, title, language, content, created_at, updated_at 
                    FROM snippets 
                    WHERE snippet_id = $1`
	err = DB.QueryRow(selectQuery, snippetID).Scan(
		&snippet.SnippetId,
		&snippet.Title,
		&snippet.Language,
		&snippet.Content,
		&snippet.CreatedAt,
		&snippet.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.Snippet{}, fmt.Errorf("snippet with ID %s not found", snippetID)
		}
		return models.Snippet{}, err
	}
	deleteQuery := "DELETE FROM snippets where snippet_id = $1"
	_, err = DB.Exec(deleteQuery, snippetID)
	if err != nil {
		return models.Snippet{}, err
	}
	return snippet, nil
}

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
