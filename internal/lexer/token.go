package lexer

import "fmt"

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
)

var stringToType = map[string]TokenType{
	"(": TokenTypeLeftParen,
	")": TokenTypeRightParen,
	"{": TokenTypeLeftBrace,
	"}": TokenTypeRightBrace,
	",": TokenTypeComma,
	".": TokenTypeDot,
	"-": TokenTypeMinus,
	"+": TokenTypePlus,
	";": TokenTypeSemicolon,
	"*": TokenTypeStar,
	"/": TokenTypeSlash,
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
