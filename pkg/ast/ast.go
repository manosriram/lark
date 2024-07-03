package ast

import (
	"fmt"
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

type Literal struct {
	Value interface{}
	Type  token.LITERAL_TYPE
}

type Number struct {
	Value interface{}
	Type  token.LITERAL_TYPE
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
	case token.PLUS:
		for a.getCurrentToken().TokenType == token.PLUS {
			a.eat(token.PLUS)
			right := a.Expr()
			left = BinOP{left: left, right: right, op: token.PLUS}
		}
	case token.MINUS:
		for a.getCurrentToken().TokenType == token.MINUS {
			a.eat(token.MINUS)
			right := a.Expr()
			left = BinOP{left: left, right: right, op: token.MINUS}
		}
	case token.EQUAL:
		a.eat(token.EQUAL)
		right := a.Expr()
		a.eat(token.SEMICOLON)
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
			a.eat(token.MULTIPLY)
			right := a.Term()
			left = BinOP{left: left, right: right, op: token.MULTIPLY}
		}
	case token.DIVIDE:
		for a.tokens[a.CurrentTokenPointer].TokenType == token.DIVIDE {
			a.eat(token.DIVIDE)
			right := a.Term()
			left = BinOP{left: left, right: right, op: token.DIVIDE}
		}
	}
	return left
}

func (a *AstBuilder) Factor() interface{} {
	c := a.CurrentTokenPointer
	switch a.tokens[c].TokenType {
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
