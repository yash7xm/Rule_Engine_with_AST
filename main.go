package main

import (
	"fmt"
	"github.com/yash7xm/Rule_Engine_with_AST/parser"
	"strconv"
	"strings"
)

// Context is a map to hold input values to be evaluated against the rule.
type Context map[string]interface{}

// Interpreter evaluates the AST based on the given context data.
func interpret(node *parser.Node, context Context) bool {
	if node == nil {
		return false
	}

	switch node.Type {
	case "LogicalAndExpression":
		// AND: Both left and right must be true
		return interpret(node.Left, context) && interpret(node.Right, context)
	case "LogicalOrExpression":
		// OR: Either left or right must be true
		return interpret(node.Left, context) || interpret(node.Right, context)
	case "BinaryExpression":
		// Binary expressions: Comparison like =, >, <, etc.
		return evaluateBinaryExpression(node, context)
	case "Identifier":
		// Lookup identifier value from the context
		return context[node.Value] != nil
	case "NumericLiteral":
		// Numeric literals will just return their value
		return context[node.Value] != nil
	case "StringLiteral":
		// String literals will be evaluated as strings
		return context[node.Value] != nil
	default:
		fmt.Printf("Unknown node type: %s\n", node.Type)
	}

	return false
}

// Helper function to evaluate binary expressions like =, >, <, <=, >= etc.
func evaluateBinaryExpression(node *parser.Node, context Context) bool {
	leftValue := evaluateExpression(node.Left, context)
	rightValue := evaluateExpression(node.Right, context)

	if leftValue == nil || rightValue == nil {
		return false
	}

	switch node.Value {
	case "=":
		// Handle equality comparison
		return leftValue == rightValue
	case ">":
		// Handle greater than comparison
		leftNum, leftOk := toNumber(leftValue)
		rightNum, rightOk := toNumber(rightValue)
		return leftOk && rightOk && leftNum > rightNum
	case "<":
		// Handle less than comparison
		leftNum, leftOk := toNumber(leftValue)
		rightNum, rightOk := toNumber(rightValue)
		return leftOk && rightOk && leftNum < rightNum
	case ">=":
		// Handle greater than or equal to comparison
		leftNum, leftOk := toNumber(leftValue)
		rightNum, rightOk := toNumber(rightValue)
		return leftOk && rightOk && leftNum >= rightNum
	case "<=":
		// Handle less than or equal to comparison
		leftNum, leftOk := toNumber(leftValue)
		rightNum, rightOk := toNumber(rightValue)
		return leftOk && rightOk && leftNum <= rightNum
	case "!=":
		// Handle not equal to comparison
		return leftValue != rightValue
	default:
		fmt.Printf("Unknown operator: %s\n", node.Value)
		return false
	}
}

// Helper function to evaluate expressions and return their values
func evaluateExpression(node *parser.Node, context Context) interface{} {
	switch node.Type {
	case "Identifier":
		// Get the value of an identifier from the context
		return context[node.Value]
	case "NumericLiteral":
		// Convert string numeric literals to a number
		num, err := strconv.Atoi(node.Value)
		if err != nil {
			return nil
		}
		return num
	case "StringLiteral":
		// Strip quotes from string literals
		return strings.Trim(node.Value, "'")
	}

	return nil
}

// Convert a value to a number if possible
func toNumber(value interface{}) (int, bool) {
	switch v := value.(type) {
	case int:
		return v, true
	case string:
		num, err := strconv.Atoi(v)
		return num, err == nil
	}

	return 0, false
}

// Helper function to print the entire AST recursively
func printAST(node *parser.Node, depth int) {
	if node == nil {
		return
	}

	// Indentation to visualize the tree structure
	indent := ""
	for i := 0; i < depth; i++ {
		indent += "  "
	}

	// Print current node
	fmt.Printf("%sNode Type: %s, Value: %s\n", indent, node.Type, node.Value)

	// Recursively print left and right child nodes
	if node.Left != nil {
		fmt.Printf("%sLeft:\n", indent)
		printAST(node.Left, depth+1)
	}
	if node.Right != nil {
		fmt.Printf("%sRight:\n", indent)
		printAST(node.Right, depth+1)
	}
}

func main() {
	// Define the rule
	rule := "((age > 30 AND department = 'Sales') OR (age < 25 AND department = 'Marketing')) AND (salary > 50000 OR experience > 5)"

	// Tokenizer and parser for the rule
	tokenizer := parser.NewTokenizer(rule)
	p := parser.NewParser(tokenizer)

	// Parse the rule into an AST
	ast := p.ParseRule()

	// Print the entire AST structure
	fmt.Println("AST Structure:")
	printAST(ast, 0)

	// Example context data to evaluate the rule
	context := Context{
		"age":        32,
		"department": "Sales",
		"salary":     60000,
		"experience": 3,
	}

	// Interpret the AST against the context data
	result := interpret(ast, context)

	fmt.Printf("Result of rule evaluation: %v\n", result)
}
