package main

import (
	"fmt"
	"lark/pkg/ast"
	token "lark/pkg/token"
	"log"
	"os"
)

func main() {
	content, err := os.ReadFile("source.lark")
	if err != nil {
		log.Fatal("error reading source file")
	}

	tokens := token.Tokenize(string(content))

	builder := ast.NewAstBuilder(tokens.Tokens)
	var tree interface{}
	var statements []ast.Statement
	for builder.CurrentTokenPointer < len(tokens.Tokens)-1 {
		tree = builder.Parse()
		if tree != nil {
			statements = append(statements, ast.Statement{
				Node: tree,
			})
		}
	}
	for _, statement := range statements {
		fmt.Println(statement)
	}
}
