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
			name:     "assignment",
			program:  "a = 5",
			expected: float64(5),
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
			if result != test.expected {
				t.Errorf("Expected %v, got %v", test.expected, result)
			}
		})
	}
}

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
			name:     "block statement",
			program:  "{var b = 5; print b + 3;}",
			expected: "8\n",
		},
		{
			name:     "block statements should create new scope",
			program:  "var a = 5; print a; {var a = 6; print a;} print a;",
			expected: "5\n6\n5\n",
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
