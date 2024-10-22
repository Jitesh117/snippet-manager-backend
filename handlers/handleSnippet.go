package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/Jitesh117/snippet-manager-backend/constants"
	"github.com/Jitesh117/snippet-manager-backend/database"
	"github.com/Jitesh117/snippet-manager-backend/helper"
	auth "github.com/Jitesh117/snippet-manager-backend/middleware"
	"github.com/Jitesh117/snippet-manager-backend/models"
	"github.com/google/uuid"
)

func HandleSnippet(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Path[len("/snippets/"):]
	snippetID, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, constants.ErrInvalidSnippetID, http.StatusBadRequest)
		return
	}
	switch r.Method {
	case http.MethodGet:
		getSnippetByID(w, r, snippetID)
	case http.MethodPut:
		updateSnippetByID(w, r, snippetID)
	case http.MethodDelete:
		deleteSnippetByID(w, r, snippetID)
	default:
		http.Error(w, constants.ErrMethodNotAllowed, http.StatusMethodNotAllowed)
	}
}

func updateSnippetByID(w http.ResponseWriter, r *http.Request, snippetID uuid.UUID) {
	var requestSnippet models.Snippet
	if err := json.NewDecoder(r.Body).Decode(&requestSnippet); err != nil {
		http.Error(w, constants.ErrInvalidPayload, http.StatusBadRequest)
		return
	}
	if err := helper.ValidateSnippet(requestSnippet); err != nil {
		http.Error(w, constants.ErrInvalidPayload+": "+err.Error(), http.StatusBadRequest)
		return
	}
	userID, err := auth.ExtractUserIDFromToken(r)
	if err != nil {
		http.Error(w, constants.ErrFailedToGetUserID+": "+err.Error(), http.StatusBadRequest)
		return
	}

	snippet, err := database.UpdateSnippet(
		requestSnippet.Title,
		requestSnippet.Language,
		requestSnippet.Content,
		snippetID,
		userID,
	)
	if err != nil {
		http.Error(w, constants.ErrFailedToUpdateSnippet, http.StatusInternalServerError)
		log.Println(err)
		return
	}

	log.Println("Updated snippet!")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(snippet)
}

func getSnippetByID(w http.ResponseWriter, r *http.Request, snippetID uuid.UUID) {
	userID, err := auth.ExtractUserIDFromToken(r)
	if err != nil {
		http.Error(w, constants.ErrFailedToGetUserID+": "+err.Error(), http.StatusBadRequest)
		return
	}

	snippet, err := database.GetSnippetByID(snippetID, userID)
	if err != nil {
		http.Error(
			w,
			constants.ErrFailedToGetSnippets+": "+err.Error(),
			http.StatusInternalServerError,
		)
		log.Println(err)
		return
	}
	log.Println("Snippet fetched from ID")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(snippet)
}

func deleteSnippetByID(w http.ResponseWriter, r *http.Request, snippetID uuid.UUID) {
	userID, err := auth.ExtractUserIDFromToken(r)
	if err != nil {
		http.Error(w, constants.ErrFailedToGetUserID+": "+err.Error(), http.StatusBadRequest)
		return
	}

	snippet, err := database.DeleteSnippetByID(snippetID, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, constants.ErrSnippetNotFound, http.StatusBadRequest)
		}
		http.Error(w, constants.ErrFailedToDeleteSnippet, http.StatusInternalServerError)
	}
	log.Println("snippet deleted from DB")
	json.NewEncoder(w).Encode(snippet)
}
