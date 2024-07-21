package lexer

import (
	"bufio"
	"io"
)

func Tokenize(file io.Reader) []Token {
	f := bufio.NewReader(file)
	tokens := make([]Token, 0)
	for {
		_, err := f.ReadString(' ')
		if err == io.EOF {
			tokens = append(tokens, Token{Type: TokenTypeEOF})
			break
		}
	}
	return tokens
}
