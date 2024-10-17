package main

import (
	"log"
	"net/http"

	"github.com/Jitesh117/snippet-manager-backend/database"
	"github.com/Jitesh117/snippet-manager-backend/handlers"
)

func main() {
	database.InitDB()
	defer database.CloseDB()

	http.HandleFunc("/snippets", handlers.HandleSnippets)
	http.HandleFunc("/snippets/", handlers.HandleSnippet)
	http.HandleFunc("/register", handlers.Register)
	http.HandleFunc("/login", handlers.Login)

	log.Println("Server is running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
