package ast

import token "lark/pkg/token"

type AstBuilder struct {
	tokens              []token.Token
	CurrentTokenPointer int
}

type Statement struct {
	Node          interface{}
	StatementType token.STATEMENT_TYPE
}

type Expression struct {
	Node           interface{}
	ExpressionType token.EXPRESSION_TYPE
}

type Assign struct {
	Statement
	Id    interface{}
	Value interface{}
}

type Literal struct {
	Value interface{}
	Type  token.LITERAL_TYPE
}

type Id struct {
	Name  string
	Value interface{}
}

type UnaryOP struct {
	left  token.TOKEN_TYPE
	right interface{}
}

type BinOP struct {
	left  interface{}
	right interface{}
	op    token.TOKEN_TYPE
}
