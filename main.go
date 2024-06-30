package main

import (
	"fmt"
	"lark/pkg/ast"
	token "lark/pkg/token"
	"log"
	"os"
)

func main() {
	content, err := os.ReadFile("./source.lark")
	if err != nil {
		log.Fatal("error reading source file")
	}

	tokens := token.Tokenize(string(content))
	builder := ast.NewAstBuilder(tokens.Tokens)
	tree := builder.Parse()

	fmt.Println(tree)
}
