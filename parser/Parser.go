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

// ParseRule parses the entire rule and returns the AST or an error.
func (p *Parser) ParseRule() (*Node, error) {
	p.lookahead = p.tokenizer.GetNextToken()

	// Handle the case where the input is empty
	if p.lookahead == nil {
		return nil, fmt.Errorf("empty input, unable to parse rule")
	}

	ast, err := p.Construct()
	if err != nil {
		return nil, err
	}

	return ast, nil
}

func (p *Parser) Construct() (*Node, error) {
	return p.LogicalOrExpression()
}

func (p *Parser) LogicalOrExpression() (*Node, error) {
	left, err := p.LogicalAndExpression()
	if err != nil {
		return nil, err
	}

	for p.lookahead != nil && p.lookahead.Type == "LOGICAL_OR" {
		operator, err := p.eat("LOGICAL_OR")
		if err != nil {
			return nil, err
		}

		right, err := p.LogicalAndExpression()
		if err != nil {
			return nil, err
		}

		left = &Node{
			Type:  "LogicalOrExpression",
			Value: operator.Value,
			Left:  left,
			Right: right,
		}
	}

	return left, nil
}

func (p *Parser) LogicalAndExpression() (*Node, error) {
	left, err := p.EqualityExpression()
	if err != nil {
		return nil, err
	}

	for p.lookahead != nil && p.lookahead.Type == "LOGICAL_AND" {
		operator, err := p.eat("LOGICAL_AND")
		if err != nil {
			return nil, err
		}

		right, err := p.EqualityExpression()
		if err != nil {
			return nil, err
		}

		left = &Node{
			Type:  "LogicalAndExpression",
			Value: operator.Value,
			Left:  left,
			Right: right,
		}
	}

	return left, nil
}

func (p *Parser) EqualityExpression() (*Node, error) {
	left, err := p.RelationalExpression()
	if err != nil {
		return nil, err
	}

	for p.lookahead != nil && p.lookahead.Type == "EQUALITY_OPERATOR" {
		operator, err := p.eat("EQUALITY_OPERATOR")
		if err != nil {
			return nil, err
		}

		right, err := p.RelationalExpression()
		if err != nil {
			return nil, err
		}

		left = &Node{
			Type:  "BinaryExpression",
			Value: operator.Value,
			Left:  left,
			Right: right,
		}
	}

	return left, nil
}

func (p *Parser) RelationalExpression() (*Node, error) {
	left, err := p.PrimaryExpression()
	if err != nil {
		return nil, err
	}

	for p.lookahead != nil && p.lookahead.Type == "RELATIONAL_OPERATOR" {
		operator, err := p.eat("RELATIONAL_OPERATOR")
		if err != nil {
			return nil, err
		}

		right, err := p.PrimaryExpression()
		if err != nil {
			return nil, err
		}

		left = &Node{
			Type:  "BinaryExpression",
			Value: operator.Value,
			Left:  left,
			Right: right,
		}
	}

	return left, nil
}

func (p *Parser) PrimaryExpression() (*Node, error) {
	if p.lookahead == nil {
		return nil, fmt.Errorf("unexpected nil lookahead in PrimaryExpression")
	}

	if p.isLiteral(p.lookahead.Type) {
		return p.Literal()
	}

	switch p.lookahead.Type {
	case "IDENTIFIER":
		return p.Identifier()
	case "(":
		return p.ParenthesizedExpression()
	default:
		return nil, fmt.Errorf("unexpected token %s in PrimaryExpression", p.lookahead.Type)
	}
}

func (p *Parser) ParenthesizedExpression() (*Node, error) {
	_, err := p.eat("(")
	if err != nil {
		return nil, fmt.Errorf("error in ParenthesizedExpression: %w", err)
	}

	exp, err := p.LogicalOrExpression()
	if err != nil {
		return nil, err
	}

	_, err = p.eat(")")
	if err != nil {
		return nil, fmt.Errorf("error in ParenthesizedExpression: %w", err)
	}

	return exp, nil
}

func (p *Parser) Identifier() (*Node, error) {
	name, err := p.eat("IDENTIFIER")
	if err != nil {
		return nil, fmt.Errorf("error in Identifier: %w", err)
	}

	return &Node{
		Type:  "Identifier",
		Value: name.Value,
	}, nil
}

func (p *Parser) isLiteral(tokenType string) bool {
	return tokenType == "NUMBER" ||
		tokenType == "STRING" ||
		tokenType == "true" ||
		tokenType == "false" ||
		tokenType == "null"
}

func (p *Parser) Literal() (*Node, error) {
	if p.lookahead == nil {
		return nil, fmt.Errorf("nil lookahead in Literal")
	}

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
		return nil, fmt.Errorf("unexpected literal type %s", p.lookahead.Type)
	}
}

func (p *Parser) NumericLiteral() (*Node, error) {
	token, err := p.eat("NUMBER")
	if err != nil {
		return nil, fmt.Errorf("error in NumericLiteral: %w", err)
	}

	return &Node{
		Type:  "NumericLiteral",
		Value: token.Value,
	}, nil
}

func (p *Parser) StringLiteral() (*Node, error) {
	token, err := p.eat("STRING")
	if err != nil {
		return nil, fmt.Errorf("error in StringLiteral: %w", err)
	}

	return &Node{
		Type:  "StringLiteral",
		Value: token.Value,
	}, nil
}

func (p *Parser) BooleanLiteral(value string) (*Node, error) {
	token, err := p.eat(value)
	if err != nil {
		return nil, fmt.Errorf("error in BooleanLiteral: %w", err)
	}

	return &Node{
		Type:  "BooleanLiteral",
		Value: token.Value,
	}, nil
}

func (p *Parser) NullLiteral() (*Node, error) {
	_, err := p.eat("null")
	if err != nil {
		return nil, fmt.Errorf("error in NullLiteral: %w", err)
	}

	return &Node{
		Type:  "NullLiteral",
		Value: "null",
	}, nil
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
