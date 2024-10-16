package parser

import (
	"regexp"
)

// Token represents a token with a type and value.
type Token struct {
	Type  string
	Value string
}

// Spec defines the regular expressions for tokens and corresponding types.
var Spec = []struct {
	Pattern   *regexp.Regexp
	TokenType string
}{
	// Whitespace
	{regexp.MustCompile(`^\s+`), ""},

	// Comments
	{regexp.MustCompile(`^\/\/.*`), ""},
	{regexp.MustCompile(`^\/\*[\s\S]*?\*\/`), ""},

	// Symbols and delimiters
	{regexp.MustCompile(`^;`), ";"},
	{regexp.MustCompile(`^{`), "{"},
	{regexp.MustCompile(`^}`), "}"},
	{regexp.MustCompile(`^\(`), "("},
	{regexp.MustCompile(`^\)`), ")"},
	{regexp.MustCompile(`^\[`), "["},
	{regexp.MustCompile(`^\]`), "]"},
	{regexp.MustCompile(`^,`), ","},
	{regexp.MustCompile(`^\.`), "."},
	{regexp.MustCompile(`^\?`), "?"},
	{regexp.MustCompile(`^:`), ":"},

	// Relational operators
	{regexp.MustCompile(`^[<>]=?`), "RELATIONAL_OPERATOR"},
	{regexp.MustCompile(`^[=!]=`), "EQUALITY_OPERATOR"},

	// Logical operators
	{regexp.MustCompile(`^&&`), "LOGICAL_AND"},
	{regexp.MustCompile(`^\|\|`), "LOGICAL_OR"},
	{regexp.MustCompile(`^!`), "LOGICAL_NOT"},

	// Keywords
	{regexp.MustCompile(`^\blet\b`), "let"},
	{regexp.MustCompile(`^\bif\b`), "if"},
	{regexp.MustCompile(`^\belif\b`), "elif"},
	{regexp.MustCompile(`^\belse\b`), "else"},
	{regexp.MustCompile(`^\btrue\b`), "true"},
	{regexp.MustCompile(`^\bfalse\b`), "false"},
	{regexp.MustCompile(`^\bnull\b`), "null"},

	// OOP keywords
	{regexp.MustCompile(`^\bclass\b`), "class"},
	{regexp.MustCompile(`^\bthis\b`), "this"},
	{regexp.MustCompile(`^\bextends\b`), "extends"},
	{regexp.MustCompile(`^\bsuper\b`), "super"},
	{regexp.MustCompile(`^\bnew\b`), "new"},
	{regexp.MustCompile(`^\bmodule\b`), "module"},
	{regexp.MustCompile(`^\bimport\b`), "import"},

	// Iterators
	{regexp.MustCompile(`^\bwhile\b`), "while"},
	{regexp.MustCompile(`^\bdo\b`), "do"},
	{regexp.MustCompile(`^\bfor\b`), "for"},

	{regexp.MustCompile(`^\bdef\b`), "def"},
	{regexp.MustCompile(`^\breturn\b`), "return"},

	// Assignment operators
	{regexp.MustCompile(`^=`), "SIMPLE_ASSIGN"},
	{regexp.MustCompile(`^[\*\\/\+\-]=`), "COMPLEX_ASSIGN"},

	// Math operators
	{regexp.MustCompile(`^[+\-]`), "ADDITIVE_OPERATOR"},
	{regexp.MustCompile(`^[*\/]`), "MULTIPLICATIVE_OPERATOR"},
	{regexp.MustCompile(`^%`), "MODULO_OPERATOR"},

	// Numbers
	{regexp.MustCompile(`^\d+`), "NUMBER"},

	// Double quoted String
	{regexp.MustCompile(`^"[^"]*"`), "STRING"},

	// Single quoted String
	{regexp.MustCompile(`^'[^']*'`), "STRING"},

	// Identifier
	{regexp.MustCompile(`^\w+`), "IDENTIFIER"},
}

// Tokenizer holds the state for the string being tokenized.
type Tokenizer struct {
	input  string
	cursor int
}

// NewTokenizer creates a new Tokenizer instance.
func NewTokenizer(input string) *Tokenizer {
	return &Tokenizer{
		input:  input,
		cursor: 0,
	}
}

// hasMoreTokens checks if there are more tokens left to process.
func (t *Tokenizer) hasMoreTokens() bool {
	return t.cursor < len(t.input)
}

// GetNextToken extracts the next token from the input string.
func (t *Tokenizer) GetNextToken() *Token {
	if !t.hasMoreTokens() {
		return nil
	}

	input := t.input[t.cursor:]

	for _, spec := range Spec {
		if matched := t.match(spec.Pattern, input); matched != "" {
			t.cursor += len(matched)

			if spec.TokenType == "" {
				return t.GetNextToken() // skip whitespace/comments
			}

			return &Token{
				Type:  spec.TokenType,
				Value: matched,
			}
		}
	}

	return nil // No matching token found
}

// match tries to match the regular expression pattern at the start of the string.
func (t *Tokenizer) match(pattern *regexp.Regexp, input string) string {
	if matched := pattern.FindString(input); matched != "" {
		return matched
	}
	return ""
}
