package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/rs/cors"
	"github.com/yash7xm/Rule_Engine_with_AST/cmd/routes"
	db "github.com/yash7xm/Rule_Engine_with_AST/internal/database"
)

func main() {
	// Load environment variables from .env file
	// err := godotenv.Load()
	// if err != nil {
	// 	log.Fatalf("Error loading .env file")
	// }

	// Initialize the database connection
	db.InitDB()
	defer db.DB.Close() // Close the DB connection when the application shuts down

	// Run database migrations
	err := db.RunMigrations()
	if err != nil {
		log.Fatalf("Error running migrations: %v", err)
	}

	// Set up the server with routes
	router := routes.NewRouter()

	// Enable CORS
	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173", "https://rule-engine-ui.vercel.app"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders:   []string{"Content-Type"},
		AllowCredentials: true,
	})

	// Wrap the router with the CORS handler
	handler := corsHandler.Handler(router)

	// Capture system signals to handle graceful shutdown
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, syscall.SIGINT, syscall.SIGTERM)

	// Use the PORT environment variable
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default port if PORT is not set (for local development)
	}
	fmt.Printf("Server is running on port %s...\n", port)

	go func() {
		if err := http.ListenAndServe(":"+port, handler); err != nil {
			log.Fatalf("Error starting server: %s\n", err)
		}
	}()

	<-stopChan // Wait for termination signal
	fmt.Println("Shutting down gracefully...")
}
