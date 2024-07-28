package evaluator_test

import (
	"bytes"
	"testing"

	"github.com/thebenkogan/lox-interpreter/internal/evaluator"
	"github.com/thebenkogan/lox-interpreter/internal/lexer"
	"github.com/thebenkogan/lox-interpreter/internal/parser"
)

func TestEvaluator(t *testing.T) {
	tests := []struct {
		name        string
		program     string
		expected    any
		expectError bool
	}{
		{
			name:     "literal",
			program:  "123",
			expected: float64(123),
		},
		{
			name:     "parentheses",
			program:  "(123)",
			expected: float64(123),
		},
		{
			name:     "unary not",
			program:  "!true",
			expected: false,
		},
		{
			name:     "unary not convert to bool",
			program:  "!nil",
			expected: true,
		},
		{
			name:        "unary minus wrong type",
			program:     "-\"hello\"",
			expectError: true,
		},
		{
			name:     "unary minus",
			program:  "-3",
			expected: float64(-3),
		},
		{
			name:     "multiply",
			program:  "2 * 3",
			expected: float64(6),
		},
		{
			name:        "multiple wrong types",
			program:     "2 * true",
			expectError: true,
		},
		{
			name:     "divide",
			program:  "2 / 3",
			expected: float64(2) / float64(3),
		},
		{
			name:        "divide wrong types",
			program:     "2 / true",
			expectError: true,
		},
		{
			name:        "divide by zero",
			program:     "2 / 0",
			expectError: true,
		},
		{
			name:     "add",
			program:  "3 + 5",
			expected: float64(8),
		},
		{
			name:     "add strings",
			program:  "\"hello\" + \" world\"",
			expected: "hello world",
		},
		{
			name:        "add wrong types",
			program:     "3 + true",
			expectError: true,
		},
		{
			name:     "subtract",
			program:  "3 - 5",
			expected: float64(-2),
		},
		{
			name:        "subtract wrong types",
			program:     "3 - true",
			expectError: true,
		},
		{
			name:     "greater",
			program:  "5 > 3",
			expected: true,
		},
		{
			name:        "greater wrong types",
			program:     "5 > true",
			expectError: true,
		},
		{
			name:     "greater equal",
			program:  "5 >= 5",
			expected: true,
		},
		{
			name:        "greater equal wrong types",
			program:     "5 >= true",
			expectError: true,
		},
		{
			name:     "less",
			program:  "5 < 3",
			expected: false,
		},
		{
			name:        "less wrong types",
			program:     "5 < true",
			expectError: true,
		},
		{
			name:     "less equal",
			program:  "5 <= 5",
			expected: true,
		},
		{
			name:        "less equal wrong types",
			program:     "5 <= true",
			expectError: true,
		},
		{
			name:     "equal",
			program:  "5 == 5",
			expected: true,
		},
		{
			name:     "not equal",
			program:  "5 != true",
			expected: true,
		},
		{
			name:        "assignment without declaration",
			program:     "a = 5",
			expectError: true,
		},
		{
			name:     "logic or",
			program:  "true or false",
			expected: true,
		},
		{
			name:     "logic or convert to truthy",
			program:  "false or 2",
			expected: float64(2),
		},
		{
			name:     "logic and",
			program:  "true and false",
			expected: false,
		},
		{
			name:     "logic and convert to falsy",
			program:  "\"hello\" and 2",
			expected: float64(2),
		},
	}

	for _, test := range tests {
		env := evaluator.NewEnvironment()
		t.Run(test.name, func(t *testing.T) {
			buf := bytes.NewBuffer([]byte(test.program + ";"))
			tokens, _ := lexer.Tokenize(buf)
			statements, _ := parser.Parse(tokens)
			expr := statements[0].(*evaluator.ExpressionStatement).Expression
			result, err := expr.Evaluate(env)
			if test.expectError && err == nil {
				t.Errorf("Expected error, got nil")
				return
			}
			if !test.expectError && err != nil {
				t.Errorf("Expected no error, got %v", err)
				return
			}
			if err != nil {
				return
			}
			val := result.(*evaluator.ValueLiteral)
			if val.Literal != test.expected {
				t.Errorf("Expected %v, got %v", test.expected, result)
			}
		})
	}
}
