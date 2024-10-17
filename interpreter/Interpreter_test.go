package interpreter

import (
	"fmt"
	"testing"
	"time"

	"github.com/yash7xm/Rule_Engine_with_AST/parser"
)

// Helper function to test rule interpretation
func runTestRule(t *testing.T, rule string, context Context, expected bool) {
	tokenizer := parser.NewTokenizer(rule)
	p := parser.NewParser(tokenizer)

	ast := p.ParseRule()

	result := Interpret(ast, context)

	if result != expected {
		t.Errorf("Rule: %s\nContext: %v\nExpected: %v, but got: %v\n", rule, context, expected, result)
	}
}

// Basic tests for simple comparisons
func TestBasicComparisons(t *testing.T) {
	// Test with a simple rule
	runTestRule(t, "age > 30", Context{"age": 32}, true)
	runTestRule(t, "age = 30", Context{"age": 30}, true)
	runTestRule(t, "age < 25", Context{"age": 24}, true)
	runTestRule(t, "salary > 50000", Context{"salary": 50001}, true)
	runTestRule(t, "salary = 60000", Context{"salary": 60000}, true)
	runTestRule(t, "salary < 50000", Context{"salary": 45000}, true)
	runTestRule(t, "department = 'Sales'", Context{"department": "Sales"}, true)
	runTestRule(t, "department = 'Sales'", Context{"department": "Marketing"}, false)
}

// Test multiple logical expressions using AND/OR
func TestComplexRules(t *testing.T) {
	// Test with AND
	runTestRule(t, "age > 30 AND salary > 50000", Context{"age": 32, "salary": 60000}, true)
	runTestRule(t, "age > 30 AND salary > 50000", Context{"age": 32, "salary": 40000}, false)

	// Test with OR
	runTestRule(t, "age > 30 OR salary > 50000", Context{"age": 25, "salary": 60000}, true)
	runTestRule(t, "age > 30 OR salary < 50000", Context{"age": 22, "salary": 30000}, true)

	// Combination of AND/OR
	runTestRule(t, "(age > 30 AND department = 'Sales') OR (salary > 50000 AND department = 'Marketing')",
		Context{"age": 32, "department": "Sales", "salary": 60000}, true)
	runTestRule(t, "(age > 30 AND department = 'Sales') OR (salary > 50000 AND department = 'Marketing')",
		Context{"age": 22, "department": "Sales", "salary": 60000}, false)
}

// Edge case tests for invalid or empty rules
func TestEdgeCases(t *testing.T) {
	// Test with empty rule
	runTestRule(t, "", Context{"age": 30}, false)

	// Test with missing context values
	runTestRule(t, "age > 30", Context{}, false)
	runTestRule(t, "salary > 50000", Context{"age": 40}, false)

	// Test with invalid field types in context
	runTestRule(t, "age > 30", Context{"age": "invalid_value"}, false)
	runTestRule(t, "salary = 'abc'", Context{"salary": 123}, false)
}

// Boundary tests for numeric comparisons
func TestBoundaryConditions(t *testing.T) {
	// Boundary test for equality
	runTestRule(t, "age = 30", Context{"age": 30}, true)
	runTestRule(t, "age = 30", Context{"age": 31}, false)

	// Greater than and less than on boundary
	runTestRule(t, "age > 30", Context{"age": 30}, false)
	runTestRule(t, "age < 30", Context{"age": 30}, false)
	runTestRule(t, "age > 30", Context{"age": 31}, true)
	runTestRule(t, "age < 30", Context{"age": 29}, true)
}

// Performance tests for large context data
func TestPerformance(t *testing.T) {
	largeContext := Context{}
	// Create large context with 1000 entries
	for i := 1; i <= 1000; i++ {
		largeContext[fmt.Sprintf("key%d", i)] = i
	}

	rule := "(key500 > 250 AND key999 < 1000) OR key100 = 100"
	runTestRule(t, rule, largeContext, true)
}

// Performance test for large rule and context data
func TestLargeRuleAndContextPerformance(t *testing.T) {
	largeContext := Context{}
	// Create large context with 1000 entries
	for i := 1; i <= 1000; i++ {
		largeContext[fmt.Sprintf("key%d", i)] = i
	}

	// Construct a large rule
	// This rule checks multiple keys in the context
	// For simplicity, we'll create a rule that checks if key1 to key10 are greater than 0
	largeRule := "key1 > 0 AND key2 > 0 AND key3 > 0 AND key4 > 0 AND key5 > 0 " +
		"AND key6 > 0 AND key7 > 0 AND key8 > 0 AND key9 > 0 AND key10 > 0"

	// Start measuring time
	start := time.Now()

	// Run the test with the large rule and context
	runTestRule(t, largeRule, largeContext, true)

	// Log the time taken for the performance test
	elapsed := time.Since(start)
	fmt.Printf("Performance test for large rule and context took: %s\n", elapsed)
}

// Test error handling and invalid inputs
func TestErrorHandling(t *testing.T) {
	// Test for invalid binary expressions
	runTestRule(t, "age >> 30", Context{"age": 32}, false) // Invalid operator

	// Test for missing operands
	runTestRule(t, "age >", Context{"age": 32}, false)

	// Test for mismatched data types in context
	runTestRule(t, "age > 30", Context{"age": "thirty"}, false) // age should be a number
}
