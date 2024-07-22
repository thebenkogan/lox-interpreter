package evaluator

import (
	"errors"
	"fmt"
)

type RuntimeError struct {
	err error
}

func NewRuntimeError(msg string) *RuntimeError {
	return &RuntimeError{err: errors.New(msg)}
}

func (e *RuntimeError) Code() int {
	return 70
}

func (e *RuntimeError) Error() string {
	return fmt.Sprintf("Runtime Error: %s", e.err.Error())
}
