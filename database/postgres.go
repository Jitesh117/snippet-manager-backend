package database

import (
	"database/sql"
	"log"

	"github.com/Jitesh117/snippet-manager/models"
	_ "github.com/lib/pq"
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
    CREATE TABLE IF NOT EXISTS snippets (
        snippet_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
        title TEXT NOT NULL,
        language TEXT NOT NULL,
        content TEXT NOT NULL,
        created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
    );
  `
	_, err = DB.Exec(initQuery)
	if err != nil {
		log.Fatal("Failed to create table: ", err)
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

func CreateSnippet(title string, language string, content string) (models.Snippet, error) {
	var snippet models.Snippet
	query := `
		INSERT INTO snippets (title, language, content) 
		VALUES ($1, $2, $3) 
		RETURNING snippet_id, title, language, content, created_at, updated_at
	`
	err := DB.QueryRow(query, title, language, content).Scan(
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
