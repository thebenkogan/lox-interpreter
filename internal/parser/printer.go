package parser

import (
	"fmt"
	"strconv"
)

func (e *Expression) String() string {
	switch e.Type {
	case ExpressionTypeLiteral:
		if e.Literal == nil {
			return "nil"
		}
		if n, ok := e.Literal.(float64); ok {
			return strconv.FormatFloat(n, 'f', -1, 64)
		}
		return fmt.Sprintf("%v", e.Literal)
	default:
		return "TODO"
	}
}
