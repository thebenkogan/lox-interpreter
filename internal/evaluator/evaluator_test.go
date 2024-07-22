package evaluator_test

import (
	"bytes"
	"testing"

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
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			buf := bytes.NewBuffer([]byte(test.program))
			tokens, _ := lexer.Tokenize(buf)
			expr, _ := parser.Parse(tokens)
			result, err := expr.Evaluate()
			if test.expectError && err == nil {
				t.Errorf("Expected error, got nil")
				return
			}
			if !test.expectError && err != nil {
				t.Errorf("Expected no error, got %v", err)
				return
			}
			if result != test.expected {
				t.Errorf("Expected %v, got %v", test.expected, result)
			}
		})
	}
}
