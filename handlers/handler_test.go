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

var jwtTokenString = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3Mjk0Mzc1NjgsInVzZXJfaWQiOiI3Y2MzOWUxYy0xYzc1LTRjYWUtOTI3Mi0wNDY3NmMxMmY1YzUifQ.ZU55nWTWBYPZXme3H7WiNJm2zTuYkYt4v0WKkNGcyE4"

func TestMain(m *testing.M) {
	database.InitDB()
	m.Run()
	database.CloseDB()
}

func TestRegister(t *testing.T) {
	user := models.User{
		UserName: "tester",
		Email:    "testing@test.com",
		Password: "password_test",
	}
	body, _ := json.Marshal(user)
	req, err := http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handlers.Register)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}

func TestLogin(t *testing.T) {
	loginData := map[string]string{
		"user_name": "tester",
		"email":     "testing@test.com",
		"password":  "password_test",
	}
	body, _ := json.Marshal(loginData)
	req, err := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handlers.Login)

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
	snippetID := "3e2272fa-e415-4d37-86c2-228f2353c152"
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
	snippetID := "3e2272fa-e415-4d37-86c2-228f2353c152"
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
		"email":    "testing@test.com",
		"password": "password_test",
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
