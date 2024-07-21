package parser

import (
	"errors"
	"strconv"

	"github.com/codecrafters-io/interpreter-starter-go/internal/lexer"
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

type Expression interface {
	String() string
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
)

type ExpressionBinary struct {
	Operator BinaryOperator
	Left     Expression
	Right    Expression
}

func Parse(tokens []lexer.Token) (Expression, error) {
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

func (p *parser) expression() (Expression, error) {
	return p.comparison()
}

// comparison     → term ( ( ">" | ">=" | "<" | "<=" ) term )* ;

func (p *parser) comparison() (Expression, error) {
	expr, err := p.term()
	if err != nil {
		return nil, err
	}
	for p.advanceMatch(lexer.TokenTypeGreater, lexer.TokenTypeGreaterEqual, lexer.TokenTypeLess, lexer.TokenTypeLessEqual) {
		var operator BinaryOperator
		switch p.previous().Type {
		case lexer.TokenTypeGreater:
			operator = BinaryOperatorGreater
		case lexer.TokenTypeGreaterEqual:
			operator = BinaryOperatorGreaterEqual
		case lexer.TokenTypeLess:
			operator = BinaryOperatorLess
		case lexer.TokenTypeLessEqual:
			operator = BinaryOperatorLessEqual
		}
		right, err := p.factor()
		if err != nil {
			return nil, err
		}
		expr = &ExpressionBinary{Operator: operator, Left: expr, Right: right}
	}
	return expr, nil
}

// term           → factor ( ( "-" | "+" ) factor )* ;

func (p *parser) term() (Expression, error) {
	expr, err := p.factor()
	if err != nil {
		return nil, err
	}
	for p.advanceMatch(lexer.TokenTypeMinus, lexer.TokenTypePlus) {
		operator := BinaryOperatorAdd
		if p.previous().Type == lexer.TokenTypeMinus {
			operator = BinaryOperatorSubtract
		}
		right, err := p.factor()
		if err != nil {
			return nil, err
		}
		expr = &ExpressionBinary{Operator: operator, Left: expr, Right: right}
	}
	return expr, nil
}

// factor         → unary ( ( "/" | "*" ) unary )* ;

func (p *parser) factor() (Expression, error) {
	expr, err := p.unary()
	if err != nil {
		return nil, err
	}
	for p.advanceMatch(lexer.TokenTypeSlash, lexer.TokenTypeStar) {
		operator := BinaryOperatorMultiply
		if p.previous().Type == lexer.TokenTypeSlash {
			operator = BinaryOperatorDivide
		}
		right, err := p.unary()
		if err != nil {
			return nil, err
		}
		expr = &ExpressionBinary{Operator: operator, Left: expr, Right: right}
	}
	return expr, nil
}

// unary          → ( "!" | "-" ) unary
//                | primary ;

func (p *parser) unary() (Expression, error) {
	if p.advanceMatch(lexer.TokenTypeMinus, lexer.TokenTypeBang) {
		operator := UnaryOperatorBang
		if p.previous().Type == lexer.TokenTypeMinus {
			operator = UnaryOperatorMinus
		}
		child, err := p.unary()
		if err != nil {
			return nil, err
		}
		return &ExpressionUnary{Operator: operator, Child: child}, nil
	}
	return p.primary()
}

// primary        → NUMBER | STRING | "true" | "false" | "nil"
//                | "(" expression ")" ;

func (p *parser) primary() (Expression, error) {
	switch {
	case p.advanceMatch(lexer.TokenTypeFalse):
		return &ExpressionLiteral{Literal: false}, nil
	case p.advanceMatch(lexer.TokenTypeTrue):
		return &ExpressionLiteral{Literal: true}, nil
	case p.advanceMatch(lexer.TokenTypeNil):
		return &ExpressionLiteral{Literal: nil}, nil
	case p.advanceMatch(lexer.TokenTypeNumber):
		n, _ := strconv.ParseFloat(p.previous().Literal, 64)
		return &ExpressionLiteral{Literal: n}, nil
	case p.advanceMatch(lexer.TokenTypeString):
		return &ExpressionLiteral{Literal: p.previous().Literal}, nil
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
		return &ExpressionGroup{Child: expr}, nil
	}
	return nil, nil
}
