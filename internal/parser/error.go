package parser

import (
	"errors"
	"fmt"
)

type ParserError struct {
	err error
}

func NewParserError(msg string) *ParserError {
	return &ParserError{err: errors.New(msg)}
}

func (e *ParserError) Code() int {
	return 65
}

func (e *ParserError) Error() string {
	return fmt.Sprintf("Parser Error: %s", e.err.Error())
}
