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
		return node.(Number).Value
	case "ast.BinOP":
		binNode := node.(BinOP)
		left := Visit(binNode.left)
		right := Visit(binNode.right)
		// fmt.Println(left.(int), binNode.op, right.(int))
		switch binNode.op {
		case token.PLUS:
			return left.(int) + right.(int)
		case token.MULTIPLY:
			return left.(int) * right.(int)
		case token.DIVIDE:
			return left.(int) / right.(int)
		case token.MINUS:
			return left.(int) - right.(int)
		}
	}
	return node
}
