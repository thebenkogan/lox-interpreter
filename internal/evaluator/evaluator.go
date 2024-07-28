package evaluator

func (e *ExpressionLiteral) Evaluate(env *Environment) (Value, *RuntimeError) {
	return &ValueLiteral{Literal: e.Literal}, nil
}

func (e *ExpressionGroup) Evaluate(env *Environment) (Value, *RuntimeError) {
	return e.Child.Evaluate(env)
}

func evaluateToLiteral(exp Expression, env *Environment) (*ValueLiteral, *RuntimeError) {
	val, err := exp.Evaluate(env)
	if err != nil {
		return nil, err
	}
	literal, ok := val.(*ValueLiteral)
	if !ok {
		return nil, NewRuntimeError("Expected literal")
	}
	return literal, nil
}

func (e *ExpressionUnary) Evaluate(env *Environment) (Value, *RuntimeError) {
	child, err := e.Child.Evaluate(env)
	if err != nil {
		return nil, err
	}
	switch e.Operator {
	case UnaryOperatorBang:
		return &ValueLiteral{Literal: !child.Bool()}, nil
	case UnaryOperatorMinus:
		val, ok := child.(*ValueLiteral)
		if !ok {
			return nil, NewRuntimeError("Expected number after '-'")
		}
		n, ok := val.Literal.(float64)
		if !ok {
			return nil, NewRuntimeError("Expected number after '-'")
		}
		return &ValueLiteral{Literal: -n}, nil
	}
	panic("Unknown unary operator")
}

func getNums(left, right *ValueLiteral) (float64, float64, *RuntimeError) {
	leftNum, ok := left.Literal.(float64)
	if !ok {
		return 0, 0, NewRuntimeError("Expected number")
	}
	rightNum, ok := right.Literal.(float64)
	if !ok {
		return 0, 0, NewRuntimeError("Expected number")
	}
	return leftNum, rightNum, nil
}

func (e *ExpressionBinary) Evaluate(env *Environment) (Value, *RuntimeError) {
	if e.Operator == BinaryOperatorAnd {
		return evalAnd(e.Left, e.Right, env)
	}
	if e.Operator == BinaryOperatorOr {
		return evalOr(e.Left, e.Right, env)
	}
	left, err := evaluateToLiteral(e.Left, env)
	if err != nil {
		return nil, err
	}
	right, err := evaluateToLiteral(e.Right, env)
	if err != nil {
		return nil, err
	}
	switch e.Operator {
	case BinaryOperatorMultiply:
		leftNum, rightNum, err := getNums(left, right)
		if err != nil {
			return nil, err
		}
		return &ValueLiteral{Literal: leftNum * rightNum}, nil
	case BinaryOperatorDivide:
		leftNum, rightNum, err := getNums(left, right)
		if err != nil {
			return nil, err
		}
		if rightNum == 0 {
			return nil, NewRuntimeError("Division by zero")
		}
		return &ValueLiteral{Literal: leftNum / rightNum}, nil
	case BinaryOperatorAdd:
		leftNum, rightNum, err := getNums(left, right)
		if err != nil {
			leftStr, ok1 := left.Literal.(string)
			rightStr, ok2 := right.Literal.(string)
			if !ok1 || !ok2 {
				return nil, NewRuntimeError("Can only add numbers or strings")
			}
			return &ValueLiteral{Literal: leftStr + rightStr}, nil
		}
		return &ValueLiteral{Literal: leftNum + rightNum}, nil
	case BinaryOperatorSubtract:
		leftNum, rightNum, err := getNums(left, right)
		if err != nil {
			return nil, err
		}
		return &ValueLiteral{Literal: leftNum - rightNum}, nil
	case BinaryOperatorGreater:
		leftNum, rightNum, err := getNums(left, right)
		if err != nil {
			return nil, err
		}
		return &ValueLiteral{Literal: leftNum > rightNum}, nil
	case BinaryOperatorGreaterEqual:
		leftNum, rightNum, err := getNums(left, right)
		if err != nil {
			return nil, err
		}
		return &ValueLiteral{Literal: leftNum >= rightNum}, nil
	case BinaryOperatorLess:
		leftNum, rightNum, err := getNums(left, right)
		if err != nil {
			return nil, err
		}
		return &ValueLiteral{Literal: leftNum < rightNum}, nil
	case BinaryOperatorLessEqual:
		leftNum, rightNum, err := getNums(left, right)
		if err != nil {
			return nil, err
		}
		return &ValueLiteral{Literal: leftNum <= rightNum}, nil
	case BinaryOperatorEqual:
		return &ValueLiteral{Literal: left.Literal == right.Literal}, nil
	case BinaryOperatorNotEqual:
		return &ValueLiteral{Literal: left.Literal != right.Literal}, nil
	}
	panic("Unknown binary operator")
}

func evalOr(left, right Expression, env *Environment) (Value, *RuntimeError) {
	leftVal, err := left.Evaluate(env)
	if err != nil {
		return nil, err
	}
	if leftVal.Bool() {
		return leftVal, nil
	}
	rightVal, err := right.Evaluate(env)
	if err != nil {
		return nil, err
	}
	if rightVal.Bool() {
		return rightVal, nil
	}
	return &ValueLiteral{Literal: false}, nil
}

func evalAnd(left, right Expression, env *Environment) (Value, *RuntimeError) {
	leftVal, err := left.Evaluate(env)
	if err != nil {
		return nil, err
	}
	if !leftVal.Bool() {
		return &ValueLiteral{Literal: false}, nil
	}
	rightVal, err := right.Evaluate(env)
	if err != nil {
		return nil, err
	}
	if !rightVal.Bool() {
		return &ValueLiteral{Literal: false}, nil
	}
	return rightVal, nil
}

func (e *ExpressionVariable) Evaluate(env *Environment) (Value, *RuntimeError) {
	return env.Get(e.Name)
}

func (e *ExpressionAssignment) Evaluate(env *Environment) (Value, *RuntimeError) {
	result, err := e.Expr.Evaluate(env)
	if err != nil {
		return nil, err
	}
	if err := env.Set(e.Name, result); err != nil {
		return nil, err
	}
	return result, nil
}
