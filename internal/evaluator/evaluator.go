package evaluator

import "errors"

func (e *ExpressionLiteral) Evaluate() (any, error) {
	return e.Literal, nil
}

func (e *ExpressionGroup) Evaluate() (any, error) {
	return e.Child.Evaluate()
}

func toBool(value any) bool {
	return value != nil && value != false
}

func (e *ExpressionUnary) Evaluate() (any, error) {
	child, err := e.Child.Evaluate()
	if err != nil {
		return nil, err
	}
	switch e.Operator {
	case UnaryOperatorBang:
		return !toBool(child), nil
	case UnaryOperatorMinus:
		n, ok := child.(float64)
		if !ok {
			return nil, errors.New("Expected number after '-'")
		}
		return -n, nil
	}
	panic("Unknown unary operator")
}

func getNums(left, right any) (float64, float64, error) {
	leftNum, ok := left.(float64)
	if !ok {
		return 0, 0, errors.New("Expected number")
	}
	rightNum, ok := right.(float64)
	if !ok {
		return 0, 0, errors.New("Expected number")
	}
	return leftNum, rightNum, nil
}

func (e *ExpressionBinary) Evaluate() (any, error) {
	left, err := e.Left.Evaluate()
	if err != nil {
		return nil, err
	}
	right, err := e.Right.Evaluate()
	if err != nil {
		return nil, err
	}
	switch e.Operator {
	case BinaryOperatorMultiply:
		leftNum, rightNum, err := getNums(left, right)
		if err != nil {
			return nil, err
		}
		return leftNum * rightNum, nil
	case BinaryOperatorDivide:
		leftNum, rightNum, err := getNums(left, right)
		if err != nil {
			return nil, err
		}
		if rightNum == 0 {
			return nil, errors.New("Division by zero")
		}
		return leftNum / rightNum, nil
	case BinaryOperatorAdd:
		leftNum, rightNum, err := getNums(left, right)
		if err != nil {
			leftStr, ok1 := left.(string)
			rightStr, ok2 := right.(string)
			if !ok1 || !ok2 {
				return nil, errors.New("Can only add numbers or strings")
			}
			return leftStr + rightStr, err
		}
		return leftNum + rightNum, nil
	case BinaryOperatorSubtract:
		leftNum, rightNum, err := getNums(left, right)
		if err != nil {
			return nil, err
		}
		return leftNum - rightNum, nil
	case BinaryOperatorGreater:
		leftNum, rightNum, err := getNums(left, right)
		if err != nil {
			return nil, err
		}
		return leftNum > rightNum, nil
	case BinaryOperatorGreaterEqual:
		leftNum, rightNum, err := getNums(left, right)
		if err != nil {
			return nil, err
		}
		return leftNum >= rightNum, nil
	case BinaryOperatorLess:
		leftNum, rightNum, err := getNums(left, right)
		if err != nil {
			return nil, err
		}
		return leftNum < rightNum, nil
	case BinaryOperatorLessEqual:
		leftNum, rightNum, err := getNums(left, right)
		if err != nil {
			return nil, err
		}
		return leftNum <= rightNum, nil
	case BinaryOperatorEqual:
		return left == right, nil
	case BinaryOperatorNotEqual:
		return left != right, nil
	}
	panic("Unknown binary operator")
}
