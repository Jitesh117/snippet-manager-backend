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
