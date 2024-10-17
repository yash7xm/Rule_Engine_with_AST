package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/yash7xm/Rule_Engine_with_AST/parser"
)

// API to create a rule
func createRuleHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		RuleString string `json:"rule_string"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Create AST and handle any errors
	ast, err := createAST(req.RuleString)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error creating AST: %s", err), http.StatusBadRequest)
		return
	}

	// Serialize the AST parser.Node to JSON and send it as a response
	response, _ := json.Marshal(ast)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

// Start the HTTP server
func main() {
	http.HandleFunc("/create_rule", createRuleHandler)

	port := ":8080"
	fmt.Printf("Server is running on port %s...\n", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		fmt.Printf("Error starting server: %s\n", err)
	}
}

// Placeholder functions for AST operations
func createAST(rule string) (*parser.Node, error) {
	tokenizer := parser.NewTokenizer(rule)
	p := parser.NewParser(tokenizer)
	ast, err := p.ParseRule()
	if err != nil {
		return nil, err
	}
	return ast, nil
}
