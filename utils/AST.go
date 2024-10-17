package utils

import (
	"fmt"

	"github.com/yash7xm/Rule_Engine_with_AST/parser"
)

// Helper function to print the entire AST recursively
func PrintAST(node *parser.Node, depth int) {
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
		PrintAST(node.Left, depth+1)
	}
	if node.Right != nil {
		fmt.Printf("%sRight:\n", indent)
		PrintAST(node.Right, depth+1)
	}
}
