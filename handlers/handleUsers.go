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
)

func RegisterUser(w http.ResponseWriter, r *http.Request) {
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
	if err := helper.ValidateUser(user); err != nil {
		http.Error(w, constants.ErrInvalidPayload+": "+err.Error(), http.StatusBadRequest)
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

func LoginUser(w http.ResponseWriter, r *http.Request) {
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
			constants.ErrFailedToGenerateToken+": "+err.Error(),
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

func ChangePassword(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, constants.ErrMethodNotAllowed, http.StatusMethodNotAllowed)
		return
	}

	var userData struct {
		Email       string `json:"email"`
		Password    string `json:"password"`
		NewPassword string `json:"new_password"`
	}

	err := json.NewDecoder(r.Body).Decode(&userData)
	if err != nil {
		http.Error(w, constants.ErrInvalidPayload, http.StatusBadRequest)
		return
	}

	if userData.Password == userData.NewPassword {
		http.Error(w, "New password must be different from the old password", http.StatusBadRequest)
		return
	}

	userID, err := database.CheckUserCredentials(userData.Email, userData.Password)
	if err != nil {
		http.Error(w, constants.ErrInvalidCredentials, http.StatusUnauthorized)
		return
	}
	err = helper.ValidatePassword(userData.NewPassword)
	if err != nil {
		http.Error(w, constants.ErrInvalidPayload+": "+err.Error(), http.StatusBadRequest)
		return
	}
	err = database.ChangePassword(userID, userData.NewPassword)
	if err != nil {
		http.Error(
			w,
			constants.ErrFailedToUpdatePassword+": "+err.Error(),
			http.StatusInternalServerError,
		)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Password updated successfully"))
}
