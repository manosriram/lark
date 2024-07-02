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
		treeType := fmt.Sprintf("%T", tree)
		// fmt.Println(treeType)
		// fmt.Println(tree)
		if tree != nil {
			statement := ast.Statement{
				Node: tree,
			}
			switch treeType {
			case "ast.Assign":
				statement.StatementType = token.ASSIGN_STATEMENT
				break
			case "ast.BinOP":
				statement.StatementType = token.EXPRESSION_STATEMENT
				break
			}

			statements = append(statements, statement)

		}
	}
	for _, statement := range statements {
		// result := ast.Evaluate(statement)
		if statement.StatementType == token.ASSIGN_STATEMENT {
			assign := statement.Node.(ast.Assign)
			fmt.Println(assign.Id.(ast.Id).Name, assign.Value.(ast.Literal).Value)
		}
		// fmt.Println(result)
	}
}
