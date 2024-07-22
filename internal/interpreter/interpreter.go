package interpreter

import (
	"io"

	"github.com/thebenkogan/lox-interpreter/internal/lexer"
	"github.com/thebenkogan/lox-interpreter/internal/parser"
)

type InterpreterError interface {
	Code() int
	Error() string
}

func Interpret(f io.Reader) InterpreterError {
	tokens, lexerErr := lexer.Tokenize(f)
	if lexerErr != nil {
		return lexerErr
	}

	statements, parserErr := parser.Parse(tokens)
	if parserErr != nil {
		return parserErr
	}

	for _, statement := range statements {
		err := statement.Execute()
		if err != nil {
			return err
		}
	}

	return nil
}
