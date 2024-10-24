package database

import (
	"database/sql"
	"fmt"

	"github.com/Jitesh117/snippet-manager-backend/constants"
	"github.com/Jitesh117/snippet-manager-backend/helper"
	"github.com/Jitesh117/snippet-manager-backend/models"
	"github.com/google/uuid"
)

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

func GetSnippetsByLanguage(language string, userID uuid.UUID) ([]models.Snippet, error) {
	var snippets []models.Snippet
	query := `
    SELECT snippet_id, title, content, language, created_at, updated_at 
    FROM snippets 
    WHERE language = $1 AND user_id = $2
    ORDER BY updated_at DESC
    `
	rows, err := DB.Query(query, language, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var snippet models.Snippet
		err := rows.Scan(
			&snippet.SnippetId,
			&snippet.Title,
			&snippet.Content,
			&snippet.Language,
			&snippet.CreatedAt,
			&snippet.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		snippets = append(snippets, snippet)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return snippets, nil
}

func GetSnippetsSorted(userID uuid.UUID, sortBy, order string) ([]models.Snippet, error) {
	if !helper.IsValidSortField(sortBy) || !helper.IsValidOrder(order) {
		return nil, fmt.Errorf("invalid sort options")
	}

	query := `
		SELECT snippet_id, title, content, language, created_at, updated_at 
		FROM snippets 
		WHERE user_id = $1
		ORDER BY %s %s
	`

	query = fmt.Sprintf(query, sortBy, order)

	rows, err := DB.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var snippets []models.Snippet
	for rows.Next() {
		var snippet models.Snippet
		err := rows.Scan(
			&snippet.SnippetId,
			&snippet.Title,
			&snippet.Content,
			&snippet.Language,
			&snippet.CreatedAt,
			&snippet.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		snippets = append(snippets, snippet)
	}
	if len(snippets) == 0 {
		return nil, fmt.Errorf(constants.ErrSnippetNotFound)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return snippets, nil
}
