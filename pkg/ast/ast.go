package ast

import (
	"fmt"
	token "lark/pkg/token"
)

var current = 0

type AstNode struct {
	Value    interface{}
	NodeType token.TOKEN_TYPE
}

type BinOP struct {
	left  interface{}
	right interface{}
	op    token.TOKEN_TYPE
}

type AstBuilder struct {
	tokens              []token.Token
	currentTokenPointer int
}

func NewAstBuilder(tokens []token.Token) AstBuilder {
	return AstBuilder{
		tokens:              tokens,
		currentTokenPointer: 0,
	}
}

func (a *AstBuilder) eat() {
	if current < len(a.tokens) {
		current++
	}
}

func (a *AstBuilder) Parse() interface{} {
	return a.Expr()
}

func (a *AstBuilder) Expr() interface{} {
	fmt.Println(a.tokens[current].TokenType)
	if current >= len(a.tokens) {
		return nil
	}
	left := a.Term()
	if a.tokens[current].TokenType == token.PLUS {
		for a.tokens[current].TokenType == token.PLUS {
			a.eat()
			right := a.Term()
			left = BinOP{left: left, right: right, op: token.PLUS}
		}
	}
	if a.tokens[current].TokenType == token.MINUS {
		for a.tokens[current].TokenType == token.MINUS {
			a.eat()
			right := a.Term()
			left = BinOP{left: left, right: right, op: token.MINUS}
		}
	}
	return left
}

func (a *AstBuilder) Term() interface{} {
	if current >= len(a.tokens) {
		return nil
	}
	left := a.Eval()
	if a.tokens[current].TokenType == token.MULTIPLY {
		for a.tokens[current].TokenType == token.MULTIPLY {
			a.eat()
			right := a.Term()
			left = BinOP{left: left, right: right, op: token.MULTIPLY}
		}
	}
	if a.tokens[current].TokenType == token.DIVIDE {
		for a.tokens[current].TokenType == token.DIVIDE {
			a.eat()
			right := a.Term()
			left = BinOP{left: left, right: right, op: token.DIVIDE}
		}
	}
	return left
}

func (a *AstBuilder) Eval() interface{} {
	c := current
	a.eat()
	return AstNode{Value: a.tokens[c].Value, NodeType: a.tokens[c].TokenType}
}
