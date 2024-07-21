package ast

import (
	"fmt"
	types "lark/pkg/types"
	"log"

	"golang.org/x/exp/constraints"
)

type Evaluator struct {
	SymbolTable      map[string]interface{}
	LocalSymbolTable map[string]interface{}
}

type Comparator interface {
	constraints.Ordered
}

type RealNumber interface {
	int | float64 | int32 | float32
}

func performBooleanComparisionOperation(left, right bool, op types.TOKEN_TYPE) bool {
	switch op {
	case types.NOT:
		return !right
	case types.EQUALS:
		return left == right
	case types.NOT_EQUAL:
		return left != right
	case types.OR:
		return left || right
	case types.AND:
		return left && right
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
		if n.Type == types.GLOBAL_ASSIGN {
			e.SymbolTable[n.Id.(types.Id).Name] = right
		} else {
			e.LocalSymbolTable[n.Id.(types.Id).Name] = right
		}
		return right
	case types.Function:
		e.SymbolTable[n.Name] = n
	case types.FunctionCall:
		_, ok := e.SymbolTable[n.Name]
		if !ok {
			log.Fatalf("variable '%s' not defined", n.Name)
		}
		if len(node.(types.FunctionCall).Arguments) != len(e.SymbolTable[n.Name].(types.Function).Arguments) {
			log.Fatalf("function arguments mismatch\n")
		}

		functionCallArgs := node.(types.FunctionCall).Arguments
		for i, v := range e.SymbolTable[n.Name].(types.Function).Arguments {
			switch functionCallArgs[i].(type) {
			case types.Literal:
				switch functionCallArgs[i].(types.Literal).Type {
				case types.IDENT:
					value := e.SymbolTable[functionCallArgs[i].(types.Literal).Value.(string)]
					e.LocalSymbolTable[v.(types.Literal).Value.(string)] = value
				default:
					e.LocalSymbolTable[v.(types.Literal).Value.(string)] = e.Visit(functionCallArgs[i])
				}
			default:
				value := e.Visit(functionCallArgs[i])
				variable := e.Visit(v)
				e.LocalSymbolTable[variable.(string)] = value
			}
		}
		defer func() {
			for _, v := range e.SymbolTable[n.Name].(types.Function).Arguments {
				delete(e.LocalSymbolTable, v.(types.Literal).Value.(string))
			}
		}()

		for _, v := range e.SymbolTable[n.Name].(types.Function).Children {
			e.Visit(v)
		}
		return e.Visit(e.SymbolTable[n.Name].(types.Function).ReturnExpression)
	case types.Array:
		val := e.SymbolTable[n.Name].([]interface{})
		if n.Index >= len(val) || n.Index < 0 {
			log.Fatalf("array index (%d) out of bounds (%d)\n", n.Index, len(val))
		}
		return val[n.Index]
	case types.Swap:
		_, ok := e.SymbolTable[n.Left.String()]
		if !ok {
			log.Fatalf("variable '%s' not defined", n.Left.String())
		}
		_, ok = e.SymbolTable[n.Right.String()]
		if !ok {
			log.Fatalf("variable '%s' not defined", n.Right.String())
		}
		e.SymbolTable[n.Left.String()], e.SymbolTable[n.Right.String()] = e.SymbolTable[n.Right.String()], e.SymbolTable[n.Left.String()]
	case types.Id:
		switch n.Type {
		case types.ARRAY_INDEX_POSITION:
			indexPosition := n.Value.(types.Literal).Value.(int)
			arr, ok := e.LocalSymbolTable[n.Name].([]interface{})
			if !ok {
				arr, ok = e.SymbolTable[n.Name].([]interface{})
				if !ok {
					log.Fatalf("variable '%s' not defined", n)
				}
				if indexPosition >= len(arr) || indexPosition < 0 {
					log.Fatalf("array index (%d) out of bounds (%d)\n", indexPosition, len(arr))
				}
				return arr[indexPosition]
			}
			return arr[indexPosition]
		default:
			value, ok := e.LocalSymbolTable[n.Name]
			if !ok {
				value, ok = e.SymbolTable[n.Name]
				if !ok {
					log.Fatalf("variable '%s' not defined", n)
				}
				return value
			}
			return value
		}
	case types.Literal:
		switch v := n.Value.(type) {
		case int, float64, string, bool:
			return v
		case []interface{}:
			return v
		}
		nodeValue := n.Value
		switch nodeValue.(type) {
		case types.Literal:
			switch nodeValue.(types.Literal).Type {
			case types.INTEGER, types.FLOAT64, types.STRING, types.BOOLEAN:
				return nodeValue.(types.Literal).Value
			default:
				log.Fatalf("unsupported type %s\n", "")
			}
		default:
			break
		}
		return node
	}
	return node
}
