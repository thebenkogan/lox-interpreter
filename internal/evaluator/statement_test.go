package evaluator_test

import (
	"bytes"
	"testing"

	"github.com/thebenkogan/lox-interpreter/internal/evaluator"
	"github.com/thebenkogan/lox-interpreter/internal/lexer"
	"github.com/thebenkogan/lox-interpreter/internal/parser"
)

func TestExecuteStatements(t *testing.T) {
	tests := []struct {
		name        string
		program     string
		expected    any
		expectError bool
	}{
		{
			name:     "expression statement does nothing",
			program:  "2 + 3;",
			expected: "",
		},
		{
			name:     "print statement",
			program:  "print 2 + 3;",
			expected: "5\n",
		},
		{
			name:     "uninitialized variable",
			program:  "var a; print a;",
			expected: "<nil>\n",
		},
		{
			name:     "initialized variable",
			program:  "var a = 123; print a;",
			expected: "123\n",
		},
		{
			name:     "initialized variable used in expression",
			program:  "var a = 123; print a + a;",
			expected: "246\n",
		},
		{
			name:        "accessing uninitialized variable fails",
			program:     "print a; var a = 123;",
			expectError: true,
		},
		{
			name:     "assignment",
			program:  "var a = 123; print a; a = 456; print a;",
			expected: "123\n456\n",
		},
		{
			name:        "assignment to uninitialized variable",
			program:     "a = 456; print a;",
			expectError: true,
		},
		{
			name:     "block statement",
			program:  "{var b = 5; print b + 3;}",
			expected: "8\n",
		},
		{
			name:     "block statements should create new scope",
			program:  "var a = 5; print a; {var a = 6; print a;} print a;",
			expected: "5\n6\n5\n",
		},
		{
			name:     "if statement",
			program:  "if (true) {print 1;}",
			expected: "1\n",
		},
		{
			name:     "if else statement",
			program:  "if (true) {print 1;} else {print 2;}",
			expected: "1\n",
		},
		{
			name:     "if else statement",
			program:  "if (false) {print 1;} else {print 2;}",
			expected: "2\n",
		},
		{
			name:     "if statement truthy condition",
			program:  "if (2) {print 1;} else {print 2;}",
			expected: "1\n",
		},
		{
			name:     "if statement false condition does nothing",
			program:  "if (false) {print 1;}",
			expected: "",
		},
		{
			name:     "while statement",
			program:  "var i = 0; while (i < 3) {print i; i = i + 1;}",
			expected: "0\n1\n2\n",
		},
		{
			name:     "for statement",
			program:  "for (var a = 1; a < 3; a = a + 1) {print a;}",
			expected: "1\n2\n",
		},
		{
			name:     "for statement assignment initializer",
			program:  "var a; for (a = 1; a < 3; a = a + 1) {print a;}",
			expected: "1\n2\n",
		},
		{
			name:     "for statement no increment",
			program:  "for (var a = 1; a < 3;) {print a; a = a + 1;}",
			expected: "1\n2\n",
		},
		{
			name:     "for statement no initializer",
			program:  "var a = 1; for (; a < 3; a = a + 1) {print a;}",
			expected: "1\n2\n",
		},
		{
			name:     "for statement should drop initializer out of scope",
			program:  "var a = 5; for (var a = 1; a < 3; a = a + 1) {print a;} print a;",
			expected: "1\n2\n5\n",
		},
		{
			name:     "fun statement",
			program:  "fun add(a, b) {print a + b;} print add;",
			expected: "<function>\n",
		},
		{
			name:     "fun call",
			program:  "fun add(a, b) {print a + b;} add(1, 2);",
			expected: "3\n",
		},
		{
			name:     "fun call references outer scope",
			program:  "var a = 1; fun add(b) {print a + b;} add(2);",
			expected: "3\n",
		},
		{
			name:     "fun call no args",
			program:  "var a = 1; var b = 2; fun add() {print a + b;} add();",
			expected: "3\n",
		},
		{
			name:     "fun call shadows outer scope",
			program:  "var a = 7; var b = 8; fun add(a, b) {print a + b;} add(1, 2); print a; print b;",
			expected: "3\n7\n8\n",
		},
		{
			name:        "fun call incorrect number of args",
			program:     "fun add(a, b) {print a + b;} add(1);",
			expectError: true,
		},
		{
			name:        "incorrect callee type",
			program:     "var add = 1; add();",
			expectError: true,
		},
		{
			name:     "fun call with no return is null",
			program:  "fun add(a, b) {a + b;} print add(1, 2);",
			expected: "<nil>\n",
		},
		{
			name:     "fun call with empty return is null",
			program:  "fun add(a, b) {a + b; return;} print add(1, 2);",
			expected: "<nil>\n",
		},
		{
			name:     "fun call with return is returned value",
			program:  "fun add(a, b) {return a + b;} print add(1, 2);",
			expected: "3\n",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			buf := bytes.NewBuffer([]byte(test.program))
			tokens, _ := lexer.Tokenize(buf)
			statements, _ := parser.Parse(tokens)
			env := evaluator.NewEnvironment()
			output := bytes.NewBuffer(nil)
			var err error
			for _, statement := range statements {
				execErr := statement.Execute(env, output)
				if execErr != nil {
					err = execErr
					break
				}
			}
			if test.expectError && err == nil {
				t.Errorf("Expected error, got nil")
			}
			if !test.expectError && err != nil {
				t.Errorf("Expected no error, got %v", err)
			}
			if err != nil {
				return
			}
			if output.String() != test.expected {
				t.Errorf("Expected %v, got %v", test.expected, output.String())
			}
		})
	}
}
