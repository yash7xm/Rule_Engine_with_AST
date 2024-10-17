package main

import (
	"fmt"

	"github.com/yash7xm/Rule_Engine_with_AST/interpreter"
	"github.com/yash7xm/Rule_Engine_with_AST/parser"
)

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
	// rule := "((age > 30 AND department = 'Sales') OR (age < 25 AND department = 'Marketing')) AND (salary > 50000 OR experience > 5)"

	rule := "age > "
	// Tokenizer and parser for the rule
	tokenizer := parser.NewTokenizer(rule)
	p := parser.NewParser(tokenizer)

	// Parse the rule into an AST
	ast := p.ParseRule()

	// Print the entire AST structure
	fmt.Println("AST Structure:")
	printAST(ast, 0)

	// Example context data to evaluate the rule
	context := interpreter.Context{
		"age":        32,
		"department": "Sales",
		"salary":     60000,
		"experience": 3,
	}

	// Interpret the AST against the context data
	result := interpreter.Interpret(ast, context)

	fmt.Printf("Result of rule evaluation: %v\n", result)
}
