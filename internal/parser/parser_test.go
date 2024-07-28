package parser

import (
	"bytes"
	"testing"

	"github.com/thebenkogan/lox-interpreter/internal/evaluator"
	"github.com/thebenkogan/lox-interpreter/internal/lexer"
)

func TestParseExpression(t *testing.T) {
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
		{
			name:     "logic or",
			program:  "true or false",
			expected: "(or true false)",
		},
		{
			name:     "logic and",
			program:  "true and false",
			expected: "(and true false)",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			tokens, _ := lexer.Tokenize(bytes.NewBuffer([]byte(test.program + ";")))
			statements, err := Parse(tokens)
			if err != nil && !test.expectError {
				t.Errorf("Expected no error, got %v", err)
			}
			if err == nil && test.expectError {
				t.Errorf("Expected error, got nil")
			}
			if err != nil {
				return
			}
			expr := statements[0].(*evaluator.ExpressionStatement).Expression
			if test.expected != "" && expr.String() != test.expected {
				t.Errorf("Expected %s, got %s", test.expected, expr.String())
			}
		})
	}
}

func TestParseStatement(t *testing.T) {
	tests := []struct {
		name        string
		program     string
		expected    string
		expectError bool
	}{
		{
			name:     "expression statement",
			program:  "2 + 3;",
			expected: "(expr (+ 2.0 3.0))",
		},
		{
			name:        "expression statement no semicolon",
			program:     "2 + 3",
			expectError: true,
		},
		{
			name:     "print statement",
			program:  "print 2 + 3;",
			expected: "print (+ 2.0 3.0)",
		},
		{
			name:        "print statement no semicolon",
			program:     "print 2 + 3",
			expectError: true,
		},
		{
			name:        "variable statement no identifier",
			program:     "var 123;",
			expectError: true,
		},
		{
			name:     "variable statement no initializer",
			program:  "var a;",
			expected: "var a",
		},
		{
			name:     "variable statement",
			program:  "var b = 2 + 3;",
			expected: "var b = (+ 2.0 3.0)",
		},
		{
			name:        "variable statement no semicolon",
			program:     "var b = 2 + 3",
			expectError: true,
		},
		{
			name:     "assignment",
			program:  "a = 5;",
			expected: "(expr (= a 5.0))",
		},
		{
			name:        "assignment to non-variable",
			program:     "3 = 5;",
			expectError: true,
		},
		{
			name:     "block statement",
			program:  "{2 + 3; var b = 5;}",
			expected: "(block (expr (+ 2.0 3.0)); var b = 5.0;)",
		},
		{
			name:        "block statement no right brace",
			program:     "{2 + 3; var b = 5;",
			expectError: true,
		},
		{
			name:        "block statement no interior semicolon",
			program:     "{2 + 3; var b = 5}",
			expectError: true,
		},
		{
			name:     "if statement",
			program:  "if (true) {2 + 3;}",
			expected: "if (true) then (block (expr (+ 2.0 3.0));)",
		},
		{
			name:     "if statement with else",
			program:  "if (true) {2 + 3;} else {4 + 5;}",
			expected: "if (true) then (block (expr (+ 2.0 3.0));) else (block (expr (+ 4.0 5.0));)",
		},
		{
			name:        "if statement no block statement",
			program:     "if (true) 2 + 3;",
			expectError: true,
		},
		{
			name:        "if statement no parens",
			program:     "if true {2 + 3;}",
			expectError: true,
		},
		{
			name:        "if statement unclosed parens",
			program:     "if (true {2 + 3;}",
			expectError: true,
		},
		{
			name:        "if statement empty else",
			program:     "if (true) {2 + 3;} else",
			expectError: true,
		},
		{
			name:        "if statement no block after else",
			program:     "if (true) {2 + 3;} else 2 + 3;",
			expectError: true,
		},
		{
			name:     "while statement",
			program:  "while (true) {2 + 3;}",
			expected: "while (true) then (block (expr (+ 2.0 3.0));)",
		},
		{
			name:        "while statement no parens",
			program:     "while true {2 + 3;}",
			expectError: true,
		},
		{
			name:        "while statement unclosed parens",
			program:     "while (true {2 + 3;}",
			expectError: true,
		},
		{
			name:        "while statement no block",
			program:     "while (true) 2 + 3;",
			expectError: true,
		},
		{
			name:     "for statement",
			program:  "for (;;) {print 1;}",
			expected: "for (; ; ) (block print 1.0;)",
		},
		{
			name:     "for statement with initializer",
			program:  "for (var a = 1;;) {print a;}",
			expected: "for (var a = 1.0; ; ) (block print a;)",
		},
		{
			name:     "for statement with condition",
			program:  "for (; a < 3;) {print a;}",
			expected: "for (; (< a 3.0); ) (block print a;)",
		},
		{
			name:     "for statement with increment",
			program:  "for (var a = 1; a < 3; a = a + 1) {print a;}",
			expected: "for (var a = 1.0; (< a 3.0); (= a (+ a 1.0))) (block print a;)",
		},
		{
			name:        "for statement no parens",
			program:     "for var a = 1; a < 3; a = a + 1 {print a;}",
			expectError: true,
		},
		{
			name:        "for statement invalid initializer",
			program:     "for (if (true) {print 1;}; a < 3; a = a + 1) {print a;}",
			expectError: true,
		},
		{
			name:        "for statement no semicolon after condition",
			program:     "for (var a = 1; a < 3) {print a;}",
			expectError: true,
		},
		{
			name:        "for statement unclosed parens",
			program:     "for (var a = 1; a < 3; a = a + 1 {print a;}",
			expectError: true,
		},
		{
			name:        "for statement no block statement",
			program:     "for (var a = 1; a < 3; a = a + 1) print a;",
			expectError: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			tokens, _ := lexer.Tokenize(bytes.NewBuffer([]byte(test.program)))
			statements, err := Parse(tokens)
			if err != nil && !test.expectError {
				t.Errorf("Expected no error, got %v", err)
			}
			if err == nil && test.expectError {
				t.Errorf("Expected error, got nil")
			}
			if err != nil {
				return
			}
			stmt := statements[0]
			if test.expected != "" && stmt.String() != test.expected {
				t.Errorf("Expected %s, got %s", test.expected, stmt.String())
			}
		})
	}
}
