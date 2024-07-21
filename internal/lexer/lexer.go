package lexer

import (
	"bufio"
	"io"
	"slices"
)

func Tokenize(file io.Reader) ([]Token, []TokenError) {
	f := bufio.NewReader(file)
	tokens := make([]Token, 0)
	errors := make([]TokenError, 0)
	line := 1
	var token string
	for {
		char, _, err := f.ReadRune()
		if err == io.EOF {
			tokens = append(tokens, Token{Type: TokenTypeEOF})
			break
		}
		if char == '\n' {
			line++
			continue
		}
		if slices.Contains(errorRunes, char) {
			errors = append(errors, TokenError{line: line, token: string(char)})
			continue
		}
		token += string(char)
		tokenType := stringToType(token, f)
		if tokenType == TokenTypeUnknown {
			continue
		}
		tokens = append(tokens, Token{Type: tokenType, Lexeme: token})
		token = ""
	}
	return tokens, errors
}
