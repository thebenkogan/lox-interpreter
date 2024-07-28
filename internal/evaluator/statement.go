package evaluator

import (
	"fmt"
	"io"
	"strings"
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
	var value Value = &ValueLiteral{Literal: nil}
	if e.Expr != nil {
		result, err := e.Expr.Evaluate(env)
		if err != nil {
			return err
		}
		value = result
	}
	env.Declare(e.Name, value)
	return nil
}

type BlockStatement struct {
	Statements []Statement
}

func (e *BlockStatement) String() string {
	statements := make([]string, 0, len(e.Statements))
	for _, stmt := range e.Statements {
		statements = append(statements, stmt.String()+";")
	}
	return fmt.Sprintf("(block %s)", strings.Join(statements, " "))
}

func (e *BlockStatement) Execute(env *Environment, output io.Writer) *RuntimeError {
	innerEnv := env.CreateScope()
	for _, stmt := range e.Statements {
		err := stmt.Execute(innerEnv, output)
		if err != nil {
			return err
		}
	}
	return nil
}

type IfStatement struct {
	Condition Expression
	Then      *BlockStatement
	Else      *BlockStatement
}

func (e *IfStatement) String() string {
	if e.Else == nil {
		return fmt.Sprintf("if (%s) then %s", e.Condition.String(), e.Then.String())
	}
	return fmt.Sprintf("if (%s) then %s else %s", e.Condition.String(), e.Then.String(), e.Else.String())
}

func (e *IfStatement) Execute(env *Environment, output io.Writer) *RuntimeError {
	condition, err := e.Condition.Evaluate(env)
	if err != nil {
		return err
	}
	if condition.Bool() {
		return e.Then.Execute(env, output)
	} else if e.Else != nil {
		return e.Else.Execute(env, output)
	}
	return nil
}

type WhileStatement struct {
	Condition Expression
	Body      *BlockStatement
}

func (e *WhileStatement) String() string {
	return fmt.Sprintf("while (%s) then %s", e.Condition.String(), e.Body.String())
}

func (e *WhileStatement) Execute(env *Environment, output io.Writer) *RuntimeError {
	for {
		condition, err := e.Condition.Evaluate(env)
		if err != nil {
			return err
		}
		if !condition.Bool() {
			break
		}
		if err := e.Body.Execute(env, output); err != nil {
			return err
		}
	}
	return nil
}

type FunStatement struct {
	Name   string
	Body   *BlockStatement
	Params []string
}

func (e *FunStatement) String() string {
	return fmt.Sprintf("fun %s(%s) %s", e.Name, strings.Join(e.Params, ", "), e.Body.String())
}

func (e *FunStatement) Execute(env *Environment, output io.Writer) *RuntimeError {
	closure := &ValueClosure{Env: env, Body: e.Body, Params: e.Params}
	env.Declare(e.Name, closure)
	return nil
}
