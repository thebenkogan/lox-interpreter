package lexer

import (
	"bufio"
	"fmt"
)

type TokenType int

const (
	TokenTypeEOF TokenType = iota
	TokenTypeLeftParen
	TokenTypeRightParen
	TokenTypeLeftBrace
	TokenTypeRightBrace
	TokenTypeComma
	TokenTypeDot
	TokenTypeMinus
	TokenTypePlus
	TokenTypeSemicolon
	TokenTypeStar
	TokenTypeSlash
	TokenTypeEqual
	TokenTypeEqualEqual
	TokenTypeUnknown
)

func stringToType(s string, stream *bufio.Reader) (TokenType, string) {
	switch s {
	case "(":
		return TokenTypeLeftParen, s
	case ")":
		return TokenTypeRightParen, s
	case "{":
		return TokenTypeLeftBrace, s
	case "}":
		return TokenTypeRightBrace, s
	case ",":
		return TokenTypeComma, s
	case ".":
		return TokenTypeDot, s
	case "-":
		return TokenTypeMinus, s
	case "+":
		return TokenTypePlus, s
	case ";":
		return TokenTypeSemicolon, s
	case "*":
		return TokenTypeStar, s
	case "/":
		return TokenTypeSlash, s
	case "=":
		next, _ := stream.Peek(1)
		if len(next) > 0 && next[0] == '=' {
			_, _ = stream.ReadByte()
			return TokenTypeEqualEqual, "=="
		}
		return TokenTypeEqual, s
	default:
		return TokenTypeUnknown, ""
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
	case TokenTypeLeftBrace:
		return "LEFT_BRACE"
	case TokenTypeRightBrace:
		return "RIGHT_BRACE"
	case TokenTypeComma:
		return "COMMA"
	case TokenTypeDot:
		return "DOT"
	case TokenTypeMinus:
		return "MINUS"
	case TokenTypePlus:
		return "PLUS"
	case TokenTypeSemicolon:
		return "SEMICOLON"
	case TokenTypeStar:
		return "STAR"
	case TokenTypeSlash:
		return "SLASH"
	case TokenTypeEqual:
		return "EQUAL"
	case TokenTypeEqualEqual:
		return "EQUAL_EQUAL"
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

type TokenError struct {
	line  int
	token string
}

var errorRunes = []rune{'@', '#', '^', '$', '%'}

func (te *TokenError) String() string {
	return fmt.Sprintf("[line %d] Error: Unexpected character: %s", te.line, te.token)
}
