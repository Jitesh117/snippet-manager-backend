package models

import (
	"time"

	"github.com/google/uuid"
)

type Snippet struct {
	SnippetId uuid.UUID `json:"snippet_id"`
	Title     string    `json:"title"`
	Language  string    `json:"language"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type User struct {
	UserID    uuid.UUID `json:"user_id"`
	UserName  string    `json:"user_name"`
	Email     string    `json:"email"`
	Password  string    `json:"password,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}
