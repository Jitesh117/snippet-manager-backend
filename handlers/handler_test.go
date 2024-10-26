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

var jwtTokenString, snippetID string

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
	var response struct {
		Token string `json:"token"`
	}

	err = json.NewDecoder(rr.Body).Decode(&response)
	if err != nil {
		t.Fatalf("failed to parse response body: %v", err)
	}

	jwtTokenString = response.Token
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
	var responseSnippet models.Snippet
	err = json.NewDecoder(rr.Body).Decode(&responseSnippet)
	if err != nil {
		t.Fatalf("failed to parse response body: %v", err)
	}

	snippetID = responseSnippet.SnippetId.String()
}

func TestGetSnippetByID(t *testing.T) {
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
