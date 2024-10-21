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

	"github.com/Jitesh117/snippet-manager-backend/constants"
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
		http.Error(w, constants.ErrMethodNotAllowed, http.StatusMethodNotAllowed)
	}
}

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
	if err := validateSnippet(requestSnippet); err != nil {
		http.Error(w, constants.ErrInvalidPayload+err.Error(), http.StatusBadRequest)
		return
	}

	userID, err := auth.ExtractUserIDFromToken(r)
	if err != nil {
		http.Error(w, constants.ErrFailedToExtractTokenID+err.Error(), http.StatusBadRequest)
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

func validateSnippet(snippet models.Snippet) error {
	if snippet.Title == "" {
		return fmt.Errorf(constants.ErrEmptyTitle)
	}
	if snippet.Language == "" {
		return fmt.Errorf(constants.ErrEmptyLanguage)
	}
	if snippet.Content == "" {
		return fmt.Errorf(constants.ErrEmptyContent)
	}
	return nil
}

func updateSnippetByID(w http.ResponseWriter, r *http.Request, snippetID uuid.UUID) {
	var requestSnippet models.Snippet
	if err := json.NewDecoder(r.Body).Decode(&requestSnippet); err != nil {
		http.Error(w, constants.ErrInvalidPayload, http.StatusBadRequest)
		return
	}
	if err := validateSnippet(requestSnippet); err != nil {
		http.Error(w, constants.ErrInvalidPayload+err.Error(), http.StatusBadRequest)
		return
	}
	userID, err := auth.ExtractUserIDFromToken(r)
	if err != nil {
		http.Error(w, constants.ErrFailedToGetUserID+err.Error(), http.StatusBadRequest)
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
		http.Error(w, constants.ErrFailedToGetUserID+err.Error(), http.StatusBadRequest)
		return
	}

	snippet, err := database.GetSnippetByID(snippetID, userID)
	if err != nil {
		http.Error(w, constants.ErrFailedToGetSnippets+err.Error(), http.StatusInternalServerError)
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
		http.Error(w, constants.ErrFailedToGetUserID+err.Error(), http.StatusBadRequest)
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

func validateEmail(email string) bool {
	emailPattern := `^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`
	return regexp.MustCompile(emailPattern).MatchString(email)
}

func validatePassword(password string) error {
	if len(password) < 8 {
		return fmt.Errorf(constants.ErrPasswordTooShort)
	}
	if len(password) > 20 {
		return fmt.Errorf(constants.ErrPasswordTooLong)
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
		return fmt.Errorf(constants.ErrEmptyUsername)
	}
	if len(user.UserName) < 3 {
		return fmt.Errorf(constants.ErrUsernameTooShort)
	}

	if len(user.UserName) > 30 {
		return fmt.Errorf(constants.ErrUsernameTooLong)
	}

	if user.Email == "" {
		return fmt.Errorf(constants.ErrEmptyContent)
	}
	if !validateEmail(user.Email) {
		return fmt.Errorf(constants.ErrInvalidEmailFormat)
	}

	if user.Password == "" {
		return fmt.Errorf(constants.ErrEmptyPassword)
	}
	if err := validatePassword(user.Password); err != nil {
		return err
	}

	return nil
}

func Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, constants.ErrMethodNotAllowed, http.StatusMethodNotAllowed)
		return
	}
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, constants.ErrInvalidPayload, http.StatusBadRequest)
		return
	}
	if err := validateUser(user); err != nil {
		http.Error(w, constants.ErrInvalidPayload+err.Error(), http.StatusBadRequest)
		return
	}

	userID, err := database.CreateUser(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	token, err := auth.GenerateJWT(userID)
	if err != nil {
		http.Error(w, constants.ErrFailedToGenerateToken, http.StatusInternalServerError)
		return
	}

	log.Println("userID: ", userID)
	log.Println("user signed in!")
	json.NewEncoder(w).Encode(token)
}

func Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, constants.ErrMethodNotAllowed, http.StatusMethodNotAllowed)
		return
	}
	var loginData struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	err := json.NewDecoder(r.Body).Decode(&loginData)
	if err != nil {
		http.Error(w, constants.ErrInvalidPayload, http.StatusBadRequest)
		return
	}

	userID, err := database.CheckUserCredentials(loginData.Email, loginData.Password)
	if err != nil {
		http.Error(w, constants.ErrInvalidCredentials, http.StatusUnauthorized)
		return
	}
	log.Println(userID)

	token, err := auth.GenerateJWT(userID)
	if err != nil {
		http.Error(
			w,
			constants.ErrFailedToGenerateToken+err.Error(),
			http.StatusInternalServerError,
		)
		return
	}

	log.Println("user logged in!")
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

func DeleteUserByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, constants.ErrMethodNotAllowed, http.StatusMethodNotAllowed)
		return
	}

	var userData struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := json.NewDecoder(r.Body).Decode(&userData)
	if err != nil {
		http.Error(w, constants.ErrInvalidPayload, http.StatusBadRequest)
		return
	}

	userID, err := database.CheckUserCredentials(userData.Email, userData.Password)
	if err != nil {
		http.Error(w, constants.ErrInvalidCredentials, http.StatusUnauthorized)
		return
	}
	deletedUserID, err := database.DeleteUser(userID)
	if err != nil {
		http.Error(w, constants.ErrFailedToDeleteUser, http.StatusInternalServerError)
		return
	}

	log.Println("user deleted!")
	json.NewEncoder(w).Encode(map[string]string{"userID": deletedUserID.String()})
}
