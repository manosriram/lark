package ast

import (
	"fmt"
	types "lark/pkg/types"
	"log"

	"golang.org/x/exp/constraints"
)

type Evaluator struct {
	SymbolTable map[string]interface{}
}

type Comparator interface {
	constraints.Ordered
}

type RealNumber interface {
	int | float64
}

func performBooleanComparisionOperation(left, right bool, op types.TOKEN_TYPE) bool {
	switch op {
	case types.NOT:
		return !right
	case types.EQUALS:
		return left == right
	case types.NOT_EQUAL:
		return left != right
	default:
		panic(fmt.Sprintf("unsupported operation: %v", op))
	}
}

func performComparisionOperation[T Comparator](left, right T, op types.TOKEN_TYPE) bool {
	switch op {
	case types.GREATER:
		return left > right
	case types.GREATER_OR_EQUAL:
		return left >= right
	case types.LESSER:
		return left < right
	case types.LESSER_OR_EQUAL:
		return left <= right
	case types.EQUALS:
		return left == right
	case types.NOT_EQUAL:
		return left != right
	case types.TRUE:
		return true
	case types.FALSE:
		return false
	default:
		panic(fmt.Sprintf("unsupported operation: %v", op))
	}
}

func performOperation[T RealNumber](left, right T, op types.TOKEN_TYPE) T {
	switch op {
	case types.PLUS:
		return left + right
	case types.MINUS:
		return left - right
	case types.MULTIPLY:
		return left * right
	case types.DIVIDE:
		return left / right
	default:
		panic(fmt.Sprintf("unsupported operation: %v", op))
	}
}

type EvalResult struct {
	Value interface{}
	Type  string
}

func (e *Evaluator) Evaluate(s types.Node) interface{} {
	return e.Visit(s)
}

func (e *Evaluator) Visit(node types.Node) interface{} {
	switch n := node.(type) {
	case types.Statement:
		return e.Visit(node.(types.Statement).Node)
	case types.Expression:
		return e.Visit(node.(types.Expression).Node)

	case types.UnaryOP:
		op := n.Left
		right := e.Visit(n.Right)

		switch right.(type) {
		case float64:
			return performOperation(0, right.(float64), op)
		case int:
			return performOperation(0, right.(int), op)
		case bool:
			return performBooleanComparisionOperation(false, right.(bool), op)
		}
	case types.BinOP:
		left := e.Visit(n.Left)
		right := e.Visit(n.Right)

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
				switch n.Op {
				case types.PLUS, types.MINUS, types.MULTIPLY, types.DIVIDE:
					return performOperation(left, right, n.Op)
				default:
					return performComparisionOperation(left, right, n.Op)
				}
			}
		case float64:
			if right, ok := right.(float64); ok {
				switch n.Op {
				case types.PLUS, types.MINUS, types.MULTIPLY, types.DIVIDE:
					return performOperation(left, right, n.Op)
				default:
					return performComparisionOperation(left, right, n.Op)
				}
			}
		case bool:
			return performBooleanComparisionOperation(left, right.(bool), n.Op)
		}

	case types.IfElseStatement:
		condition := e.Visit(n.Condition)
		if condition.(bool) {
			for _, statement := range n.IfChildren {
				e.Visit(statement)
			}
		} else {
			for _, statement := range n.ElseChildren {
				e.Visit(statement)
			}
		}
	case types.Assign:
		right := e.Visit(n.Value)
		e.SymbolTable[n.Id.(types.Id).Name] = right
		return right
	case types.Id:
		value, ok := e.SymbolTable[n.Name]
		if !ok {
			log.Fatalf("variable '%s' not defined", n)
		}
		return value
	case types.Literal:
		nodeValue := n.Value.(types.Literal).Value
		switch v := nodeValue.(type) {
		case int, float64, string, bool:
			return v
		default:
			log.Fatalf("unsupported type %s\n", v)

		}
		return node
	}
	return node
}
