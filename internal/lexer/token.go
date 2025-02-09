package lexer

import (
	"bufio"
	"errors"
	"fmt"
	"strings"
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
	TokenTypeString
	TokenTypeNumber
	TokenTypeIdentifier
	TokenTypeAnd
	TokenTypeClass
	TokenTypeElse
	TokenTypeFalse
	TokenTypeFor
	TokenTypeFun
	TokenTypeIf
	TokenTypeNil
	TokenTypeOr
	TokenTypePrint
	TokenTypeReturn
	TokenTypeSuper
	TokenTypeThis
	TokenTypeTrue
	TokenTypeVar
	TokenTypeWhile
	TokenTypeUnknown
)

var reserved = map[string]TokenType{
	"and":    TokenTypeAnd,
	"class":  TokenTypeClass,
	"else":   TokenTypeElse,
	"false":  TokenTypeFalse,
	"for":    TokenTypeFor,
	"fun":    TokenTypeFun,
	"if":     TokenTypeIf,
	"nil":    TokenTypeNil,
	"or":     TokenTypeOr,
	"print":  TokenTypePrint,
	"return": TokenTypeReturn,
	"super":  TokenTypeSuper,
	"this":   TokenTypeThis,
	"true":   TokenTypeTrue,
	"var":    TokenTypeVar,
	"while":  TokenTypeWhile,
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
	case TokenTypeString:
		return "STRING"
	case TokenTypeNumber:
		return "NUMBER"
	case TokenTypeIdentifier:
		return "IDENTIFIER"
	case TokenTypeAnd:
		return "AND"
	case TokenTypeClass:
		return "CLASS"
	case TokenTypeElse:
		return "ELSE"
	case TokenTypeFalse:
		return "FALSE"
	case TokenTypeFor:
		return "FOR"
	case TokenTypeFun:
		return "FUN"
	case TokenTypeIf:
		return "IF"
	case TokenTypeNil:
		return "NIL"
	case TokenTypeOr:
		return "OR"
	case TokenTypePrint:
		return "PRINT"
	case TokenTypeReturn:
		return "RETURN"
	case TokenTypeSuper:
		return "SUPER"
	case TokenTypeThis:
		return "THIS"
	case TokenTypeTrue:
		return "TRUE"
	case TokenTypeVar:
		return "VAR"
	case TokenTypeWhile:
		return "WHILE"
	default:
		panic("Unknown token type")
	}
}

type Token struct {
	Type    TokenType
	Lexeme  string
	Literal string
}

func peekNext(stream *bufio.Reader) rune {
	next, _ := stream.Peek(1)
	if len(next) > 0 {
		return rune(next[0])
	}
	return rune(0)
}

func readToken(s rune, stream *bufio.Reader) (*Token, error) {
	switch s {
	case '@':
		fallthrough
	case '#':
		fallthrough
	case '^':
		fallthrough
	case '$':
		fallthrough
	case '&':
		fallthrough
	case '%':
		return nil, fmt.Errorf("Unexpected character: %s", string(s))
	case '(':
		return &Token{Type: TokenTypeLeftParen, Lexeme: string(s)}, nil
	case ')':
		return &Token{Type: TokenTypeRightParen, Lexeme: string(s)}, nil
	case '{':
		return &Token{Type: TokenTypeLeftBrace, Lexeme: string(s)}, nil
	case '}':
		return &Token{Type: TokenTypeRightBrace, Lexeme: string(s)}, nil
	case ',':
		return &Token{Type: TokenTypeComma, Lexeme: string(s)}, nil
	case '.':
		return &Token{Type: TokenTypeDot, Lexeme: string(s)}, nil
	case '-':
		return &Token{Type: TokenTypeMinus, Lexeme: string(s)}, nil
	case '+':
		return &Token{Type: TokenTypePlus, Lexeme: string(s)}, nil
	case ';':
		return &Token{Type: TokenTypeSemicolon, Lexeme: string(s)}, nil
	case '*':
		return &Token{Type: TokenTypeStar, Lexeme: string(s)}, nil
	case '/':
		return &Token{Type: TokenTypeSlash, Lexeme: string(s)}, nil
	case '=':
		if peekNext(stream) == '=' {
			_, _ = stream.Discard(1)
			return &Token{Type: TokenTypeEqualEqual, Lexeme: "=="}, nil
		}
		return &Token{Type: TokenTypeEqual, Lexeme: string(s)}, nil
	case '!':
		if peekNext(stream) == '=' {
			_, _ = stream.Discard(1)
			return &Token{Type: TokenTypeBangEqual, Lexeme: "!="}, nil
		}
		return &Token{Type: TokenTypeBang, Lexeme: string(s)}, nil
	case '<':
		if peekNext(stream) == '=' {
			_, _ = stream.Discard(1)
			return &Token{Type: TokenTypeLessEqual, Lexeme: "<="}, nil
		}
		return &Token{Type: TokenTypeLess, Lexeme: string(s)}, nil
	case '>':
		if peekNext(stream) == '=' {
			_, _ = stream.Discard(1)
			return &Token{Type: TokenTypeGreaterEqual, Lexeme: ">="}, nil
		}
		return &Token{Type: TokenTypeGreater, Lexeme: string(s)}, nil
	case '"':
		return readString(s, stream)
	default:
		if isDigit(s) {
			return readNumber(s, stream)
		}
		if isAlpha(s) {
			return readIdentifierOrReserved(s, stream)
		}
		return &Token{Type: TokenTypeUnknown, Lexeme: string(s)}, nil
	}
}

func readString(_ rune, stream *bufio.Reader) (*Token, error) {
	rest, err := stream.ReadString('"')
	if err != nil {
		return nil, errors.New("Unterminated string.")
	}
	literal := strings.TrimSuffix(rest, "\"")
	return &Token{Type: TokenTypeString, Lexeme: fmt.Sprintf("\"%s\"", literal), Literal: literal}, nil
}

func readNumber(s rune, stream *bufio.Reader) (*Token, error) {
	number := string(s)
	for isDigit(peekNext(stream)) {
		number += string(peekNext(stream))
		_, _ = stream.Discard(1)
	}
	if peekNext(stream) == '.' {
		nextTwo, _ := stream.Peek(2)
		if len(nextTwo) == 2 && isDigit(rune(nextTwo[1])) {
			_, _ = stream.Discard(1)
			number += "."
			for isDigit(peekNext(stream)) {
				number += string(peekNext(stream))
				_, _ = stream.Discard(1)
			}
		}
	}
	literal := number
	if !strings.Contains(number, ".") {
		literal += ".0"
	}
	for strings.HasSuffix(literal, "00") {
		literal = strings.TrimSuffix(literal, "0")
	}
	return &Token{Type: TokenTypeNumber, Lexeme: number, Literal: literal}, nil
}

func readIdentifierOrReserved(s rune, stream *bufio.Reader) (*Token, error) {
	ident := string(s)
	for isAlphaNumeric(peekNext(stream)) {
		ident += string(peekNext(stream))
		_, _ = stream.Discard(1)
	}
	tokenType := TokenTypeIdentifier
	if reservedType, ok := reserved[ident]; ok {
		tokenType = reservedType
	}
	return &Token{Type: tokenType, Lexeme: ident}, nil
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
	line int
	msg  string
}

func (te *TokenError) String() string {
	return fmt.Sprintf("[line %d] Error: %s", te.line, te.msg)
}

func isWhitespace(r rune) bool {
	return r == ' ' || r == '\t'
}

func isDigit(r rune) bool {
	return r >= '0' && r <= '9'
}

func isAlpha(r rune) bool {
	return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || r == '_'
}

func isAlphaNumeric(r rune) bool {
	return isAlpha(r) || isDigit(r)
}
