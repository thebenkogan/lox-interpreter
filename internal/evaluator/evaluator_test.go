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
			if result != test.expected {
				t.Errorf("Expected %v, got %v", test.expected, result)
			}
		})
	}
}
