package parser

import (
	"fmt"
	"math"
	"strconv"
)

func (e *Expression) String() string {
	switch e.Type {
	case ExpressionTypeLiteral:
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
	default:
		return "TODO"
	}
}
