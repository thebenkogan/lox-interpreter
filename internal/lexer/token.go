package lexer

import "fmt"

type TokenType int

const (
	TokenTypeEOF TokenType = iota
	TokenTypeLeftParen
	TokenTypeRightParen
	TokenTypeUnknown
)

func typeFromString(s string) TokenType {
	switch s {
	case "(":
		return TokenTypeLeftParen
	case ")":
		return TokenTypeRightParen
	default:
		return TokenTypeUnknown
	}
}

func (t TokenType) String() string {
	switch t {
	case TokenTypeEOF:
		return "EOF"
	case TokenTypeLeftParen:
		return "LEFT_PAREN"
	case TokenTypeRightParen:
		return "RIGHT_PAREN"
	default:
		panic("Unknown token type")
	}
}

type Token struct {
	Type    TokenType
	Lexeme  string
	Literal string
}

func (t *Token) String() string {
	switch t.Type {
	case TokenTypeEOF:
		return "EOF  null"
	default:
		literal := t.Literal
		if literal == "" {
			literal = "null"
		}
		return fmt.Sprintf("%s %s %s", t.Type.String(), t.Lexeme, literal)
	}
}
