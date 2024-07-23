package evaluator

import "fmt"

type Environment struct {
	mem map[string]any
}

func NewEnvironment() *Environment {
	return &Environment{mem: make(map[string]any)}
}

func (e *Environment) Get(name string) (any, *RuntimeError) {
	if val, ok := e.mem[name]; ok {
		return val, nil
	}
	return nil, NewRuntimeError(fmt.Sprintf("Undefined variable: %q", name))
}

func (e *Environment) Set(name string, val any) {
	e.mem[name] = val
}
