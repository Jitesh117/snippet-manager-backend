package main

import (
	"log"
	"net/http"

	"github.com/Jitesh117/snippet-manager-backend/database"
	"github.com/Jitesh117/snippet-manager-backend/handlers"
	auth "github.com/Jitesh117/snippet-manager-backend/middleware"
)

func main() {
	database.InitDB()
	defer database.CloseDB()

	// protected endpoints
	http.HandleFunc("/snippets", auth.JWTAuthMiddleware(handlers.HandleSnippets))
	http.HandleFunc("/snippets/", auth.JWTAuthMiddleware(handlers.HandleSnippet))

	// open endpoints
	http.HandleFunc("/register", handlers.Register)
	http.HandleFunc("/login", handlers.Login)

	log.Println("Server is running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
