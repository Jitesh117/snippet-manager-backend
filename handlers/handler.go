package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"
	"unicode"

	"github.com/Jitesh117/snippet-manager-backend/database"
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

func createSnippet(w http.ResponseWriter, r *http.Request) {
	var requestSnippet models.Snippet
	var userID uuid.UUID

	if err := json.NewDecoder(r.Body).Decode(&requestSnippet); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	if err := validateSnippet(requestSnippet); err != nil {
		http.Error(w, "Invalid request payload: "+err.Error(), http.StatusBadRequest)
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

func validateSnippet(snippet models.Snippet) error {
	if snippet.Title == "" {
		return fmt.Errorf("title can't be empty!")
	}
	if snippet.Language == "" {
		return fmt.Errorf("language can't be empty!")
	}
	if snippet.Content == "" {
		return fmt.Errorf("content can't be empty!")
	}
	return nil
}

func updateSnippetByID(w http.ResponseWriter, r *http.Request, snippetID uuid.UUID) {
	var requestSnippet models.Snippet
	if err := json.NewDecoder(r.Body).Decode(&requestSnippet); err != nil {
		http.Error(w, "Invalid request payload: ", http.StatusBadRequest)
		return
	}
	if err := validateSnippet(requestSnippet); err != nil {
		http.Error(w, "Invalid request payload: "+err.Error(), http.StatusBadRequest)
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

func validateEmail(email string) bool {
	emailPattern := `^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`
	return regexp.MustCompile(emailPattern).MatchString(email)
}

func validatePassword(password string) error {
	if len(password) < 8 {
		return fmt.Errorf("password must be at least 8 characters long")
	}

	var (
		hasUpper   bool
		hasLower   bool
		hasNumber  bool
		hasSpecial bool
	)

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		case strings.ContainsRune(`!@#$%^&*(),.?":{}|<>`, char):
			hasSpecial = true
		}
	}

	var missing []string
	if !hasUpper {
		missing = append(missing, "uppercase letter")
	}
	if !hasLower {
		missing = append(missing, "lowercase letter")
	}
	if !hasNumber {
		missing = append(missing, "number")
	}
	if !hasSpecial {
		missing = append(missing, "special character")
	}

	if len(missing) > 0 {
		return fmt.Errorf("password must contain at least one %s", strings.Join(missing, ", "))
	}

	return nil
}

func validateUser(user models.User) error {
	if user.UserName == "" {
		return fmt.Errorf("username can't be empty")
	}
	if len(user.UserName) < 3 {
		return fmt.Errorf("username must be at least 3 characters long")
	}

	if user.Email == "" {
		return fmt.Errorf("email can't be empty")
	}
	if !validateEmail(user.Email) {
		return fmt.Errorf("invalid email format")
	}

	if user.Password == "" {
		return fmt.Errorf("password can't be empty")
	}
	if err := validatePassword(user.Password); err != nil {
		return err
	}

	return nil
}

func Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	if err := validateUser(user); err != nil {
		http.Error(w, "Invalid request payload: "+err.Error(), http.StatusBadRequest)
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
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var loginData struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	err := json.NewDecoder(r.Body).Decode(&loginData)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	userID, err := database.CheckUserCredentials(loginData.Email, loginData.Password)
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}
	log.Println(userID)

	token, err := auth.GenerateJWT(userID)
	if err != nil {
		http.Error(w, "Failed to generated JWT: "+err.Error(), http.StatusInternalServerError)
		return
	}

	log.Println("user logged in!")
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

func DeleteUserByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var userData struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := json.NewDecoder(r.Body).Decode(&userData)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	userID, err := database.CheckUserCredentials(userData.Email, userData.Password)
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}
	deletedUserID, err := database.DeleteUser(userID)
	if err != nil {
		http.Error(w, "Failed to Delete user", http.StatusInternalServerError)
		return
	}

	log.Println("user deleted!")
	json.NewEncoder(w).Encode(map[string]string{"userID": deletedUserID.String()})
}
