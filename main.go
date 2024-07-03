package main

import (
	"fmt"
	"lark/pkg/ast"
	token "lark/pkg/token"
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
	var tree interface{}
	var statements []ast.Statement
	for builder.CurrentTokenPointer < len(tokens.Tokens)-1 {
		tree = builder.Parse()
		treeType := fmt.Sprintf("%T", tree)
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
		result := ast.Evaluate(statement)
		switch statement.StatementType {
		case token.ASSIGN_STATEMENT:
			assign := statement.Node.(ast.Assign)
			id := assign.Id.(ast.Id).Name
			symbolTable[id] = result
			break
		}
	}

	for k, v := range symbolTable {
		fmt.Printf("var_name: %s, var_value: %v\n", k, v)
	}
}
