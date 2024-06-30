package ast

import token "lark/pkg/token"

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

func eat(t []token.Token) {
	if current < len(t) {
		current++
	}
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

func (a *AstBuilder) Parse() interface{} {
	return a.Expr()
}

func (a *AstBuilder) Expr() interface{} {
	if current >= len(a.tokens) {
		return nil
	}
	left := a.Term()
	for a.tokens[current].TokenType == token.PLUS {
		eat(a.tokens)
		right := a.Term()
		left = BinOP{left: left, right: right, op: token.PLUS}
	}
	return left
}

func (a *AstBuilder) Term() interface{} {
	if current >= len(a.tokens) {
		return nil
	}
	left := a.Eval()
	for a.tokens[current].TokenType == token.MULTIPLY {
		eat(a.tokens)
		right := a.Eval()
		left = BinOP{left: left, right: right, op: token.MULTIPLY}
	}
	return left
}

func (a *AstBuilder) Eval() interface{} {
	c := current
	eat(a.tokens)
	return AstNode{Value: a.tokens[c].Value, NodeType: a.tokens[c].TokenType}
}
