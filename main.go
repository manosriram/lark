package main

import (
	"fmt"
	"lark/pkg/ast"
	token "lark/pkg/token"
	"lark/pkg/types"
	"log"
	"os"
)

var symbolTable map[string]interface{}

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("not enough arguments\n")
	}
	file := os.Args[1]
	content, err := os.ReadFile(file)
	if err != nil {
		log.Fatal("error reading source file")
	}

	root := types.Compound{Children: []types.Node{}}

	symbolTable = make(map[string]interface{})
	tokens := token.Tokenize(string(content))
	builder := ast.NewAstBuilder(tokens.Tokens)
	var tree types.Node
	for builder.CurrentTokenPointer < len(tokens.Tokens)-1 {
		tree = builder.Parse()
		if tree != nil {
			root.Children = append(root.Children, tree)
		}
	}
	evaluator := ast.Evaluator{
		SymbolTable: symbolTable,
	}

	for _, node := range root.Children {
		evaluator.Evaluate(node)
	}

	for k, v := range evaluator.SymbolTable {
		fmt.Printf("var_name: %s, var_value: %v\n", k, v)
	}
}
