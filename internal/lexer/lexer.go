package lexer

import (
	"bufio"
	"io"
)

func Tokenize(file io.Reader) ([]Token, *LexerError) {
	f := bufio.NewReader(file)
	tokens := make([]Token, 0)
	errors := make([]TokenError, 0)
	line := 1
	for {
		char, _, err := f.ReadRune()
		if err == io.EOF {
			tokens = append(tokens, Token{Type: TokenTypeEOF})
			break
		}
		if isWhitespace(char) {
			continue
		}
		if char == '\n' {
			line++
			continue
		}
		if char == '/' && peekNext(f) == '/' {
			_, _ = f.ReadString('\n')
			line++
			continue
		}
		parsed, err := readToken(char, f)
		if err != nil {
			errors = append(errors, TokenError{line: line, msg: err.Error()})
			continue
		}
		if parsed.Type == TokenTypeUnknown {
			continue
		}
		tokens = append(tokens, *parsed)
	}
	var err *LexerError
	if len(errors) > 0 {
		err = &LexerError{Errors: errors}
	}
	return tokens, err
}
