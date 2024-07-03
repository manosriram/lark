package ast

import (
	"fmt"
	token "lark/pkg/token"
	"log"
)

var store map[string]interface{}

type RealNumber interface {
	int | float64
}

func performOperation[T RealNumber](left, right T, op token.TOKEN_TYPE) T {
	switch op {
	case token.PLUS:
		return left + right
	case token.MINUS:
		return left - right
	case token.MULTIPLY:
		return left * right
	case token.DIVIDE:
		return left / right
	default:
		panic(fmt.Sprintf("unsupported operation: %v", op))
	}
}
func Evaluate(s Statement) interface{} {
	return Visit(s.Node)
}

func Visit(node interface{}) interface{} {
	switch n := node.(type) {
	case UnaryOP:
		op := n.left
		right := Visit(n.right)

		switch right.(type) {
		case float64:
			return performOperation(0, right.(float64), op)
		case int:
			return performOperation(0, right.(int), op)
		}
	case BinOP:
		left := Visit(n.left)
		right := Visit(n.right)

		// TODO: make this better
		leftType := fmt.Sprintf("%T", left)
		rightType := fmt.Sprintf("%T", right)
		if leftType == "string" || rightType == "string" {
			log.Fatalf("operation cannot be performed on type '%s'\n", leftType)
		}
		if leftType != rightType {
			log.Fatalf("expression type mismatch\n")
		}
		switch left := left.(type) {
		case int:
			if right, ok := right.(int); ok {
				return performOperation(left, right, n.op)
			}
		case float64:
			if right, ok := right.(float64); ok {
				return performOperation(left, right, n.op)
			}
		}

	case Assign:
		return Visit(n.Value)
	case Literal:
		nodeValue := n.Value
		switch v := nodeValue.(type) {
		case int, float64, string:
			return v
		default:
			log.Fatalf("unsupported type %s\n", v)

		}
		return node
	}
	return node
}
