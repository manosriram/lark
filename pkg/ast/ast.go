package ast

import (
	"fmt"
	token "lark/pkg/token"
	"log"
)

var current = 0

func NewAstBuilder(tokens []token.Token) *AstBuilder {
	return &AstBuilder{
		tokens:              tokens,
		CurrentTokenPointer: 0,
	}
}

func (a *AstBuilder) getCurrentToken() token.Token {
	if a.CurrentTokenPointer < len(a.tokens) {
		return a.tokens[a.CurrentTokenPointer]
	}
	return token.Token{}
}

func (a *AstBuilder) eat(t token.TOKEN_TYPE) bool {
	if a.CurrentTokenPointer < len(a.tokens) && a.getCurrentToken().TokenType == t {
		a.CurrentTokenPointer++
		return true
	} else {
		log.Fatalf("syntax error\n")
	}
	return false
}

func (a *AstBuilder) Parse() interface{} {
	return a.Expr()
}

func (a *AstBuilder) Expr() interface{} {
	if a.CurrentTokenPointer >= len(a.tokens) {
		return nil
	}
	left := a.Term()
	switch a.getCurrentToken().TokenType {
	case token.PLUS, token.MINUS:
		for a.getCurrentToken().TokenType == token.PLUS || a.getCurrentToken().TokenType == token.MINUS {
			op := a.getCurrentToken().TokenType
			a.eat(op)
			right := a.Term()
			left = BinOP{left: left, right: right, op: op}
		}
	case token.EQUAL:
		a.eat(token.EQUAL)
		right := a.Expr()
		a.eat(token.SEMICOLON)
		// return Statement{Node: Assign{Id: left, Value: right}, StatementType: token.ASSIGN_STATEMENT}
		return Assign{Id: left, Value: right}
	}
	return left
}

func (a *AstBuilder) Term() interface{} {
	left := a.Factor()
	for a.getCurrentToken().TokenType == token.MULTIPLY || a.getCurrentToken().TokenType == token.DIVIDE {
		op := a.getCurrentToken().TokenType
		a.eat(op)
		right := a.Factor()
		left = BinOP{left: left, right: right, op: op}
	}
	return left
}

func (a *AstBuilder) Factor() interface{} {
	c := a.CurrentTokenPointer
	switch a.tokens[c].TokenType {
	case token.MINUS:
		a.eat(token.MINUS)
		// return Expression{Node: UnaryOP{left: token.MINUS, right: a.Expr()}, ExpressionType: token.UNARY_OP}
		return UnaryOP{left: token.MINUS, right: a.Expr()}
	case token.LITERAL:
		a.eat(token.LITERAL)
		return Literal{Value: a.tokens[c].Value, Type: a.tokens[c].LiteralType}
	case token.ID:
		a.eat(token.ID)
		return Id{Name: a.tokens[c].Value.(string)}
	case token.SEMICOLON:
		a.eat(token.SEMICOLON)
	case token.LBRACE:
		a.eat(token.LBRACE)
		expr := a.Expr()
		a.eat(token.RBRACE)
		return expr
	}
	fmt.Println("received nil ", a.tokens[c].TokenType)
	return nil
}
