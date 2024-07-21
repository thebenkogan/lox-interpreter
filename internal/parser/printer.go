package parser

import (
	"fmt"
	"math"
	"strconv"
)

func (e *ExpressionLiteral) String() string {
	if e.Literal == nil {
		return "nil"
	}
	if n, ok := e.Literal.(float64); ok {
		if math.Floor(n) == n {
			return fmt.Sprintf("%.1f", n)
		}
		return strconv.FormatFloat(n, 'f', -1, 64)
	}
	return fmt.Sprintf("%v", e.Literal)
}

func (e *ExpressionGroup) String() string {
	return fmt.Sprintf("(group %s)", e.Child.String())
}

func (e *ExpressionUnary) String() string {
	switch e.Operator {
	case UnaryOperatorMinus:
		return fmt.Sprintf("(- %s)", e.Child.String())
	case UnaryOperatorBang:
		return fmt.Sprintf("(! %s)", e.Child.String())
	}
	panic("Unknown unary operator")
}

func (e *ExpressionBinary) String() string {
	var op string
	switch e.Operator {
	case BinaryOperatorMultiply:
		op = "*"
	case BinaryOperatorDivide:
		op = "/"
	case BinaryOperatorAdd:
		op = "+"
	case BinaryOperatorSubtract:
		op = "-"
	case BinaryOperatorGreater:
		op = ">"
	case BinaryOperatorGreaterEqual:
		op = ">="
	case BinaryOperatorLess:
		op = "<"
	case BinaryOperatorLessEqual:
		op = "<="
	}
	return fmt.Sprintf("(%s %s %s)", op, e.Left.String(), e.Right.String())
}
