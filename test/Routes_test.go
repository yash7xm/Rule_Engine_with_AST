package Test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/yash7xm/Rule_Engine_with_AST/cmd/routes"
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
	err := json.NewDecoder(rr.Body).Decode(&response)
	if err != nil {
		t.Errorf("Expected valid JSON response, got error: %v", err)
	}
}

// Test for combineRulesHandler
func TestCombineRulesHandler(t *testing.T) {
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
	err := json.NewDecoder(rr.Body).Decode(&response)
	if err != nil {
		t.Errorf("Expected valid JSON response, got error: %v", err)
	}

	// Ensure the combined rule contains the expected "OR" operator
	combinedRule, ok := response["Value"].(string)
	if !ok || !strings.Contains(combinedRule, "OR") {
		t.Errorf("Expected combined rule with OR operator, got %v", combinedRule)
	}
}

// Test for evaluateRuleHandler
func TestEvaluateRuleHandler(t *testing.T) {
	// Mock request data
	reqBody := map[string]interface{}{
		"ast": map[string]interface{}{
			"Type":  "BinaryExpression",
			"Left":  map[string]interface{}{"Type": "Identifier", "Value": "a"},
			"Right": map[string]interface{}{"Type": "NumericLiteral", "Value": "1"},
			"Value": ">",
		},
		"data": map[string]interface{}{"a": "10"},
	}

	// Create a new request
	req := newJSONRequest(t, "POST", "/evaluate_rule", reqBody)

	// Create a new response recorder
	rr := httptest.NewRecorder()

	// Call the handler
	handler := http.HandlerFunc(routes.EvaluateRuleHandler)
	handler.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, status)
	}

	// Check the response body for valid JSON
	var response map[string]interface{}
	err := json.NewDecoder(rr.Body).Decode(&response)
	if err != nil {
		t.Errorf("Expected valid JSON response, got error: %v", err)
	}

	// Check if the result is true, which means the rule was evaluated correctly
	if result, ok := response["result"].(bool); !ok || !result {
		t.Errorf("Expected result true, got %v", response["result"])
	}
}
