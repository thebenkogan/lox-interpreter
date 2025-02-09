package evaluator

import "io"

type Expression interface {
	String() string
	Evaluate(env *Environment, output io.Writer) (Value, *RuntimeError)
}

type ExpressionLiteral struct {
	Literal any // number, string, bool, nil
}

type ExpressionGroup struct {
	Child Expression
}

type UnaryOperator int

const (
	UnaryOperatorBang UnaryOperator = iota
	UnaryOperatorMinus
)

type ExpressionUnary struct {
	Operator UnaryOperator
	Child    Expression
}

type BinaryOperator int

const (
	BinaryOperatorMultiply BinaryOperator = iota
	BinaryOperatorDivide
	BinaryOperatorAdd
	BinaryOperatorSubtract
	BinaryOperatorGreater
	BinaryOperatorGreaterEqual
	BinaryOperatorLess
	BinaryOperatorLessEqual
	BinaryOperatorEqual
	BinaryOperatorNotEqual
	BinaryOperatorAnd
	BinaryOperatorOr
)

type ExpressionBinary struct {
	Operator BinaryOperator
	Left     Expression
	Right    Expression
}

type ExpressionVariable struct {
	Name string
}

type ExpressionAssignment struct {
	Name string
	Expr Expression
}

type ExpressionCall struct {
	Callee Expression
	Args   []Expression
}
