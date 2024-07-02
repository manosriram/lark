package ast

import (
	token "lark/pkg/token"
	"log"
)

var current = 0

type Statement struct {
	StatementType token.STATEMENT_TYPE
	Node          interface{}
}

type Expression struct {
	Statement
}

type AstNode struct {
	Value    interface{}
	NodeType token.TOKEN_TYPE
}

type Number struct {
	Value interface{}
}

type Id struct {
	Name  string
	Value interface{}
}

type Assign struct {
	Statement
	Id    interface{}
	Value interface{}
}

type BinOP struct {
	left  interface{}
	right interface{}
	op    token.TOKEN_TYPE
}

type AstBuilder struct {
	tokens              []token.Token
	CurrentTokenPointer int
}

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

func (a *AstBuilder) eat() bool {
	if a.CurrentTokenPointer < len(a.tokens) {
		a.CurrentTokenPointer++
		return true
	}
	return false
}

func (a *AstBuilder) Parse() interface{} {
	expr := a.Expr()
	if a.getCurrentToken().TokenType != token.SEMICOLON {
		log.Fatalf("syntax error: missing ;")
	}
	a.eat()
	return expr
}

func (a *AstBuilder) Expr() interface{} {
	if a.CurrentTokenPointer >= len(a.tokens) {
		return nil
	}
	left := a.Term()
	switch a.getCurrentToken().TokenType {
	case token.PLUS:
		for a.getCurrentToken().TokenType == token.PLUS {
			a.eat()
			right := a.Expr()
			left = BinOP{left: left, right: right, op: token.PLUS}
		}
	case token.MINUS:
		for a.getCurrentToken().TokenType == token.MINUS {
			a.eat()
			right := a.Expr()
			left = BinOP{left: left, right: right, op: token.MINUS}
		}
	case token.EQUAL:
		a.eat()
		right := a.Expr()
		return Assign{Id: left, Value: right}
	}
	return left
}

func (a *AstBuilder) Term() interface{} {
	if a.CurrentTokenPointer >= len(a.tokens) {
		return nil
	}
	left := a.Factor()
	switch a.tokens[a.CurrentTokenPointer].TokenType {
	case token.MULTIPLY:
		for a.tokens[a.CurrentTokenPointer].TokenType == token.MULTIPLY {
			a.eat()
			right := a.Term()
			left = BinOP{left: left, right: right, op: token.MULTIPLY}
		}
	case token.DIVIDE:
		for a.tokens[a.CurrentTokenPointer].TokenType == token.DIVIDE {
			a.eat()
			right := a.Term()
			// left := a.Factor()
			left = BinOP{left: left, right: right, op: token.DIVIDE}
		}
	}
	return left
}

func (a *AstBuilder) Factor() interface{} {
	c := a.CurrentTokenPointer
	switch a.tokens[c].TokenType {
	case token.NUMBER:
		a.eat()
		return Number{Value: a.tokens[c].Value}
	case token.ID:
		a.eat()
		return Id{Name: a.tokens[c].Value.(string)}
	case token.LBRACE:
		a.eat()
		expr := a.Expr()
		if a.getCurrentToken().TokenType != token.RBRACE {
			log.Fatalf("syntax error: missing ) \n")
		}
		a.eat()
		return expr
	}
	return nil
}
