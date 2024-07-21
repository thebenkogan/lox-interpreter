package parser

import (
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

type ExpressionType int

const (
	ExpressionTypeLiteral ExpressionType = iota
)

type Expression struct {
	Type     ExpressionType
	Literal  any // number, string, bool, nil
	Children []Expression
}

func Parse(tokens []lexer.Token) Expression {
	p := &parser{tokens: tokens}
	return p.primary()
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

func (p *parser) primary() Expression {
	if p.advanceMatch(lexer.TokenTypeFalse) {
		return Expression{Type: ExpressionTypeLiteral, Literal: false}
	}
	if p.advanceMatch(lexer.TokenTypeTrue) {
		return Expression{Type: ExpressionTypeLiteral, Literal: true}
	}
	if p.advanceMatch(lexer.TokenTypeNil) {
		return Expression{Type: ExpressionTypeLiteral, Literal: nil}
	}
	if p.advanceMatch(lexer.TokenTypeNumber) {
		n, _ := strconv.ParseFloat(p.previous().Literal, 64)
		return Expression{Type: ExpressionTypeLiteral, Literal: n}
	}
	panic("TODO")
}
