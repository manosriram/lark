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
	content, err := os.ReadFile("source.lark")
	if err != nil {
		log.Fatal("error reading source file")
	}

	symbolTable = make(map[string]interface{})
	tokens := token.Tokenize(string(content))
	builder := ast.NewAstBuilder(tokens.Tokens)
	var tree types.Node
	// var statements []types.Statement
	var nodes []types.Node
	for builder.CurrentTokenPointer < len(tokens.Tokens)-1 {
		tree = builder.Parse()
		// treeType := fmt.Sprintf("%T", tree)
		if tree != nil {
			// statement := types.Statement{
			// Node: tree,
			// }
			nodes = append(nodes, tree)
			// switch treeType := tree.(type) {
			// case types.BinOP:
			// break
			// }
			// statements = append(statements, statement)
		}
	}
	fmt.Println(nodes)
	for _, node := range nodes {
		result := ast.Evaluate(node)
		fmt.Println(result)
		switch nType := node.(type) {
		case types.Statement:
			switch node.(types.Statement).StatementType {
			case types.AssignType:
				fmt.Println("n = ", node.(types.Statement).Node)
				// assign := node.(types.Assign).Id.String()
				// symbolTable[assign] = result
				break

			}
		default:
			fmt.Println(nType)
		}
	}

	for k, v := range symbolTable {
		fmt.Printf("var_name: %s, var_value: %v\n", k, v)
	}
}
