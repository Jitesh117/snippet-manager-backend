package handlers

import (
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

func HandleSnippets(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getAllSnippets(w)
	case http.MethodPost:
		createSnippet(w, r)
	default:
		http.Error(w, constants.ErrMethodNotAllowed, http.StatusMethodNotAllowed)
	}
}

func getAllSnippets(w http.ResponseWriter) {
	snippets, err := database.GetAllSnippets()
	if err != nil {
		http.Error(w, constants.ErrFailedToGetSnippets, http.StatusInternalServerError)
		return
	}
	log.Println("Got all Snippets")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(snippets)
}

func createSnippet(w http.ResponseWriter, r *http.Request) {
	var requestSnippet models.Snippet
	var userID uuid.UUID

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
		http.Error(w, constants.ErrFailedToExtractTokenID+": "+err.Error(), http.StatusBadRequest)
		return
	}

	snippet, err := database.CreateSnippet(
		requestSnippet.Title,
		requestSnippet.Language,
		requestSnippet.Content,
		userID,
	)
	if err != nil {
		http.Error(w, constants.ErrFailedToCreateSnippet, http.StatusInternalServerError)
		log.Println(err)
		return
	}
	log.Println("Created snippet!")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(snippet)
}
