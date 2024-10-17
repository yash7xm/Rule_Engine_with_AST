package routes

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/yash7xm/Rule_Engine_with_AST/internal/interpreter"
	"github.com/yash7xm/Rule_Engine_with_AST/internal/parser"
	"github.com/yash7xm/Rule_Engine_with_AST/internal/utils"
)

// NewRouter creates a new HTTP router and registers routes.
func NewRouter() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/create_rule", createRuleHandler)
	mux.HandleFunc("/combine_rules", combineRulesHandler)
	mux.HandleFunc("/evaluate_rule", evaluateRuleHandler)
	return mux
}

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
		http.Error(w, "Error creating AST: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Serialize the AST parser.Node to JSON and send it as a response
	response, err := json.Marshal(ast)
	if err != nil {
		http.Error(w, "Error marshaling AST to JSON", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

// API to combine rules
func combineRulesHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Rules []string `json:"rules"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Use combineAST to combine rules into a single AST
	combinedAST, err := combineAST(req.Rules)
	if err != nil {
		http.Error(w, "Error combining rules: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Serialize the combined AST node to JSON and send it as a response
	response, err := json.Marshal(combinedAST)
	if err != nil {
		http.Error(w, "Error marshaling combined AST to JSON", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

// API to evaluate a rule
func evaluateRuleHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		AST  map[string]interface{} `json:"ast"`
		Data map[string]interface{} `json:"data"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Convert the incoming AST JSON to your internal AST structure
	ast, err := utils.ConvertToASTNode(req.AST)
	if err != nil {
		http.Error(w, "Error converting AST: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Assume evaluateAST is your function that evaluates the AST against the provided data
	result := evaluateAST(ast, req.Data)

	// Send the result as a response
	response := map[string]interface{}{"result": result, "ast": ast}
	responseData, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Error marshaling result to JSON", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(responseData)
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

// combineAST combines multiple rule strings into a single AST using the OR operator.
func combineAST(rules []string) (*parser.Node, error) {
	if len(rules) == 0 {
		return nil, fmt.Errorf("no rules provided")
	}

	// Concatenate rules with OR operator
	combinedRule := rules[0]
	for _, rule := range rules[1:] {
		combinedRule += " OR " + rule
	}

	// Create AST from the combined rule
	ast, err := createAST(combinedRule)
	if err != nil {
		return nil, err
	}

	return ast, nil
}

// evaluateAST evaluates the given AST node using the provided context.
func evaluateAST(ast *parser.Node, data map[string]interface{}) bool {
	context := interpreter.Context(data)
	result := interpreter.Interpret(ast, context)

	return result
}
