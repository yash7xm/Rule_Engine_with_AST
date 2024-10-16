package main

import (
	"fmt"
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
	rule := "((age > 30 AND department = 'Sales') OR (age < 25 AND department = 'Marketing')) AND (salary > 50000 OR experience > 5)"
	tokenizer := parser.NewTokenizer(rule)
	p := parser.NewParser(tokenizer)

	// Parse the rule into an AST
	ast := p.ParseRule()

	// Print the entire AST starting from the root node
	fmt.Println("AST Structure:")
	printAST(ast, 0)
}
