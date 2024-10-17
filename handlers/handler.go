package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/Jitesh117/snippet-manager/database"
	auth "github.com/Jitesh117/snippet-manager/middleware"
	"github.com/Jitesh117/snippet-manager/models"
	"github.com/google/uuid"
)

func HandleSnippets(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getAllSnippets(w)
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
	case http.MethodDelete:
		deleteSnippetByID(w, r, snippetID)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func getAllSnippets(w http.ResponseWriter) {
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
	var userID uuid.UUID

	if err := json.NewDecoder(r.Body).Decode(&requestSnippet); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	userID, err := auth.ExtractUserIDFromToken(r)
	if err != nil {
		http.Error(w, "failed to extract ID from token: "+err.Error(), http.StatusBadRequest)
		return
	}

	snippet, err := database.CreateSnippet(
		requestSnippet.Title,
		requestSnippet.Language,
		requestSnippet.Content,
		userID,
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
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	userID, err := auth.ExtractUserIDFromToken(r)
	if err != nil {
		http.Error(w, "failed to get userID: "+err.Error(), http.StatusBadRequest)
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
		http.Error(w, "failed to update snippet", http.StatusInternalServerError)
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
		http.Error(w, "Failed to get userID"+err.Error(), http.StatusBadRequest)
		return
	}

	snippet, err := database.GetSnippetByID(snippetID, userID)
	if err != nil {
		http.Error(w, "Failed to get snippet"+err.Error(), http.StatusInternalServerError)
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
		http.Error(w, "Failed to get userID: "+err.Error(), http.StatusBadRequest)
		return
	}

	snippet, err := database.DeleteSnippetByID(snippetID, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "snippet not found", http.StatusBadRequest)
		}
		http.Error(w, "Failed to delete snippet", http.StatusInternalServerError)
	}
	log.Println("snippet deleted from DB")
	json.NewEncoder(w).Encode(snippet)
}

func Register(w http.ResponseWriter, r *http.Request) {
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	userID, err := database.CreateUser(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	token, err := auth.GenerateJWT(userID)
	if err != nil {
		http.Error(w, "Failed to create jwt token", http.StatusInternalServerError)
		return
	}

	log.Println("userID: ", userID)
	log.Println("user signed in!")
	json.NewEncoder(w).Encode(token)
}

func Login(w http.ResponseWriter, r *http.Request) {
	var loginData struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	err := json.NewDecoder(r.Body).Decode(&loginData)
	if err != nil {
		http.Error(w, "invalid request payload", http.StatusBadRequest)
		return
	}

	userID, err := database.CheckUserCredentials(loginData.Email, loginData.Password)
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	token, err := auth.GenerateJWT(userID)
	if err != nil {
		http.Error(w, "Failed to generated JWT", http.StatusInternalServerError)
		return
	}

	log.Println("user logged in!")
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}
