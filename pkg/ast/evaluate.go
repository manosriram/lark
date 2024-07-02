package ast

import (
	"fmt"
	token "lark/pkg/token"
)

var store map[string]interface{}

func Evaluate(s Statement) interface{} {
	// fmt.Printf("%T\n", s.Node)
	return Visit(s.Node)
}

func Visit(node interface{}) interface{} {
	nodeType := fmt.Sprintf("%T", node)
	// fmt.Println("visiting ", node)
	switch nodeType {
	case "ast.Number":
		return node.(Literal)
	case "ast.BinOP":
		binNode := node.(BinOP)
		left := Visit(binNode.left)
		right := Visit(binNode.right)

		switch left.(Literal).Type {
		case token.FLOAT64:
			switch binNode.op {
			case token.PLUS:
				return left.(Literal).Value.(float64) + right.(Literal).Value.(float64)
			case token.MULTIPLY:
				return left.(Literal).Value.(float64) * right.(Literal).Value.(float64)
			case token.DIVIDE:
				return left.(Literal).Value.(float64) / right.(Literal).Value.(float64)
			case token.MINUS:
				return left.(Literal).Value.(float64) - right.(Literal).Value.(float64)
			}
		case token.INTEGER:
			switch binNode.op {
			case token.PLUS:
				return left.(Literal).Value.(int) + right.(Literal).Value.(int)
			case token.MULTIPLY:
				return left.(Literal).Value.(int) * right.(Literal).Value.(int)
			case token.DIVIDE:
				return left.(Literal).Value.(int) / right.(Literal).Value.(int)
			case token.MINUS:
				return left.(Literal).Value.(int) - right.(Literal).Value.(int)
			}
		}
	}
	return node
}
