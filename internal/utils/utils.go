package utils

import (
	"fmt"

	"github.com/yash7xm/Rule_Engine_with_AST/internal/parser"
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

func ConvertToASTNode(astJSON map[string]interface{}) (*parser.Node, error) {
	nodeType := astJSON["Type"].(string)

	switch nodeType {
	case "LogicalOrExpression":
		left, err := ConvertToASTNode(astJSON["Left"].(map[string]interface{}))
		if err != nil {
			return nil, err
		}
		right, err := ConvertToASTNode(astJSON["Right"].(map[string]interface{}))
		if err != nil {
			return nil, err
		}
		return &parser.Node{
			Type:  nodeType,
			Value: astJSON["Value"].(string),
			Left:  left,
			Right: right,
		}, nil

	case "LogicalAndExpression":
		left, err := ConvertToASTNode(astJSON["Left"].(map[string]interface{}))
		if err != nil {
			return nil, err
		}
		right, err := ConvertToASTNode(astJSON["Right"].(map[string]interface{}))
		if err != nil {
			return nil, err
		}
		return &parser.Node{
			Type:  nodeType,
			Value: astJSON["Value"].(string),
			Left:  left,
			Right: right,
		}, nil

	case "BinaryExpression":
		left, err := ConvertToASTNode(astJSON["Left"].(map[string]interface{}))
		if err != nil {
			return nil, err
		}
		right, err := ConvertToASTNode(astJSON["Right"].(map[string]interface{}))
		if err != nil {
			return nil, err
		}
		return &parser.Node{
			Type:  nodeType,
			Value: astJSON["Value"].(string),
			Left:  left,
			Right: right,
		}, nil

	case "Identifier":
		return &parser.Node{
			Type:  nodeType,
			Value: astJSON["Value"].(string),
		}, nil

	case "NumericLiteral":
		return &parser.Node{
			Type:  nodeType,
			Value: astJSON["Value"].(string),
		}, nil

	case "StringLiteral":
		return &parser.Node{
			Type:  nodeType,
			Value: astJSON["Value"].(string),
		}, nil

	default:
		return nil, fmt.Errorf("unknown node type: %s", nodeType)
	}
}
