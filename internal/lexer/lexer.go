package lexer

import (
	"bufio"
	"io"
)

func Tokenize(file io.Reader) []Token {
	f := bufio.NewReader(file)
	tokens := make([]Token, 0)
	var token string
	for {
		char, _, err := f.ReadRune()
		if err == io.EOF {
			tokens = append(tokens, Token{Type: TokenTypeEOF})
			break
		}
		token += string(char)
		tokenType := typeFromString(token)
		if tokenType == TokenTypeUnknown {
			continue
		}
		tokens = append(tokens, Token{Type: tokenType, Lexeme: token})
		token = ""
	}
	return tokens
}
