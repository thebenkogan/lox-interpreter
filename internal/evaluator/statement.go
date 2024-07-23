package evaluator

import (
	"fmt"
	"io"
)

type Statement interface {
	String() string
	Execute(env *Environment, output io.Writer) *RuntimeError
}

type ExpressionStatement struct {
	Expression Expression
}

func (e *ExpressionStatement) String() string {
	return fmt.Sprintf("(expr %s)", e.Expression.String())
}

func (e *ExpressionStatement) Execute(env *Environment, _ io.Writer) *RuntimeError {
	_, err := e.Expression.Evaluate(env)
	return err
}

type PrintStatement struct {
	Expression Expression
}

func (e *PrintStatement) String() string {
	return fmt.Sprintf("print %s", e.Expression.String())
}

func (e *PrintStatement) Execute(env *Environment, output io.Writer) *RuntimeError {
	result, err := e.Expression.Evaluate(env)
	if err != nil {
		return err
	}
	fmt.Fprintln(output, result)
	return nil
}

type VarStatement struct {
	Name string
	Expr Expression
}

func (e *VarStatement) String() string {
	if e.Expr == nil {
		return fmt.Sprintf("var %s", e.Name)
	}
	return fmt.Sprintf("var %s = %s", e.Name, e.Expr.String())
}

func (e *VarStatement) Execute(env *Environment, _ io.Writer) *RuntimeError {
	var value any = nil
	if e.Expr != nil {
		result, err := e.Expr.Evaluate(env)
		if err != nil {
			return err
		}
		value = result
	}
	env.Set(e.Name, value)
	return nil
}
