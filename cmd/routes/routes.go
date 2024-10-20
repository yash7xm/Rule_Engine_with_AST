package routes

import (
	"encoding/json"
	"fmt"
	"net/http"

	db "github.com/yash7xm/Rule_Engine_with_AST/internal/database"
	"github.com/yash7xm/Rule_Engine_with_AST/internal/interpreter"
	"github.com/yash7xm/Rule_Engine_with_AST/internal/parser"
	"github.com/yash7xm/Rule_Engine_with_AST/internal/utils"
)

// NewRouter creates a new HTTP router and registers routes.
func NewRouter() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/create_rule", CreateRuleHandler)
	mux.HandleFunc("/combine_rules", CombineRulesHandler)
	mux.HandleFunc("/evaluate_rule", EvaluateRuleHandler)
	mux.HandleFunc("/ping", PongHandler)
	return mux
}

// CreateRuleHandler handles the creation of a rule and stores it in the database.
func CreateRuleHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		RuleString string `json:"rule_string"`
	}

	// Decode the incoming request JSON body
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		SendErrorResponse(w, http.StatusBadRequest, "Invalid request payload", err)
		return
	}

	// Create AST (Abstract Syntax Tree) from the rule string
	ast, err := createAST(req.RuleString)
	if err != nil {
		SendErrorResponse(w, http.StatusBadRequest, "Error creating AST", err)
		return
	}

	// Convert AST to JSON format to store in the database
	astJSON, err := json.Marshal(ast)
	if err != nil {
		SendErrorResponse(w, http.StatusInternalServerError, "Error marshaling AST to JSON", err)
		return
	}

	// Insert the rule and the AST into the database
	query := `INSERT INTO rules (rule_string, ast) VALUES ($1, $2) RETURNING id`
	var ruleID int
	err = db.DB.QueryRow(query, req.RuleString, astJSON).Scan(&ruleID)
	if err != nil {
		SendErrorResponse(w, http.StatusInternalServerError, "Error storing rule", err)
		return
	}

	// Prepare response data
	responseData := map[string]interface{}{
		"rule_id": ruleID,
		"node":    ast,
	}

	// Send success response
	SendSuccessResponse(w, http.StatusOK, "Rule created successfully", responseData)
}

// CombineRulesHandler combines multiple rules into one and stores the result in the database.
func CombineRulesHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Rules []string `json:"rules"`
	}

	// Decode the request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		SendErrorResponse(w, http.StatusBadRequest, "Invalid request payload", err)
		return
	}

	// Use combineAST to combine rules into a single AST
	combinedAST, err := combineAST(req.Rules)
	if err != nil {
		SendErrorResponse(w, http.StatusInternalServerError, "Error combining rules", err)
		return
	}

	// Convert combined AST to JSON format
	astJSON, err := json.Marshal(combinedAST)
	if err != nil {
		SendErrorResponse(w, http.StatusInternalServerError, "Error marshaling combined AST to JSON", err)
		return
	}

	// Create the combined rule string (concatenating the rules)
	combinedRuleString := req.Rules[0]
	for _, rule := range req.Rules[1:] {
		combinedRuleString += " OR " + rule
	}

	// Insert the combined rule and the AST into the database
	query := `INSERT INTO rules (rule_string, ast) VALUES ($1, $2) RETURNING id`
	var ruleID int
	err = db.DB.QueryRow(query, combinedRuleString, astJSON).Scan(&ruleID)
	if err != nil {
		SendErrorResponse(w, http.StatusInternalServerError, "Error storing combined rule", err)
		return
	}

	// Prepare response data
	responseData := map[string]interface{}{
		"rule_id":       ruleID,
		"combined_rule": combinedRuleString,
		"node":          combinedAST,
	}

	// Send success response
	SendSuccessResponse(w, http.StatusOK, "Combined rule created successfully", responseData)
}

// EvaluateRuleHandler evaluates a rule's AST against provided data.
func EvaluateRuleHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		AST  map[string]interface{} `json:"ast"`
		Data map[string]interface{} `json:"data"`
	}

	// Decode the request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		SendErrorResponse(w, http.StatusBadRequest, "Invalid request payload", err)
		return
	}

	// Convert the incoming AST JSON to your internal AST structure
	ast, err := utils.ConvertToASTNode(req.AST)
	if err != nil {
		SendErrorResponse(w, http.StatusBadRequest, "Error converting AST", err)
		return
	}

	// Evaluate the AST with the provided data
	result := evaluateAST(ast, req.Data)

	// Prepare response data
	responseData := map[string]interface{}{
		"result": result,
	}

	// Send success response
	SendSuccessResponse(w, http.StatusOK, "Rule evaluated successfully", responseData)
}

func PongHandler(w http.ResponseWriter, r *http.Request) {
	responseData := map[string]interface{}{
		"result": "Pong",
	}

	// Send success response
	SendSuccessResponse(w, http.StatusOK, "Rule evaluated successfully", responseData)
}

// Placeholder functions for AST operations

// createAST creates an AST from a rule string.
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
