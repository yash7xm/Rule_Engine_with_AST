package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/joho/godotenv"
	"github.com/yash7xm/Rule_Engine_with_AST/cmd/routes"
	db "github.com/yash7xm/Rule_Engine_with_AST/internal/database"
)

func main() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	db.InitDB()

	// Set up the server with routes
	router := routes.NewRouter()

	port := ":8080"
	fmt.Printf("Server is running on port %s...\n", port)
	if err := http.ListenAndServe(port, router); err != nil {
		log.Fatalf("Error starting server: %s\n", err)
	}
}
