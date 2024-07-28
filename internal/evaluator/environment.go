package evaluator

import "fmt"

type Environment struct {
	mem    map[string]Value
	parent *Environment
}

func NewEnvironment() *Environment {
	return &Environment{mem: make(map[string]Value)}
}

func (e *Environment) Get(name string) (Value, *RuntimeError) {
	if val, ok := e.mem[name]; ok {
		return val, nil
	}
	if e.parent != nil {
		return e.parent.Get(name)
	}
	return nil, NewRuntimeError(fmt.Sprintf("Undefined variable: %q", name))
}

func (e *Environment) Declare(name string, val Value) {
	e.mem[name] = val
}

func (e *Environment) Set(name string, val Value) *RuntimeError {
	if _, ok := e.mem[name]; ok {
		e.mem[name] = val
		return nil
	}
	if e.parent != nil {
		return e.parent.Set(name, val)
	}
	return NewRuntimeError(fmt.Sprintf("Undefined variable: %q", name))
}

func (e *Environment) CreateScope() *Environment {
	return &Environment{mem: make(map[string]Value), parent: e}
}
