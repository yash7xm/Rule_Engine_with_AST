package Test

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/joho/godotenv"
	"github.com/yash7xm/Rule_Engine_with_AST/cmd/routes"
	db "github.com/yash7xm/Rule_Engine_with_AST/internal/database"
)

// Helper function to create a new HTTP request with JSON body
func newJSONRequest(t *testing.T, method, url string, body interface{}) *http.Request {
	jsonBody, err := json.Marshal(body)
	if err != nil {
		t.Fatalf("Failed to marshal JSON body: %v", err)
	}
	req, err := http.NewRequest(method, url, bytes.NewBuffer(jsonBody))
	if err != nil {
		t.Fatalf("Failed to create HTTP request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	return req
}

// Test for createRuleHandler
func TestCreateRuleHandler(t *testing.T) {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	// Initialize the database connection
	db.InitDB()
	defer db.DB.Close() // Close the DB connection when the application shuts down

	// Mock request data
	reqBody := map[string]string{"rule_string": "a = 1"}

	// Create a new request
	req := newJSONRequest(t, "POST", "/create_rule", reqBody)

	// Create a new response recorder
	rr := httptest.NewRecorder()

	// Call the handler
	handler := http.HandlerFunc(routes.CreateRuleHandler)
	handler.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, status)
	}

	// Check the response body for valid JSON
	var response map[string]interface{}
	err = json.NewDecoder(rr.Body).Decode(&response)
	if err != nil {
		t.Errorf("Expected valid JSON response, got error: %v", err)
	}

	// Check if rule_id exists in the response
	if _, ok := response["data"]; !ok {
		t.Errorf("Expected rule_id in response, but got %v", response["data"])
	}
}

// Test for combineRulesHandler
func TestCombineRulesHandler(t *testing.T) {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	// Initialize the database connection
	db.InitDB()
	defer db.DB.Close() // Close the DB connection when the application shuts down

	// Mock request data
	reqBody := map[string]interface{}{
		"rules": []string{"a = 1", "b = 2"},
	}

	// Create a new request
	req := newJSONRequest(t, "POST", "/combine_rules", reqBody)

	// Create a new response recorder
	rr := httptest.NewRecorder()

	// Call the handler
	handler := http.HandlerFunc(routes.CombineRulesHandler)
	handler.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, status)
	}

	// Check the response body for valid JSON
	var response map[string]interface{}
	err = json.NewDecoder(rr.Body).Decode(&response)
	if err != nil {
		t.Errorf("Expected valid JSON response, got error: %v", err)
	}

	// Ensure the combined rule contains the expected "OR" operator
	combinedRule, ok := response["data"].(map[string]interface{})["combined_rule"].(string)
	if !ok || !strings.Contains(combinedRule, "OR") {
		t.Errorf("Expected combined rule with OR operator, got %v", response)
	}
}

// Test for evaluateRuleHandler
func TestEvaluateRuleHandler(t *testing.T) {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	// Initialize the database connection
	db.InitDB()
	defer db.DB.Close() // Close the DB connection when the application shuts down

	// Mock request data, representing AST as a JSON-like map structure
	reqBody := `{
		"ast": {
		  "Type": "LogicalOrExpression",
		  "Value": "OR",
		  "Left": {
			"Type": "LogicalAndExpression",
			"Value": "AND",
			"Left": {
			  "Type": "BinaryExpression",
			  "Value": ">",
			  "Left": { "Type": "Identifier", "Value": "age" },
			  "Right": { "Type": "NumericLiteral", "Value": "20" }
			},
			"Right": {
			  "Type": "BinaryExpression",
			  "Value": "<",
			  "Left": { "Type": "Identifier", "Value": "age" },
			  "Right": { "Type": "NumericLiteral", "Value": "30" }
			}
		  },
		  "Right": {
			"Type": "BinaryExpression",
			"Value": "=",
			"Left": { "Type": "Identifier", "Value": "status" },
			"Right": { "Type": "StringLiteral", "Value": "active" }
		  }
		},
		"data": { "age": 25, "status": "active" }
	  }`

	// Create a new request with the mocked request body
	req := httptest.NewRequest("POST", "/evaluate_rule", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")

	// Create a new response recorder to capture the handler's response
	rr := httptest.NewRecorder()

	// Call the handler (adjust the handler call as per your routes)
	handler := http.HandlerFunc(routes.EvaluateRuleHandler)
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect (200 OK)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, status)
	}

	// Check the response body for valid JSON
	var response map[string]interface{}
	err = json.NewDecoder(rr.Body).Decode(&response)
	if err != nil {
		t.Errorf("Expected valid JSON response, got error: %v", err)
	}

	// Check if the result is true, which means the rule was evaluated correctly
	if result, ok := response["data"].(map[string]interface{})["result"].(bool); !ok || !result {
		t.Errorf("Expected result true, got %v", response["result"])
	}
}
