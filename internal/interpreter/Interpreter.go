package interpreter

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/yash7xm/Rule_Engine_with_AST/internal/parser"
)

// Context is a map to hold input values to be evaluated against the rule.
type Context map[string]interface{}

// Interpreter evaluates the AST based on the given context data.
func Interpret(node *parser.Node, context Context) bool {
	if node == nil {
		return false
	}

	switch node.Type {
	case "LogicalAndExpression":
		// AND: Both left and right must be true
		return Interpret(node.Left, context) && Interpret(node.Right, context)
	case "LogicalOrExpression":
		// OR: Either left or right must be true
		return Interpret(node.Left, context) || Interpret(node.Right, context)
	case "BinaryExpression":
		// Binary expressions: Comparison like =, >, <, etc.
		return evaluateBinaryExpression(node, context)
	case "Identifier":
		// Lookup identifier value from the context
		value := context[node.Value]
		return value != nil
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
	var leftValue, rightValue interface{}
	if node.Left != nil {
		leftValue = evaluateExpression(node.Left, context)
	}

	if node.Right != nil {
		rightValue = evaluateExpression(node.Right, context)
	}

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
			fmt.Printf("Error converting NumericLiteral '%s' to int: %v\n", node.Value, err)
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
