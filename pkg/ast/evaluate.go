package ast

import (
	"fmt"
	token "lark/pkg/token"
	"log"
)

var store map[string]interface{}

func Evaluate(s Statement) interface{} {
	return Visit(s.Node)
}

func Visit(node interface{}) interface{} {
	nodeType := fmt.Sprintf("%T", node)
	switch nodeType {
	case "ast.BinOP":
		left := node.(BinOP).left
		right := node.(BinOP).right
		op := node.(BinOP).op
		left = Visit(left)
		right = Visit(right)

		leftType := fmt.Sprintf("%T", left)
		rightType := fmt.Sprintf("%T", right)
		if leftType == "string" || rightType == "string" {
			log.Fatalf("operation cannot be performed on type '%s'\n", leftType)
		}
		if leftType != rightType {
			log.Fatalf("expression type mismatch\n")
		}
		if leftType == "int" {
			switch op {
			case token.PLUS:
				return left.(int) + right.(int)
			case token.MINUS:
				return left.(int) - right.(int)
			case token.MULTIPLY:
				return left.(int) * right.(int)
			case token.DIVIDE:
				return left.(int) / right.(int)
			}
		} else if leftType == "float64" {
			switch op {
			case token.PLUS:
				return left.(float64) + right.(float64)
			case token.MINUS:
				return left.(float64) - right.(float64)
			case token.MULTIPLY:
				return left.(float64) * right.(float64)
			case token.DIVIDE:
				return left.(float64) / right.(float64)
			}
		}

	case "ast.Assign":
		n := node.(Assign)
		return Visit(n.Value)
	case "ast.Literal":
		nodeValue := node.(Literal).Value
		nodeType := fmt.Sprintf("%T", nodeValue)
		switch nodeType {
		case "int":
			return nodeValue.(int)
		case "float64":
			return nodeValue.(float64)
		case "string":
			return nodeValue.(string)
		}
	}
	return node
}
