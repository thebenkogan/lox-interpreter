package lexer

import (
	"bytes"
	"reflect"
	"slices"
	"testing"
)

func TestLexer(t *testing.T) {
	tests := []struct {
		name           string
		program        string
		expected       []TokenType
		expectedErrors []TokenError
	}{
		{
			name:    "single characters",
			program: "(){},.-+;*/=",
			expected: []TokenType{
				TokenTypeLeftParen,
				TokenTypeRightParen,
				TokenTypeLeftBrace,
				TokenTypeRightBrace,
				TokenTypeComma,
				TokenTypeDot,
				TokenTypeMinus,
				TokenTypePlus,
				TokenTypeSemicolon,
				TokenTypeStar,
				TokenTypeSlash,
				TokenTypeEqual,
				TokenTypeEOF,
			},
		},
		{
			name:     "single characters with errors",
			program:  "*%(\n^\n$)",
			expected: []TokenType{TokenTypeStar, TokenTypeLeftParen, TokenTypeRightParen, TokenTypeEOF},
			expectedErrors: []TokenError{
				{line: 1, token: "%"},
				{line: 2, token: "^"},
				{line: 3, token: "$"},
			},
		},
		{
			name:     "equal equal",
			program:  "===\n==\n=",
			expected: []TokenType{TokenTypeEqualEqual, TokenTypeEqual, TokenTypeEqualEqual, TokenTypeEqual, TokenTypeEOF},
		},
		{
			name:     "negation and inequality",
			program:  "!!===",
			expected: []TokenType{TokenTypeBang, TokenTypeBangEqual, TokenTypeEqualEqual, TokenTypeEOF},
		},
		{
			name:    "relational operators",
			program: "<<=>>=",
			expected: []TokenType{
				TokenTypeLess,
				TokenTypeLessEqual,
				TokenTypeGreater,
				TokenTypeGreaterEqual,
				TokenTypeEOF,
			},
		},
		{
			name:     "comments",
			program:  "// this is a comment\n==\n!%// this is another comment\n=//last one",
			expected: []TokenType{TokenTypeEqualEqual, TokenTypeBang, TokenTypeEqual, TokenTypeEOF},
			expectedErrors: []TokenError{
				{line: 3, token: "%"},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			stream := bytes.NewBuffer([]byte(test.program))
			tokens, errors := Tokenize(stream)

			tokenTypes := make([]TokenType, 0)
			for _, token := range tokens {
				tokenTypes = append(tokenTypes, token.Type)
			}
			if !slices.Equal(tokenTypes, test.expected) {
				t.Errorf("Expected token types %v, got %v", test.expected, tokenTypes)
			}

			if test.expectedErrors != nil && !reflect.DeepEqual(errors, test.expectedErrors) {
				t.Errorf("Expected errors %v, got %v", test.expectedErrors, errors)
			}
		})
	}
}
