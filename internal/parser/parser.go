package parser

import (
	"strconv"

	"github.com/thebenkogan/lox-interpreter/internal/evaluator"
	"github.com/thebenkogan/lox-interpreter/internal/lexer"
)

func Parse(tokens []lexer.Token) ([]evaluator.Statement, *ParserError) {
	p := &parser{tokens: tokens}
	statements := make([]evaluator.Statement, 0)
	for !p.isAtEnd() {
		statement, err := p.statement()
		if err != nil {
			return nil, err
		}
		statements = append(statements, statement)
	}
	return statements, nil
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

func (p *parser) statement() (evaluator.Statement, *ParserError) {
	if p.advanceMatch(lexer.TokenTypeVar) {
		return p.varStatement()
	}
	if p.advanceMatch(lexer.TokenTypePrint) {
		return p.printStatement()
	}
	return p.expressionStatement()
}

// varDecl        → "var" IDENTIFIER ( "=" expression )? ";" ;

func (p *parser) varStatement() (*evaluator.VarStatement, *ParserError) {
	if !p.advanceMatch(lexer.TokenTypeIdentifier) {
		return nil, NewParserError("Expected variable name")
	}

	varStmt := &evaluator.VarStatement{Name: p.previous().Lexeme}
	if p.advanceMatch(lexer.TokenTypeEqual) {
		expr, err := p.expression()
		if err != nil {
			return nil, err
		}
		varStmt.Expr = expr
	}

	if !p.advanceMatch(lexer.TokenTypeSemicolon) {
		return nil, NewParserError("Expected semicolon after var statement")
	}

	return varStmt, nil
}

// printStmt      → "print" expression ";" ;

func (p *parser) printStatement() (*evaluator.PrintStatement, *ParserError) {
	expr, err := p.expression()
	if err != nil {
		return nil, err
	}
	if !p.advanceMatch(lexer.TokenTypeSemicolon) {
		return nil, NewParserError("Expected semicolon after print statement")
	}
	return &evaluator.PrintStatement{Expression: expr}, nil
}

// exprStmt       → expression ";" ;

func (p *parser) expressionStatement() (*evaluator.ExpressionStatement, *ParserError) {
	expr, err := p.expression()
	if err != nil {
		return nil, err
	}
	if !p.advanceMatch(lexer.TokenTypeSemicolon) {
		return nil, NewParserError("Expected semicolon after expression statement")
	}
	return &evaluator.ExpressionStatement{Expression: expr}, nil
}

// expression     → assignment ;
// assignment     → IDENTIFIER "=" assignment
//                | equality ;
// equality       → comparison ( ( "!=" | "==" ) comparison )* ;
// comparison     → term ( ( ">" | ">=" | "<" | "<=" ) term )* ;
// term           → factor ( ( "-" | "+" ) factor )* ;
// factor         → unary ( ( "/" | "*" ) unary )* ;
// unary          → ( "!" | "-" ) unary
//                | primary ;
// primary        → NUMBER | STRING | "true" | "false" | "nil"
//                | "(" expression ")" ;

func (p *parser) expression() (evaluator.Expression, *ParserError) {
	return p.assignment()
}

// assignment     → IDENTIFIER "=" assignment
//                | equality ;

func (p *parser) assignment() (evaluator.Expression, *ParserError) {
	expr, err := p.equality()
	if err != nil {
		return nil, err
	}
	if p.advanceMatch(lexer.TokenTypeEqual) {
		variable, ok := expr.(*evaluator.ExpressionVariable)
		if !ok {
			return nil, NewParserError("Can only assign to variables")
		}
		right, err := p.assignment()
		if err != nil {
			return nil, err
		}
		return &evaluator.ExpressionAssignment{Name: variable.Name, Expr: right}, nil
	}
	return expr, nil
}

// equality       → comparison ( ( "!=" | "==" ) comparison )* ;

func (p *parser) equality() (evaluator.Expression, *ParserError) {
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

func (p *parser) comparison() (evaluator.Expression, *ParserError) {
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

func (p *parser) term() (evaluator.Expression, *ParserError) {
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

func (p *parser) factor() (evaluator.Expression, *ParserError) {
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

func (p *parser) unary() (evaluator.Expression, *ParserError) {
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
//                | "(" expression ")" | IDENTIFIER ;

func (p *parser) primary() (evaluator.Expression, *ParserError) {
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
		if !p.advanceMatch(lexer.TokenTypeRightParen) {
			return nil, NewParserError("Unmatched parentheses.")
		}
		return &evaluator.ExpressionGroup{Child: expr}, nil
	case p.advanceMatch(lexer.TokenTypeIdentifier):
		return &evaluator.ExpressionVariable{Name: p.previous().Lexeme}, nil
	}
	return nil, NewParserError("Expected expression.")
}
