package parser

import "fmt"

func (e *Expression) String() string {
	switch e.Type {
	case ExpressionTypeLiteral:
		return fmt.Sprintf("%v", e.Literal)
	default:
		return "TODO"
	}
}
