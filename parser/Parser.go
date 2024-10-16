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

	for p.lookahead != nil && p.lookahead.Type == "LOGICAL_OR" {
		operator, _ := p.eat("LOGICAL_OR")
		right := p.LogicalAndExpression()

		left = &Node{
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

	for p.lookahead != nil && p.lookahead.Type == "LOGICAL_AND" {
		operator, _ := p.eat("LOGICAL_AND")
		right := p.EqualityExpression()

		left = &Node{
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

	for p.lookahead != nil && p.lookahead.Type == "EQUALITY_OPERATOR" {
		operator, _ := p.eat("EQUALITY_OPERATOR")
		right := p.RelationalExpression()

		left = &Node{
			Type:  "BinaryExpression",
			Value: operator.Value,
			Left:  left,
			Right: right,
		}
	}

	return left
}

func (p *Parser) RelationalExpression() *Node {
	left := p.PrimaryExpression()

	for p.lookahead != nil && p.lookahead.Type == "RELATIONAL_OPERATOR" {
		operator, _ := p.eat("RELATIONAL_OPERATOR")
		right := p.PrimaryExpression()

		left = &Node{
			Type:  "BinaryExpression",
			Value: operator.Value,
			Left:  left,
			Right: right,
		}
	}

	return left
}

func (p *Parser) PrimaryExpression() *Node {
	if p.isLiteral(p.lookahead.Type) {
		return p.Literal()
	}

	switch p.lookahead.Type {
	case "IDENTIFIER":
		return p.Identifier()
	case "(":
		return p.ParenthesizedExpression()
	default:
		fmt.Printf("PrimaryExpression: unexpected token %s\n", p.lookahead.Type)
		return nil
	}
}

func (p *Parser) ParenthesizedExpression() *Node {
	p.eat("(") // Consume '('
	exp := p.LogicalOrExpression()
	p.eat(")") // Consume ')'
	return exp
}

func (p *Parser) Identifier() *Node {
	name, _ := p.eat("IDENTIFIER")
	return &Node{
		Type:  "Identifier",
		Value: name.Value,
	}
}

func (p *Parser) isLiteral(tokenType string) bool {
	return tokenType == "NUMBER" ||
		tokenType == "STRING" ||
		tokenType == "true" ||
		tokenType == "false" ||
		tokenType == "null"
}

func (p *Parser) Literal() *Node {
	switch p.lookahead.Type {
	case "NUMBER":
		return p.NumericLiteral()
	case "STRING":
		return p.StringLiteral()
	case "true":
		return p.BooleanLiteral("true")
	case "false":
		return p.BooleanLiteral("false")
	case "null":
		return p.NullLiteral()
	default:
		fmt.Println("Literal: unexpected literal production.")
		return nil
	}
}

func (p *Parser) NumericLiteral() *Node {
	token, _ := p.eat("NUMBER")
	return &Node{
		Type:  "NumericLiteral",
		Value: token.Value,
	}
}

func (p *Parser) StringLiteral() *Node {
	token, _ := p.eat("STRING")
	return &Node{
		Type:  "StringLiteral",
		Value: token.Value,
	}
}

func (p *Parser) BooleanLiteral(value string) *Node {
	token, _ := p.eat(value)
	return &Node{
		Type:  "BooleanLiteral",
		Value: token.Value,
	}
}

func (p *Parser) NullLiteral() *Node {
	p.eat("null")
	return &Node{
		Type:  "NullLiteral",
		Value: "null",
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

	// Move to the next token.
	p.lookahead = p.tokenizer.GetNextToken()
	return token, nil
}
