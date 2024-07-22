package lexer

import "strings"

type LexerError struct {
	Errors []TokenError
}

func (e *LexerError) Code() int {
	return 65
}

func (e *LexerError) Error() string {
	log := strings.Builder{}
	for _, err := range e.Errors {
		log.WriteString(err.String())
		log.WriteString("\n")
	}
	return log.String()
}
