package lexer

type TokenType int

const (
	TokenTypeEOF TokenType = iota
)

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
		panic("Unknown token type")
	}
}
