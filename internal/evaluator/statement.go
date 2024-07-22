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
	return e.Expression.String()
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
