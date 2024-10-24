package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Jitesh117/snippet-manager-backend/database"
	"github.com/Jitesh117/snippet-manager-backend/handlers"
	"github.com/Jitesh117/snippet-manager-backend/models"
)

var jwtTokenString = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3Mjk2NzE3NTEsInVzZXJfaWQiOiJmOTI3ZTgzNy1kNTVmLTQ1YjAtODM4Ni1mZjQ3NjQ0OGZjNTcifQ.wgKQAKM5u09MN1izKPKUGFrzM1LHzi_MJETnNTRREw4"

func TestMain(m *testing.M) {
	database.InitDB()
	m.Run()
	database.CloseDB()
}

func TestRegister(t *testing.T) {
	user := models.User{
		UserName: "testerTestNew",
		Email:    "testingTest@testNew.com",
		Password: "Password@123",
	}
	body, _ := json.Marshal(user)
	req, err := http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handlers.RegisterUser)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}

func TestLogin(t *testing.T) {
	loginData := models.User{
		UserName: "testerTestNew",
		Email:    "testingTest@testNew.com",
		Password: "Password@123",
	}
	body, _ := json.Marshal(loginData)
	req, err := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handlers.LoginUser)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}

func TestChangePassword(t *testing.T) {
	changeData := map[string]string{
		"email":        "testingTest@testNew.com",
		"password":     "Password@123",
		"new_password": "NewPassword@123",
	}

	body, _ := json.Marshal(changeData)

	req, err := http.NewRequest(http.MethodPut, "/changePassword", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handlers.ChangePassword)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}

func TestGetAllSnippets(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "/snippets", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Authorization", "Bearer "+jwtTokenString)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handlers.HandleSnippets)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}

func TestCreateSnippet(t *testing.T) {
	snippet := models.Snippet{
		Title:    "test snippet",
		Language: "Go",
		Content:  "fmt.Println('Hello tests!')",
	}

	body, _ := json.Marshal(snippet)
	req, err := http.NewRequest(http.MethodPost, "/snippets", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Authorization", "Bearer "+jwtTokenString)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handlers.HandleSnippets)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusCreated)
	}
}

func TestGetSnippetByID(t *testing.T) {
	snippetID := "2d8b9384-a551-4a09-96d4-f9c1678a778c"
	req, err := http.NewRequest(http.MethodGet, "/snippets/"+snippetID, nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Authorization", "Bearer "+jwtTokenString)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handlers.HandleSnippet)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}

func TestUpdateSnippetByID(t *testing.T) {
	snippetID := "2d8b9384-a551-4a09-96d4-f9c1678a778c"
	snippet := models.Snippet{
		Title:    "Updated Snippet",
		Language: "Go",
		Content:  "fmt.Println('Updated content')",
	}

	body, _ := json.Marshal(snippet)
	req, err := http.NewRequest(
		http.MethodPut,
		"/snippets/"+snippetID,
		bytes.NewBuffer(body),
	)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Authorization", "Bearer "+jwtTokenString)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handlers.HandleSnippet)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}

func TestDeleteUser(t *testing.T) {
	userData := map[string]string{
		"email":    "testingTest@testNew.com",
		"password": "NewPassword@123",
	}
	body, _ := json.Marshal(userData)

	req, err := http.NewRequest(http.MethodDelete, "/deleteUser", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handlers.DeleteUserByID)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}
