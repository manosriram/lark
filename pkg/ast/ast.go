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

func Expr(t []token.Token) interface{} {
	if current >= len(t) {
		return nil
	}
	left := Term(t)
	for t[current].TokenType == token.PLUS {
		eat(t)
		right := Term(t)
		left = BinOP{left: left, right: right, op: token.PLUS}
	}
	return left
}

func Term(t []token.Token) interface{} {
	if current >= len(t) {
		return nil
	}
	left := Eval(t)
	for t[current].TokenType == token.MULTIPLY {
		eat(t)
		right := Eval(t)
		left = BinOP{left: left, right: right, op: token.MULTIPLY}
	}
	return left
}

func Eval(t []token.Token) interface{} {
	c := current
	eat(t)
	return AstNode{Value: t[c].Value, NodeType: t[c].TokenType}
}
