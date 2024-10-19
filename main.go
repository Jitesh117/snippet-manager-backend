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

	// Protected endpoints with rate limiter and JWT middleware
	http.Handle(
		"/snippets",
		auth.RateLimiter(auth.JWTAuthMiddleware(http.HandlerFunc(handlers.HandleSnippets))),
	)
	http.Handle(
		"/snippets/",
		auth.RateLimiter(auth.JWTAuthMiddleware(http.HandlerFunc(handlers.HandleSnippet))),
	)

	// Open endpoints with just rate limiter
	http.HandleFunc("/register", auth.RateLimiter(handlers.Register))
	http.HandleFunc("/login", auth.RateLimiter(handlers.Login))
	http.HandleFunc("/deleteUser", auth.RateLimiter(handlers.DeleteUserByID))

	log.Println("Server is running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
