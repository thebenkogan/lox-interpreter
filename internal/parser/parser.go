package parser

import (
	"errors"
	"strconv"

	"github.com/thebenkogan/lox-interpreter/internal/evaluator"
	"github.com/thebenkogan/lox-interpreter/internal/lexer"
)

// expression     → equality ;
// equality       → comparison ( ( "!=" | "==" ) comparison )* ;
// comparison     → term ( ( ">" | ">=" | "<" | "<=" ) term )* ;
// term           → factor ( ( "-" | "+" ) factor )* ;
// factor         → unary ( ( "/" | "*" ) unary )* ;
// unary          → ( "!" | "-" ) unary
//                | primary ;
// primary        → NUMBER | STRING | "true" | "false" | "nil"
//                | "(" expression ")" ;

func Parse(tokens []lexer.Token) (evaluator.Expression, error) {
	p := &parser{tokens: tokens}
	return p.expression()
}

type parser struct {
	tokens []lexer.Token
	index  int
}

func (p *parser) isAtEnd() bool {
	return p.peek().Type == lexer.TokenTypeEOF
}

func (p *parser) peek() lexer.Token {
	return p.tokens[p.index]
}

func (p *parser) previous() lexer.Token {
	return p.tokens[p.index-1]
}

func (p *parser) advance() lexer.Token {
	if !p.isAtEnd() {
		p.index++
	}
	return p.previous()
}

func (p *parser) check(t lexer.TokenType) bool {
	if p.isAtEnd() {
		return false
	}
	return p.peek().Type == t
}

func (p *parser) advanceMatch(types ...lexer.TokenType) bool {
	for _, t := range types {
		if p.check(t) {
			p.advance()
			return true
		}
	}
	return false
}

func (p *parser) expression() (evaluator.Expression, error) {
	return p.equality()
}

// equality       → comparison ( ( "!=" | "==" ) comparison )* ;

func (p *parser) equality() (evaluator.Expression, error) {
	expr, err := p.comparison()
	if err != nil {
		return nil, err
	}
	for p.advanceMatch(lexer.TokenTypeBangEqual, lexer.TokenTypeEqualEqual) {
		operator := evaluator.BinaryOperatorEqual
		if p.previous().Type == lexer.TokenTypeBangEqual {
			operator = evaluator.BinaryOperatorNotEqual
		}
		right, err := p.factor()
		if err != nil {
			return nil, err
		}
		expr = &evaluator.ExpressionBinary{Operator: operator, Left: expr, Right: right}
	}
	return expr, nil
}

// comparison     → term ( ( ">" | ">=" | "<" | "<=" ) term )* ;

func (p *parser) comparison() (evaluator.Expression, error) {
	expr, err := p.term()
	if err != nil {
		return nil, err
	}
	for p.advanceMatch(lexer.TokenTypeGreater, lexer.TokenTypeGreaterEqual, lexer.TokenTypeLess, lexer.TokenTypeLessEqual) {
		var operator evaluator.BinaryOperator
		switch p.previous().Type {
		case lexer.TokenTypeGreater:
			operator = evaluator.BinaryOperatorGreater
		case lexer.TokenTypeGreaterEqual:
			operator = evaluator.BinaryOperatorGreaterEqual
		case lexer.TokenTypeLess:
			operator = evaluator.BinaryOperatorLess
		case lexer.TokenTypeLessEqual:
			operator = evaluator.BinaryOperatorLessEqual
		}
		right, err := p.factor()
		if err != nil {
			return nil, err
		}
		expr = &evaluator.ExpressionBinary{Operator: operator, Left: expr, Right: right}
	}
	return expr, nil
}

// term           → factor ( ( "-" | "+" ) factor )* ;

func (p *parser) term() (evaluator.Expression, error) {
	expr, err := p.factor()
	if err != nil {
		return nil, err
	}
	for p.advanceMatch(lexer.TokenTypeMinus, lexer.TokenTypePlus) {
		operator := evaluator.BinaryOperatorAdd
		if p.previous().Type == lexer.TokenTypeMinus {
			operator = evaluator.BinaryOperatorSubtract
		}
		right, err := p.factor()
		if err != nil {
			return nil, err
		}
		expr = &evaluator.ExpressionBinary{Operator: operator, Left: expr, Right: right}
	}
	return expr, nil
}

// factor         → unary ( ( "/" | "*" ) unary )* ;

func (p *parser) factor() (evaluator.Expression, error) {
	expr, err := p.unary()
	if err != nil {
		return nil, err
	}
	for p.advanceMatch(lexer.TokenTypeSlash, lexer.TokenTypeStar) {
		operator := evaluator.BinaryOperatorMultiply
		if p.previous().Type == lexer.TokenTypeSlash {
			operator = evaluator.BinaryOperatorDivide
		}
		right, err := p.unary()
		if err != nil {
			return nil, err
		}
		expr = &evaluator.ExpressionBinary{Operator: operator, Left: expr, Right: right}
	}
	return expr, nil
}

// unary          → ( "!" | "-" ) unary
//                | primary ;

func (p *parser) unary() (evaluator.Expression, error) {
	if p.advanceMatch(lexer.TokenTypeMinus, lexer.TokenTypeBang) {
		operator := evaluator.UnaryOperatorBang
		if p.previous().Type == lexer.TokenTypeMinus {
			operator = evaluator.UnaryOperatorMinus
		}
		child, err := p.unary()
		if err != nil {
			return nil, err
		}
		return &evaluator.ExpressionUnary{Operator: operator, Child: child}, nil
	}
	return p.primary()
}

// primary        → NUMBER | STRING | "true" | "false" | "nil"
//                | "(" expression ")" ;

func (p *parser) primary() (evaluator.Expression, error) {
	switch {
	case p.advanceMatch(lexer.TokenTypeFalse):
		return &evaluator.ExpressionLiteral{Literal: false}, nil
	case p.advanceMatch(lexer.TokenTypeTrue):
		return &evaluator.ExpressionLiteral{Literal: true}, nil
	case p.advanceMatch(lexer.TokenTypeNil):
		return &evaluator.ExpressionLiteral{Literal: nil}, nil
	case p.advanceMatch(lexer.TokenTypeNumber):
		n, _ := strconv.ParseFloat(p.previous().Literal, 64)
		return &evaluator.ExpressionLiteral{Literal: n}, nil
	case p.advanceMatch(lexer.TokenTypeString):
		return &evaluator.ExpressionLiteral{Literal: p.previous().Literal}, nil
	case p.advanceMatch(lexer.TokenTypeLeftParen):
		expr, err := p.expression()
		if err != nil {
			return nil, err
		}
		if expr == nil {
			return nil, errors.New("Expected expression after '('")
		}
		if !p.advanceMatch(lexer.TokenTypeRightParen) {
			return nil, errors.New("Unmatched parentheses.")
		}
		return &evaluator.ExpressionGroup{Child: expr}, nil
	}
	return nil, errors.New("Expected expression.")
}
