package evaluator

import "fmt"

type mem map[string]Value

type Environment struct {
	scopes []mem
}

func NewEnvironment() *Environment {
	initialScope := make(mem)
	return &Environment{scopes: []mem{initialScope}}
}

func (e *Environment) Get(name string) (Value, *RuntimeError) {
	for i := len(e.scopes) - 1; i >= 0; i-- {
		scope := e.scopes[i]
		if val, ok := scope[name]; ok {
			return val, nil
		}
	}
	return nil, NewRuntimeError(fmt.Sprintf("Undefined variable: %q", name))
}

func (e *Environment) Declare(name string, val Value) {
	scope := e.scopes[len(e.scopes)-1]
	scope[name] = val
}

func (e *Environment) Set(name string, val Value) *RuntimeError {
	for i := len(e.scopes) - 1; i >= 0; i-- {
		scope := e.scopes[i]
		if _, ok := scope[name]; ok {
			scope[name] = val
			return nil
		}
	}
	return NewRuntimeError(fmt.Sprintf("Undefined variable: %q", name))
}

func (e *Environment) CreateScope() {
	e.scopes = append(e.scopes, make(mem))
}

func (e *Environment) ExitScope() {
	if len(e.scopes) == 0 {
		panic("Cannot exit scope, no scopes to exit")
	}
	e.scopes = e.scopes[:len(e.scopes)-1]
}
