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
	switch {
	case p.advanceMatch(lexer.TokenTypeFun):
		return p.funStatement()
	case p.advanceMatch(lexer.TokenTypeVar):
		return p.varStatement()
	case p.advanceMatch(lexer.TokenTypeFor):
		return p.forStatement()
	case p.advanceMatch(lexer.TokenTypeIf):
		return p.ifStatement()
	case p.advanceMatch(lexer.TokenTypePrint):
		return p.printStatement()
	case p.advanceMatch(lexer.TokenTypeWhile):
		return p.whileStatement()
	case p.advanceMatch(lexer.TokenTypeLeftBrace):
		return p.blockStatement()
	default:
		return p.expressionStatement()
	}
}

// funDecl        → "fun" function ;
// function       → IDENTIFIER "(" parameters? ")" block ;
// parameters     → IDENTIFIER ( "," IDENTIFIER )* ;

func (p *parser) funStatement() (*evaluator.FunStatement, *ParserError) {
	if !p.advanceMatch(lexer.TokenTypeIdentifier) {
		return nil, NewParserError("Expected function name")
	}
	name := p.previous().Lexeme

	if !p.advanceMatch(lexer.TokenTypeLeftParen) {
		return nil, NewParserError("Expected '(' after function name")
	}

	params := make([]string, 0)
	for p.advanceMatch(lexer.TokenTypeIdentifier) {
		params = append(params, p.previous().Lexeme)
		if !p.advanceMatch(lexer.TokenTypeComma) {
			break
		}
	}

	if !p.advanceMatch(lexer.TokenTypeRightParen) {
		return nil, NewParserError("Expected ')' after function parameters")
	}

	if !p.advanceMatch(lexer.TokenTypeLeftBrace) {
		return nil, NewParserError("Expected '{' after function parameters")
	}

	body, err := p.blockStatement()
	if err != nil {
		return nil, err
	}

	return &evaluator.FunStatement{Name: name, Params: params, Body: body}, nil
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

// forStmt        → "for" "(" ( varDecl | exprStmt | ";" )
//                  expression? ";"
//                  expression? ")" blockStmt ;
// desugar to block statement with initializer and while statement

func (p *parser) forStatement() (*evaluator.BlockStatement, *ParserError) {
	if !p.advanceMatch(lexer.TokenTypeLeftParen) {
		return nil, NewParserError("Expected '(' after 'for'")
	}

	var init evaluator.Statement
	if !p.advanceMatch(lexer.TokenTypeSemicolon) {
		statement, err := p.statement()
		if err != nil {
			return nil, err
		}
		_, isExprStmt := statement.(*evaluator.ExpressionStatement)
		_, isVarStmt := statement.(*evaluator.VarStatement)
		if !isExprStmt && !isVarStmt {
			return nil, NewParserError("Expected expression or variable declaration in for loop initializer")
		}
		init = statement
	}

	var condition evaluator.Expression
	if !p.advanceMatch(lexer.TokenTypeSemicolon) {
		cond, err := p.expression()
		if err != nil {
			return nil, err
		}
		condition = cond
		if !p.advanceMatch(lexer.TokenTypeSemicolon) {
			return nil, NewParserError("Expected ';' after for condition")
		}
	}

	var increment evaluator.Expression
	if !p.advanceMatch(lexer.TokenTypeRightParen) {
		inc, err := p.expression()
		if err != nil {
			return nil, err
		}
		increment = inc
		if !p.advanceMatch(lexer.TokenTypeRightParen) {
			return nil, NewParserError("Expected ')' after for increment")
		}
	}

	if !p.advanceMatch(lexer.TokenTypeLeftBrace) {
		return nil, NewParserError("Expected '{' after for header")
	}

	body, err := p.blockStatement()
	if err != nil {
		return nil, err
	}
	if increment != nil {
		body.Statements = append(body.Statements, &evaluator.ExpressionStatement{Expression: increment})
	}

	whileStmt := &evaluator.WhileStatement{Condition: condition, Body: body}
	if condition == nil {
		whileStmt.Condition = &evaluator.ExpressionLiteral{Literal: true}
	}
	blockStmts := make([]evaluator.Statement, 0)
	if init != nil {
		blockStmts = append(blockStmts, init)
	}
	blockStmts = append(blockStmts, whileStmt)
	return &evaluator.BlockStatement{Statements: blockStmts}, nil
}

// ifStmt         → "if" "(" expression ")" blockStmt
//                ( "else" blockStmt )? ;

func (p *parser) ifStatement() (*evaluator.IfStatement, *ParserError) {
	if !p.advanceMatch(lexer.TokenTypeLeftParen) {
		return nil, NewParserError("Expected '(' after 'if'")
	}
	condition, err := p.expression()
	if err != nil {
		return nil, err
	}
	if !p.advanceMatch(lexer.TokenTypeRightParen) {
		return nil, NewParserError("Expected ')' after if condition")
	}
	if !p.advanceMatch(lexer.TokenTypeLeftBrace) {
		return nil, NewParserError("Expected '{' after if condition")
	}
	then, err := p.blockStatement()
	if err != nil {
		return nil, err
	}
	var elseStmt *evaluator.BlockStatement
	if p.advanceMatch(lexer.TokenTypeElse) {
		if !p.advanceMatch(lexer.TokenTypeLeftBrace) {
			return nil, NewParserError("Expected '{' after else")
		}
		elseStmt, err = p.blockStatement()
		if err != nil {
			return nil, err
		}
	}
	return &evaluator.IfStatement{Condition: condition, Then: then, Else: elseStmt}, nil
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

// whileStmt      → "while" "(" expression ")" blockStmt ;

func (p *parser) whileStatement() (*evaluator.WhileStatement, *ParserError) {
	if !p.advanceMatch(lexer.TokenTypeLeftParen) {
		return nil, NewParserError("Expected '(' after 'while'")
	}
	condition, err := p.expression()
	if err != nil {
		return nil, err
	}
	if !p.advanceMatch(lexer.TokenTypeRightParen) {
		return nil, NewParserError("Expected ')' after while condition")
	}
	if !p.advanceMatch(lexer.TokenTypeLeftBrace) {
		return nil, NewParserError("Expected '{' after while condition")
	}
	body, err := p.blockStatement()
	if err != nil {
		return nil, err
	}
	return &evaluator.WhileStatement{Condition: condition, Body: body}, nil
}

// blockStmt          → "{" declaration* "}" ;

func (p *parser) blockStatement() (*evaluator.BlockStatement, *ParserError) {
	statements := make([]evaluator.Statement, 0)
	for !p.isAtEnd() && p.peek().Type != lexer.TokenTypeRightBrace {
		statement, err := p.statement()
		if err != nil {
			return nil, err
		}
		statements = append(statements, statement)
	}
	if !p.advanceMatch(lexer.TokenTypeRightBrace) {
		return nil, NewParserError("Expected right brace after block statement")
	}
	return &evaluator.BlockStatement{Statements: statements}, nil
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
//                | logic_or ;
// logic_or       → logic_and ( "or" logic_and )* ;
// logic_and      → equality ( "and" equality )* ;
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
//                | logic_or ;

func (p *parser) assignment() (evaluator.Expression, *ParserError) {
	expr, err := p.logicOr()
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

// logic_or       → logic_and ( "or" logic_and )* ;

func (p *parser) logicOr() (evaluator.Expression, *ParserError) {
	expr, err := p.logicAnd()
	if err != nil {
		return nil, err
	}
	for p.advanceMatch(lexer.TokenTypeOr) {
		right, err := p.logicAnd()
		if err != nil {
			return nil, err
		}
		expr = &evaluator.ExpressionBinary{Operator: evaluator.BinaryOperatorOr, Left: expr, Right: right}
	}
	return expr, nil
}

// logic_and      → equality ( "and" equality )* ;

func (p *parser) logicAnd() (evaluator.Expression, *ParserError) {
	expr, err := p.equality()
	if err != nil {
		return nil, err
	}
	for p.advanceMatch(lexer.TokenTypeAnd) {
		right, err := p.equality()
		if err != nil {
			return nil, err
		}
		expr = &evaluator.ExpressionBinary{Operator: evaluator.BinaryOperatorAnd, Left: expr, Right: right}
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
