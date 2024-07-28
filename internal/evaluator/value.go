package evaluator

import (
	"fmt"
)

type Value interface {
	String() string
	Bool() bool
}

type ValueLiteral struct {
	Literal any // number, string, bool, nil
}

func (v *ValueLiteral) String() string {
	return fmt.Sprintf("%v", v.Literal)
}

func (v *ValueLiteral) Bool() bool {
	return v.Literal != nil && v.Literal != false
}

type ValueClosure struct {
	Env    Environment
	Body   *BlockStatement
	Params []string
}

func (v *ValueClosure) String() string {
	return "<closure>"
}

func (v *ValueClosure) Bool() bool {
	return true
}
