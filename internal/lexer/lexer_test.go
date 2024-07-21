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
				{line: 1, msg: "Unexpected character: %"},
				{line: 2, msg: "Unexpected character: ^"},
				{line: 3, msg: "Unexpected character: $"},
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
				{line: 3, msg: "Unexpected character: %"},
			},
		},
		{
			name:    "whitespace",
			program: "( \t    )  \t  !\n   \t !=",
			expected: []TokenType{
				TokenTypeLeftParen,
				TokenTypeRightParen,
				TokenTypeBang,
				TokenTypeBangEqual,
				TokenTypeEOF,
			},
		},
		{
			name:    "strings",
			program: "\"hello\" \"world\" \"!\"\n\"unterminated",
			expected: []TokenType{
				TokenTypeString,
				TokenTypeString,
				TokenTypeString,
				TokenTypeEOF,
			},
			expectedErrors: []TokenError{
				{line: 2, msg: "Unterminated string."},
			},
		},
		{
			name:    "numbers",
			program: "1234 1234.5 . 1.234 12.34\n.1234\n1234.\n1234.1234.1234.",
			expected: []TokenType{
				TokenTypeNumber,
				TokenTypeNumber,
				TokenTypeDot,
				TokenTypeNumber,
				TokenTypeNumber,

				TokenTypeDot,
				TokenTypeNumber,

				TokenTypeNumber,
				TokenTypeDot,

				TokenTypeNumber,
				TokenTypeDot,
				TokenTypeNumber,
				TokenTypeDot,
				TokenTypeEOF,
			},
		},
		{
			name:    "identifiers",
			program: "foo bar baz",
			expected: []TokenType{
				TokenTypeIdentifier,
				TokenTypeIdentifier,
				TokenTypeIdentifier,
				TokenTypeEOF,
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
