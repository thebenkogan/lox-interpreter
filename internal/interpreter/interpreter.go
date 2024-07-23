package interpreter

import (
	"io"

	"github.com/thebenkogan/lox-interpreter/internal/evaluator"
	"github.com/thebenkogan/lox-interpreter/internal/lexer"
	"github.com/thebenkogan/lox-interpreter/internal/parser"
)

type InterpreterError interface {
	Code() int
	Error() string
}

type Interpreter struct {
	env    *evaluator.Environment
	output io.Writer
}

func NewInterpreter(output io.Writer) *Interpreter {
	return &Interpreter{env: evaluator.NewEnvironment(), output: output}
}

func (i *Interpreter) Interpret(f io.Reader) InterpreterError {
	tokens, lexerErr := lexer.Tokenize(f)
	if lexerErr != nil {
		return lexerErr
	}

	statements, parserErr := parser.Parse(tokens)
	if parserErr != nil {
		return parserErr
	}

	for _, statement := range statements {
		err := statement.Execute(i.env, i.output)
		if err != nil {
			return err
		}
	}

	return nil
}
