package ast

import (
	"fmt"
	types "lark/pkg/types"
	"log"
)

var current = 0

type AstBuilder struct {
	tokens              []types.Token
	CurrentTokenPointer int
}

func NewAstBuilder(typess []types.Token) *AstBuilder {
	return &AstBuilder{
		tokens:              typess,
		CurrentTokenPointer: 0,
	}
}

func (a *AstBuilder) getCurrentToken() types.Token {
	if a.CurrentTokenPointer < len(a.tokens) {
		return a.tokens[a.CurrentTokenPointer]
	}
	return types.Token{}
}

func (a *AstBuilder) eat(t types.TOKEN_TYPE) bool {
	if a.CurrentTokenPointer < len(a.tokens) && a.getCurrentToken().TokenType == t {
		a.CurrentTokenPointer++
		return true
	} else {
		log.Fatalf("syntax error\n")
	}
	return false
}

func (a *AstBuilder) Parse() types.Node {
	return a.Expr()
}

func (a *AstBuilder) Expr() types.Node {
	if a.CurrentTokenPointer >= len(a.tokens) {
		return nil
	}
	left := a.Term()
	switch a.getCurrentToken().TokenType {
	case types.GREATER, types.GREATER_OR_EQUAL, types.LESSER, types.LESSER_OR_EQUAL, types.EQUALS, types.NOT_EQUAL:
		op := a.getCurrentToken().TokenType
		a.eat(op)
		right := a.Term()
		left = types.BinOP{Left: left, Right: right, Op: op}
	case types.PLUS, types.MINUS:
		for a.getCurrentToken().TokenType == types.PLUS || a.getCurrentToken().TokenType == types.MINUS {
			op := a.getCurrentToken().TokenType
			a.eat(op)
			right := a.Term()
			left = types.BinOP{Left: left, Right: right, Op: op}
		}
	case types.ASSIGN:
		a.eat(types.ASSIGN)
		right := a.Expr()
		a.eat(types.SEMICOLON)
		return types.Assign{Id: left, Value: right}
	}
	return left
}

func (a *AstBuilder) Term() types.Node {
	left := a.Factor()
	for a.getCurrentToken().TokenType == types.MULTIPLY || a.getCurrentToken().TokenType == types.DIVIDE {
		op := a.getCurrentToken().TokenType
		a.eat(op)
		right := a.Factor()
		left = types.BinOP{Left: left, Right: right, Op: op}
	}
	return left
}

func (a *AstBuilder) Factor() types.Node {
	c := a.CurrentTokenPointer
	switch a.tokens[c].TokenType {
	case types.MINUS:
		a.eat(types.MINUS)
		return types.UnaryOP{Left: types.MINUS, Right: a.Expr()}
	case types.LITERAL:
		a.eat(types.LITERAL)
		return types.Literal{Value: a.tokens[c].Value, Type: a.tokens[c].LiteralType}
	case types.ID:
		a.eat(types.ID)
		return types.Id{Name: a.tokens[c].Value.(types.Literal).Value.(string)}
	case types.SEMICOLON:
		a.eat(types.SEMICOLON)
	case types.LBRACE:
		a.eat(types.LBRACE)
		expr := a.Expr()
		a.eat(types.RBRACE)
		return expr
		// case types.EQUALS:
		// a.eat(types.EQUALS)
		// case types.GREATER:
		// fmt.Println("got ", a.tokens[c].TokenType)
		// a.eat(types.GREATER)
		// return types.UnaryOP{Left: types.GREATER}
		// case types.GREATER_OR_EQUAL:
		// a.eat(types.GREATER_OR_EQUAL)
		// case types.LESSER:
		// a.eat(types.LESSER)
		// case types.LESSER_OR_EQUAL:
		// a.eat(types.LESSER_OR_EQUAL)
		// case types.NOT:
		// a.eat(types.NOT)
		// case types.NOT_EQUAL:
		// a.eat(types.NOT_EQUAL)
	}
	fmt.Println("received nil ", a.tokens[c].TokenType == types.GREATER)
	return nil
}
