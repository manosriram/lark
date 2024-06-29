package main

import (
	"fmt"
	token "lark/pkg/token"
	"log"
	"os"
)

func main() {
	content, err := os.ReadFile("./source.lrk")
	if err != nil {
		log.Fatal("error reading source file")
	}

	tokens := token.Tokenize(string(content))
	for _, x := range tokens.Tokens {
		fmt.Println(x)
	}

}
