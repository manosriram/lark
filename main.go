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
	tr := ast.Expr(tokens.Tokens)
	fmt.Println(tr)
	// for _, x := range tokens.Tokens {
	// fmt.Println(x)
	// }

}
