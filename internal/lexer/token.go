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
	TokenTypeBang
	TokenTypeBangEqual
	TokenTypeLess
	TokenTypeLessEqual
	TokenTypeGreater
	TokenTypeGreaterEqual
	TokenTypeUnknown
)

func peekNext(stream *bufio.Reader) rune {
	next, _ := stream.Peek(1)
	if len(next) > 0 {
		return rune(next[0])
	}
	return rune(0)
}

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
		if peekNext(stream) == '=' {
			_, _ = stream.ReadByte()
			return TokenTypeEqualEqual, "=="
		}
		return TokenTypeEqual, s
	case "!":
		if peekNext(stream) == '=' {
			_, _ = stream.ReadByte()
			return TokenTypeBangEqual, "!="
		}
		return TokenTypeBang, s
	case "<":
		if peekNext(stream) == '=' {
			_, _ = stream.ReadByte()
			return TokenTypeLessEqual, "<="
		}
		return TokenTypeLess, s
	case ">":
		if peekNext(stream) == '=' {
			_, _ = stream.ReadByte()
			return TokenTypeGreaterEqual, ">="
		}
		return TokenTypeGreater, s
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
	case TokenTypeBang:
		return "BANG"
	case TokenTypeBangEqual:
		return "BANG_EQUAL"
	case TokenTypeLess:
		return "LESS"
	case TokenTypeLessEqual:
		return "LESS_EQUAL"
	case TokenTypeGreater:
		return "GREATER"
	case TokenTypeGreaterEqual:
		return "GREATER_EQUAL"
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

func isWhitespace(r rune) bool {
	return r == ' ' || r == '\t'
}
