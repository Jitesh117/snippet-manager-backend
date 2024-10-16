package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Jitesh117/snippet-manager/database"
	"github.com/Jitesh117/snippet-manager/models"
)

func HandleSnippets(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getAllSnippets(w, r)
	case http.MethodPost:
		createSnippets(w, r)
	}
}

func getAllSnippets(w http.ResponseWriter, r *http.Request) {
	snippets, err := database.GetAllSnippets()
	if err != nil {
		http.Error(w, "Failed to get snippets", http.StatusInternalServerError)
		return
	}
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
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(snippet)
}
