package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func InitDB() {
	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		log.Fatalf("DATABASE_URL is not set")
	}

	var err error
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}

	// Test the database connection
	if err := DB.Ping(); err != nil {
		log.Fatalf("Error pinging database: %v", err)
	}

	fmt.Println("Connected to the database successfully!")
}
