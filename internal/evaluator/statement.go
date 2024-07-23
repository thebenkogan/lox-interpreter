package evaluator

import "fmt"

type Statement interface {
	String() string
	Execute() *RuntimeError
}

type ExpressionStatement struct {
	Expression Expression
}

func (e *ExpressionStatement) String() string {
	return fmt.Sprintf("expr %s", e.Expression.String())
}

func (e *ExpressionStatement) Execute() *RuntimeError {
	_, err := e.Expression.Evaluate()
	return err
}

type PrintStatement struct {
	Expression Expression
}

func (e *PrintStatement) String() string {
	return fmt.Sprintf("print %s", e.Expression.String())
}

func (e *PrintStatement) Execute() *RuntimeError {
	result, err := e.Expression.Evaluate()
	if err != nil {
		return err
	}
	fmt.Println(result)
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

func (e *VarStatement) Execute() *RuntimeError {
	panic("todo")
}
