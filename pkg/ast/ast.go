package ast

import (
	types "lark/pkg/types"
	"log"
)

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

func (a *AstBuilder) peek(forward int) types.Token {
	for a.CurrentTokenPointer < len(a.tokens) {
		return a.tokens[a.CurrentTokenPointer+forward]
	}
	return types.Token{}
}

func (a *AstBuilder) eat(t types.TOKEN_TYPE) bool {
	if a.CurrentTokenPointer < len(a.tokens) && a.getCurrentToken().TokenType == t {
		a.CurrentTokenPointer++
		return true
	} else {
		log.Fatalf("syntax error: expected %s\n", t)
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
	// z := a.getCurrentToken().TokenType
	switch a.getCurrentToken().TokenType {
	case types.TRUE, types.FALSE, types.NOT:
		op := a.getCurrentToken().TokenType
		a.eat(op)
		right := a.Expr()
		left = types.BinOP{Left: left, Right: right, Op: op}

	case types.PLUS, types.MINUS, types.EQUALS, types.GREATER, types.GREATER_OR_EQUAL, types.LESSER, types.LESSER_OR_EQUAL, types.NOT_EQUAL:
		for a.getCurrentToken().TokenType == types.PLUS || a.getCurrentToken().TokenType == types.MINUS || a.getCurrentToken().TokenType == types.EQUALS || a.getCurrentToken().TokenType == types.GREATER || a.getCurrentToken().TokenType == types.GREATER_OR_EQUAL || a.getCurrentToken().TokenType == types.LESSER || a.getCurrentToken().TokenType == types.LESSER_OR_EQUAL || a.getCurrentToken().TokenType == types.NOT_EQUAL {
			op := a.getCurrentToken().TokenType
			a.eat(op)
			right := a.Term()
			left = types.BinOP{Left: left, Right: right, Op: op}
		}
	case types.ASSIGN, types.LOCAL:
		assignType := types.GLOBAL_ASSIGN
		if a.getCurrentToken().TokenType == types.LOCAL {
			a.eat(types.LOCAL)
			assignType = types.LOCAL_ASSIGN
			left = a.Term()
		}

		a.eat(types.ASSIGN)
		right := a.Expr()
		switch right.(type) {
		case types.FunctionCall:
			break
		default:
			a.eat(types.SEMICOLON)
		}
		return types.Assign{Id: left, Value: right, Type: assignType}
	case types.FUNCTION_CALL:
		fn := types.FunctionCall{Name: a.getCurrentToken().Value.(types.Literal).Value.(string)}
		a.eat(types.FUNCTION_CALL)
		for a.getCurrentToken().TokenType == types.FUNCTION_ARGUMENT {
			fn.Arguments = append(fn.Arguments, a.getCurrentToken().Value)
			a.eat(types.FUNCTION_ARGUMENT)
		}
		a.eat(types.SEMICOLON)
		return fn
	case types.FUNCTION:
		a.eat(types.FUNCTION)
		functionName := a.getCurrentToken().Value.(types.Literal).Value
		a.eat(types.ID)
		function := types.Function{
			Name: functionName.(string),
		}
		// a.eat(types.FUNCTION_ARGUMENT_OPEN)
		for a.getCurrentToken().TokenType == types.FUNCTION_ARGUMENT {
			function.Arguments = append(function.Arguments, a.getCurrentToken().Value)
			a.eat(types.FUNCTION_ARGUMENT)
		}
		// a.eat(types.FUNCTION_ARGUMENT_CLOSE)
		a.eat(types.FUNCTION_OPEN)
		for a.getCurrentToken().TokenType != types.FUNCTION_RETURN && a.getCurrentToken().TokenType != types.FUNCTION_CLOSE {
			node := a.Expr()
			function.Children = append(function.Children, node)
		}
		if a.getCurrentToken().TokenType == types.FUNCTION_RETURN {
			a.eat(types.FUNCTION_RETURN)
			function.ReturnExpression = a.Expr()
			a.eat(types.SEMICOLON)
		}
		a.eat(types.FUNCTION_CLOSE)
		return function

	case types.SWAP:
		a.eat(types.SWAP)
		right := a.Expr()
		a.eat(types.SEMICOLON)
		return types.Swap{Left: left, Right: right}
	case types.IF:
		a.eat(types.IF)
		condition := a.Expr()
		ifStatement := types.IfElseStatement{Condition: condition}

		a.eat(types.LPAREN)
		for a.getCurrentToken().TokenType != types.RPAREN {
			node := a.Expr()
			ifStatement.IfChildren = append(ifStatement.IfChildren, node)
		}
		a.eat(types.RPAREN)
		if a.getCurrentToken().TokenType == types.ELSE {
			a.eat(types.ELSE)
			a.eat(types.LPAREN)
			for a.getCurrentToken().TokenType != types.RPAREN {
				node := a.Expr()
				ifStatement.ElseChildren = append(ifStatement.ElseChildren, node)
			}
			a.eat(types.RPAREN)
		}
		return ifStatement
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
	case types.NOT:
		a.eat(types.NOT)
		if a.peek(0).Value.(types.Literal).Type != types.BOOLEAN {
			log.Fatalf("unexpected token")
		}
		return types.UnaryOP{Left: types.NOT, Right: a.Expr()}
	case types.COMMENT:
		a.eat(types.COMMENT)
	case types.MINUS:
		a.eat(types.MINUS)
		right := a.Expr()
		return types.UnaryOP{Left: types.MINUS, Right: right}
	case types.LITERAL:
		a.eat(types.LITERAL)
		return types.Literal{Value: a.tokens[c].Value, Type: a.tokens[c].LiteralType}
	case types.ID:
		a.eat(types.ID)
		return types.Id{Name: a.tokens[c].Value.(types.Literal).Value.(string)}
	case types.LOCAL:
		// a.eat(types.LOCAL)
		return types.Literal{Value: "local", Type: types.KEYWORD}
	case types.SEMICOLON:
		a.eat(types.SEMICOLON)
	case types.LBRACE:
		a.eat(types.LBRACE)
		expr := a.Expr()
		a.eat(types.RBRACE)
		return expr
	case types.EQUALS, types.GREATER, types.GREATER_OR_EQUAL, types.LESSER, types.LESSER_OR_EQUAL, types.NOT_EQUAL:
		a.eat(a.getCurrentToken().TokenType)
		right := a.Expr()
		return right
	}
	return nil
}
