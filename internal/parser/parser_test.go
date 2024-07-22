package parser

import (
	"bytes"
	"testing"

	"github.com/thebenkogan/lox-interpreter/internal/evaluator"
	"github.com/thebenkogan/lox-interpreter/internal/lexer"
)

func TestParser(t *testing.T) {
	tests := []struct {
		name        string
		program     string
		expected    string
		expectError bool
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
			name:        "unmatched parentheses",
			program:     "(\"hello\"",
			expectError: true,
		},
		{
			name:        "empty parentheses",
			program:     "()",
			expectError: true,
		},
		{
			name:     "unary bang",
			program:  "!true",
			expected: "(! true)",
		},
		{
			name:     "unary minus",
			program:  "-3",
			expected: "(- 3.0)",
		},
		{
			name:     "multiply",
			program:  "2 * 3",
			expected: "(* 2.0 3.0)",
		},
		{
			name:     "divide",
			program:  "2 / 3",
			expected: "(/ 2.0 3.0)",
		},
		{
			name:     "multiply and divide",
			program:  "16 * 38 / 58",
			expected: "(/ (* 16.0 38.0) 58.0)",
		},
		{
			name:     "add",
			program:  "3 + 5",
			expected: "(+ 3.0 5.0)",
		},
		{
			name:     "subtract",
			program:  "3 - 5",
			expected: "(- 3.0 5.0)",
		},
		{
			name:     "add and subtract",
			program:  "3 + 5 - 9",
			expected: "(- (+ 3.0 5.0) 9.0)",
		},
		{
			name:     "greater",
			program:  "3 > 5",
			expected: "(> 3.0 5.0)",
		},
		{
			name:     "greater equal",
			program:  "3 >= 5",
			expected: "(>= 3.0 5.0)",
		},
		{
			name:     "less",
			program:  "3 < 5",
			expected: "(< 3.0 5.0)",
		},
		{
			name:     "less equal",
			program:  "3 <= 5",
			expected: "(<= 3.0 5.0)",
		},
		{
			name:     "equal",
			program:  "3 == 5",
			expected: "(== 3.0 5.0)",
		},
		{
			name:     "not equal",
			program:  "3 != 5",
			expected: "(!= 3.0 5.0)",
		},
		{
			name:        "error add",
			program:     "(72 +)",
			expectError: true,
		},
		{
			name:        "error subtract",
			program:     "(72 -)",
			expectError: true,
		},
		{
			name:        "error multiply",
			program:     "(72 *)",
			expectError: true,
		},
		{
			name:        "error divide",
			program:     "(72 /)",
			expectError: true,
		},
		{
			name:        "error greater",
			program:     "(72 >)",
			expectError: true,
		},
		{
			name:        "error greater equal",
			program:     "(72 >=)",
			expectError: true,
		},
		{
			name:        "error less",
			program:     "(72 <)",
			expectError: true,
		},
		{
			name:        "error less equal",
			program:     "(72 <=)",
			expectError: true,
		},
		{
			name:        "error equal",
			program:     "(72 ==)",
			expectError: true,
		},
		{
			name:        "error not equal",
			program:     "(72 !=)",
			expectError: true,
		},
		{
			name:        "error bang",
			program:     "!",
			expectError: true,
		},
		{
			name:        "error minus",
			program:     "-",
			expectError: true,
		},
		{
			name:     "all types",
			program:  "((3.0 + 5.1) * 7) >= (-3.1 / -5.5) < (true != false) != \"hello world\" == nil",
			expected: "(== (!= (< (>= (group (* (group (+ 3.0 5.1)) 7.0)) (group (/ (- 3.1) (- 5.5)))) (group (!= true false))) hello world) nil)",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			tokens, _ := lexer.Tokenize(bytes.NewBuffer([]byte(test.program + ";")))
			statements, err := Parse(tokens)
			if test.expectError {
				if err == nil {
					t.Errorf("Expected error, got nil")
				}
				return
			}
			expr := statements[0].(*evaluator.ExpressionStatement).Expression
			if test.expected != "" && expr.String() != test.expected {
				t.Errorf("Expected %s, got %s", test.expected, expr.String())
			}
		})
	}
}
