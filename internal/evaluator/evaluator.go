package evaluator

func (e *ExpressionLiteral) Evaluate(env *Environment) (any, *RuntimeError) {
	return e.Literal, nil
}

func (e *ExpressionGroup) Evaluate(env *Environment) (any, *RuntimeError) {
	return e.Child.Evaluate(env)
}

func toBool(value any) bool {
	return value != nil && value != false
}

func (e *ExpressionUnary) Evaluate(env *Environment) (any, *RuntimeError) {
	child, err := e.Child.Evaluate(env)
	if err != nil {
		return nil, err
	}
	switch e.Operator {
	case UnaryOperatorBang:
		return !toBool(child), nil
	case UnaryOperatorMinus:
		n, ok := child.(float64)
		if !ok {
			return nil, NewRuntimeError("Expected number after '-'")
		}
		return -n, nil
	}
	panic("Unknown unary operator")
}

func getNums(left, right any) (float64, float64, *RuntimeError) {
	leftNum, ok := left.(float64)
	if !ok {
		return 0, 0, NewRuntimeError("Expected number")
	}
	rightNum, ok := right.(float64)
	if !ok {
		return 0, 0, NewRuntimeError("Expected number")
	}
	return leftNum, rightNum, nil
}

func (e *ExpressionBinary) Evaluate(env *Environment) (any, *RuntimeError) {
	if e.Operator == BinaryOperatorAnd {
		return evalAnd(e.Left, e.Right, env)
	}
	if e.Operator == BinaryOperatorOr {
		return evalOr(e.Left, e.Right, env)
	}
	left, err := e.Left.Evaluate(env)
	if err != nil {
		return nil, err
	}
	right, err := e.Right.Evaluate(env)
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
			return nil, NewRuntimeError("Division by zero")
		}
		return leftNum / rightNum, nil
	case BinaryOperatorAdd:
		leftNum, rightNum, err := getNums(left, right)
		if err != nil {
			leftStr, ok1 := left.(string)
			rightStr, ok2 := right.(string)
			if !ok1 || !ok2 {
				return nil, NewRuntimeError("Can only add numbers or strings")
			}
			return leftStr + rightStr, nil
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

func evalOr(left, right Expression, env *Environment) (any, *RuntimeError) {
	leftVal, err := left.Evaluate(env)
	if err != nil {
		return nil, err
	}
	if toBool(leftVal) {
		return leftVal, nil
	}
	rightVal, err := right.Evaluate(env)
	if err != nil {
		return nil, err
	}
	if toBool(rightVal) {
		return rightVal, nil
	}
	return false, nil
}

func evalAnd(left, right Expression, env *Environment) (any, *RuntimeError) {
	leftVal, err := left.Evaluate(env)
	if err != nil {
		return nil, err
	}
	if !toBool(leftVal) {
		return false, nil
	}
	rightVal, err := right.Evaluate(env)
	if err != nil {
		return nil, err
	}
	if !toBool(rightVal) {
		return false, nil
	}
	return rightVal, nil
}

func (e *ExpressionVariable) Evaluate(env *Environment) (any, *RuntimeError) {
	return env.Get(e.Name)
}

func (e *ExpressionAssignment) Evaluate(env *Environment) (any, *RuntimeError) {
	result, err := e.Expr.Evaluate(env)
	if err != nil {
		return nil, err
	}
	env.Set(e.Name, result)
	return result, nil
}
