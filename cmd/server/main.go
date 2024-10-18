package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	"github.com/rs/cors"
	"github.com/yash7xm/Rule_Engine_with_AST/cmd/routes"
	db "github.com/yash7xm/Rule_Engine_with_AST/internal/database"
)

func main() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	// Initialize the database connection
	db.InitDB()
	defer db.DB.Close() // Close the DB connection when the application shuts down

	// Set up the server with routes
	router := routes.NewRouter()

	// Enable CORS
	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders:   []string{"Content-Type"},
		AllowCredentials: true,
	})

	// Wrap the router with the CORS handler
	handler := corsHandler.Handler(router)

	// Capture system signals to handle graceful shutdown
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, syscall.SIGINT, syscall.SIGTERM)

	port := ":8080"
	fmt.Printf("Server is running on port %s...\n", port)

	go func() {
		if err := http.ListenAndServe(port, handler); err != nil {
			log.Fatalf("Error starting server: %s\n", err)
		}
	}()

	<-stopChan // Wait for termination signal
	fmt.Println("Shutting down gracefully...")
}
