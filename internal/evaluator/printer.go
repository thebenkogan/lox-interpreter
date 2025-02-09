package evaluator

import (
	"fmt"
	"math"
	"strconv"
	"strings"
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
	case BinaryOperatorEqual:
		op = "=="
	case BinaryOperatorNotEqual:
		op = "!="
	case BinaryOperatorAnd:
		op = "and"
	case BinaryOperatorOr:
		op = "or"
	}
	return fmt.Sprintf("(%s %s %s)", op, e.Left.String(), e.Right.String())
}

func (e *ExpressionVariable) String() string {
	return e.Name
}

func (e *ExpressionAssignment) String() string {
	return fmt.Sprintf("(= %s %s)", e.Name, e.Expr.String())
}

func (e *ExpressionCall) String() string {
	callee := e.Callee.String()
	if len(e.Args) == 0 {
		return fmt.Sprintf("%s()", callee)
	}
	args := make([]string, 0, len(e.Args))
	for _, arg := range e.Args {
		args = append(args, arg.String())
	}
	return fmt.Sprintf("%s(%s)", callee, strings.Join(args, ", "))
}
