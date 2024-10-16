package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Jitesh117/snippet-manager/database"
	"github.com/Jitesh117/snippet-manager/models"
	"github.com/google/uuid"
)

func HandleSnippets(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getAllSnippets(w, r)
	case http.MethodPost:
		createSnippets(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func HandleSnippet(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Path[len("/snippets/"):]
	snippetID, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}
	switch r.Method {
	case http.MethodGet:
		getSnippetByID(w, r, snippetID)
	case http.MethodPut:
		updateSnippetByID(w, r, snippetID)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func getAllSnippets(w http.ResponseWriter, r *http.Request) {
	snippets, err := database.GetAllSnippets()
	if err != nil {
		http.Error(w, "Failed to get snippets", http.StatusInternalServerError)
		return
	}
	log.Println("Got all Snippets")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(snippets)
}

func createSnippets(w http.ResponseWriter, r *http.Request) {
	var requestSnippet models.Snippet

	if err := json.NewDecoder(r.Body).Decode(&requestSnippet); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}
	snippet, err := database.CreateSnippet(
		requestSnippet.Title,
		requestSnippet.Language,
		requestSnippet.Content,
	)
	if err != nil {
		http.Error(w, "failed to create snippet", http.StatusInternalServerError)
		log.Println(err)
		return
	}
	log.Println("Created snippet!")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(snippet)
}

func updateSnippetByID(w http.ResponseWriter, r *http.Request, snippetID uuid.UUID) {
	var requestSnippet models.Snippet
	if err := json.NewDecoder(r.Body).Decode(&requestSnippet); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	snippet, err := database.UpdateSnippet(
		requestSnippet.Title,
		requestSnippet.Language,
		requestSnippet.Content,
		snippetID,
	)
	if err != nil {
		http.Error(w, "failed to update snippet", http.StatusInternalServerError)
		log.Println(err)
		return
	}

	log.Println("Updated snippet!")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(snippet)
}

func getSnippetByID(w http.ResponseWriter, r *http.Request, snippetID uuid.UUID) {
	snippet, err := database.GetSnippetByID(snippetID)
	if err != nil {
		http.Error(w, "Failed to get snippet", http.StatusInternalServerError)
		log.Println(err)
	}
	log.Println("Snippet fetched from ID")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(snippet)
}
