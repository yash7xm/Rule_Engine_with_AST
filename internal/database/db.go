package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/golang-migrate/migrate/source/file"
	_ "github.com/lib/pq"
)

var DB *sql.DB

// Initialize the database connection
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

	log.Println("Connected to the database successfully!")
}

// Run database migrations
func RunMigrations() error {
	migrationQuery := `
	CREATE TABLE IF NOT EXISTS rules (
		id SERIAL PRIMARY KEY,
		rule_string TEXT NOT NULL,
		ast JSONB NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);
	`

	_, err := DB.Exec(migrationQuery)
	if err != nil {
		log.Fatalf("Error running migrations: %v", err)
		return err
	}

	fmt.Println("Database migrations completed successfully!")
	return nil
}

// InsertRule inserts a new rule into the "rules" table in the database.
func InsertRule(rule string) error {
	query := `INSERT INTO rules (rule_string) VALUES ($1)`
	_, err := DB.Exec(query, rule)
	if err != nil {
		log.Printf("Error inserting rule: %v", err)
		return err
	}

	log.Printf("Successfully inserted rule: %s", rule)
	return nil
}

// GetRules retrieves all rules from the database
func GetRules() ([]string, error) {
	rows, err := DB.Query("SELECT rule_string FROM rules")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rules []string
	for rows.Next() {
		var rule string
		if err := rows.Scan(&rule); err != nil {
			return nil, err
		}
		rules = append(rules, rule)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return rules, nil
}
