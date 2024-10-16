package parser

import (
	"fmt"
)

// Node represents a node in the Abstract Syntax Tree.
type Node struct {
	Type  string
	Value string
	Left  *Node
	Right *Node
}

// Parser holds the tokenizer and the current token being processed.
type Parser struct {
	tokenizer *Tokenizer
	lookahead *Token
}

// NewParser creates a new Parser instance with the given tokenizer.
func NewParser(tokenizer *Tokenizer) *Parser {
	return &Parser{
		tokenizer: tokenizer,
	}
}

// ParseRule parses the entire rule and returns the AST.
func (p *Parser) ParseRule() *Node {
	p.lookahead = p.tokenizer.GetNextToken() // Initialize the first token.
	return p.Construct()
}

func (p *Parser) Construct() *Node {
	return p.LogicalOrExpression()
}

func (p *Parser) LogicalOrExpression() *Node {
	left := p.LogicalAndExpression()

	for p.lookahead.Type == "LOGICAL_OR" {
		operator, _ := p.eat("LOGICAL_OR")
		right := p.LogicalAndExpression()

		return &Node{
			Type:  "LogicalOrExpression",
			Value: operator.Value,
			Left:  left,
			Right: right,
		}
	}

	return left
}

func (p *Parser) LogicalAndExpression() *Node {
	left := p.EqualityExpression()

	for p.lookahead.Type == "LOGICAL_AND" {
		operator, _ := p.eat("LOGICAL_AND")
		right := p.EqualityExpression()

		return &Node{
			Type:  "LogicalAndExpression",
			Value: operator.Value,
			Left:  left,
			Right: right,
		}

	}

	return left
}

func (p *Parser) EqualityExpression() *Node {
	left := p.RelationalExpression()

	for p.lookahead.Type == "EQUALITY_OPERATOR" {
		operator, _ := p.eat("EQUALITY_OPERATOR")
		right := p.RelationalExpression()

		return &Node{
			Type:  "BinaryExpression",
			Value: operator.Value,
			Left:  left,
			Right: right,
		}

	}

	return left
}

func (p *Parser) RelationalExpression() *Node {
	left := p.UnaryExpression()

	for p.lookahead.Type == "RELATIONAL_OPERATOR" {
		operator, _ := p.eat("RELATIONAL_OPERATOR")
		right := p.UnaryExpression()

		return &Node{
			Type:  "BinaryExpression",
			Value: operator.Value,
			Left:  left,
			Right: right,
		}

	}

	return left
}

func (p *Parser) UnaryExpression() *Node {

	return &Node{
		Type:  "LogicalAndExpression",
		Value: "someValue",
	}
}

// eat consumes the current token if it matches the expected type and returns it.
func (p *Parser) eat(tokenType string) (*Token, error) {
	token := p.lookahead
	if token == nil {
		return nil, fmt.Errorf("unexpected end of input, expected: %s", tokenType)
	}

	if token.Type != tokenType {
		return nil, fmt.Errorf("unexpected token: %s, expected: %s", token.Value, tokenType)
	}

	p.lookahead = p.tokenizer.GetNextToken() // Move to the next token.
	return token, nil
}
