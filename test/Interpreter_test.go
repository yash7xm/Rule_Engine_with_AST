package interpreter

import (
	"fmt"
	"testing"
	"time"

	"github.com/yash7xm/Rule_Engine_with_AST/internal/interpreter" // Importing interpreter
	"github.com/yash7xm/Rule_Engine_with_AST/internal/parser"
)

type Context = interpreter.Context // Alias for easier usage

// Helper function to test rule interpretation
func runTestRule(t *testing.T, rule string, ctx Context, expected bool) {
	tokenizer := parser.NewTokenizer(rule)
	p := parser.NewParser(tokenizer)

	ast, err := p.ParseRule()
	if err != nil {
		fmt.Println(err)
	}

	result := interpreter.Interpret(ast, ctx)
	if result != expected {
		t.Errorf("Rule: %s\nContext: %v\nExpected: %v, but got: %v\n", rule, ctx, expected, result)
	}
}

// Test cases
func TestBasicComparisons(t *testing.T) {
	tests := []struct {
		rule     string
		context  Context
		expected bool
	}{
		{"age > 30", Context{"age": 32}, true},
		{"age = 30", Context{"age": 30}, true},
		{"age < 25", Context{"age": 24}, true},
		{"salary > 50000", Context{"salary": 50001}, true},
		{"salary = 60000", Context{"salary": 60000}, true},
		{"salary < 50000", Context{"salary": 45000}, true},
		{"department = 'Sales'", Context{"department": "Sales"}, true},
		{"department = 'Sales'", Context{"department": "Marketing"}, false},
	}

	for _, test := range tests {
		runTestRule(t, test.rule, test.context, test.expected)
	}
}

func TestComplexRules(t *testing.T) {
	tests := []struct {
		rule     string
		context  Context
		expected bool
	}{
		{"age > 30 AND salary > 50000", Context{"age": 32, "salary": 60000}, true},
		{"age > 30 AND salary > 50000", Context{"age": 32, "salary": 40000}, false},
		{"age > 30 OR salary > 50000", Context{"age": 25, "salary": 60000}, true},
		{"age > 30 OR salary < 50000", Context{"age": 22, "salary": 30000}, true},
		{"(age > 30 AND department = 'Sales') OR (salary > 50000 AND department = 'Marketing')",
			Context{"age": 32, "department": "Sales", "salary": 60000}, true},
		{"(age > 30 AND department = 'Sales') OR (salary > 50000 AND department = 'Marketing')",
			Context{"age": 22, "department": "Sales", "salary": 60000}, false},
	}

	for _, test := range tests {
		runTestRule(t, test.rule, test.context, test.expected)
	}
}

// Edge case tests
func TestEdgeCases(t *testing.T) {
	tests := []struct {
		rule     string
		context  Context
		expected bool
	}{
		{"", Context{"age": 30}, false},
		{"age > 30", Context{}, false},
		{"salary > 50000", Context{"age": 40}, false},
		{"age > 30", Context{"age": "invalid_value"}, false},
		{"salary = 'abc'", Context{"salary": 123}, false},
	}

	for _, test := range tests {
		runTestRule(t, test.rule, test.context, test.expected)
	}
}

// Performance test
func TestPerformance(t *testing.T) {
	largeContext := Context{}
	for i := 1; i <= 1000; i++ {
		largeContext[fmt.Sprintf("key%d", i)] = i
	}

	rule := "(key500 > 250 AND key999 < 1000) OR key100 = 100"
	runTestRule(t, rule, largeContext, true)
}

func TestLargeRuleAndContextPerformance(t *testing.T) {
	largeContext := Context{}
	for i := 1; i <= 1000; i++ {
		largeContext[fmt.Sprintf("key%d", i)] = i
	}

	largeRule := "key1 > 0 AND key2 > 0 AND key3 > 0 AND key4 > 0 AND key5 > 0 " +
		"AND key6 > 0 AND key7 > 0 AND key8 > 0 AND key9 > 0 AND key10 > 0"

	start := time.Now()
	runTestRule(t, largeRule, largeContext, true)
	elapsed := time.Since(start)
	fmt.Printf("Performance test for large rule and context took: %s\n", elapsed)
}

// Test error handling and invalid inputs
func TestErrorHandling(t *testing.T) {
	tests := []struct {
		rule     string
		context  Context
		expected bool
	}{
		{"age >> 30", Context{"age": 32}, false}, // Invalid operator
		{"age >", Context{"age": 32}, false},
		{"age > 30", Context{"age": "thirty"}, false},
	}

	for _, test := range tests {
		runTestRule(t, test.rule, test.context, test.expected)
	}
}
