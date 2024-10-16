package main

import (
	"log"
	"net/http"

	"github.com/Jitesh117/snippet-manager/database"
	"github.com/Jitesh117/snippet-manager/handlers"
)

func main() {
	database.InitDB()
	defer database.CloseDB()

	http.HandleFunc("/snippets", handlers.HandleSnippets)

	log.Println("Server is running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
