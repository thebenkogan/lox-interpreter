package parser

import (
	"bytes"
	"testing"

	"github.com/codecrafters-io/interpreter-starter-go/internal/lexer"
)

func TestParser(t *testing.T) {
	tests := []struct {
		name     string
		program  string
		expected string
	}{
		{
			name:     "nil",
			program:  "nil",
			expected: "nil",
		},
		{
			name:     "false",
			program:  "false",
			expected: "false",
		},
		{
			name:     "true",
			program:  "true",
			expected: "true",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			tokens, _ := lexer.Tokenize(bytes.NewBuffer([]byte(test.program)))
			expr := Parse(tokens)
			if expr.String() != test.expected {
				t.Errorf("Expected %s, got %s", test.expected, expr.String())
			}
		})
	}
}
