package parser

import (
	"bytes"
	"testing"

	"github.com/codecrafters-io/interpreter-starter-go/internal/lexer"
)

func TestParser(t *testing.T) {
	tests := []struct {
		name          string
		program       string
		expected      string
		expectedError string
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
		{
			name:     "integer",
			program:  "123",
			expected: "123.0",
		},
		{
			name:     "float",
			program:  "123.23",
			expected: "123.23",
		},
		{
			name:     "string",
			program:  "\"hello\"",
			expected: "hello",
		},
		{
			name:     "parentheses",
			program:  "(123)",
			expected: "(group 123.0)",
		},
		{
			name:          "unmatched parentheses",
			program:       "(\"hello\"",
			expectedError: "Unmatched parentheses.",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			tokens, _ := lexer.Tokenize(bytes.NewBuffer([]byte(test.program)))
			expr, err := Parse(tokens)
			if test.expected != "" && expr.String() != test.expected {
				t.Errorf("Expected %s, got %s", test.expected, expr.String())
			}
			if test.expectedError != "" && test.expectedError != err.Error() {
				t.Errorf("Expected error %s, got %s", test.expectedError, err.Error())
			}
		})
	}
}
