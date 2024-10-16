package main

import (
	"fmt"
	"github.com/yash7xm/Rule_Engine_with_AST/parser"
)

func main() {
	// rule := "((age > 30 AND department = 'Sales') OR (age < 25 AND department = 'Marketing')) AND (salary > 50000 OR experience > 5)"
	rule := "(age > 30 AND department = 'Sales') OR (department = 'Sales')"
	tokenizer := parser.NewTokenizer(rule)
	p := parser.NewParser(tokenizer)

	ast := p.ParseRule()
	fmt.Printf("AST: %+v\n", ast)
}
