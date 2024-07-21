package main

import (
	"lark/pkg/ast"
	token "lark/pkg/token"
	"lark/pkg/types"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func evaluate(t *testing.T, sourceFile string) map[string]interface{} {
	content, err := os.ReadFile(sourceFile)
	assert.Equal(t, nil, err)

	root := types.Compound{Children: []types.Node{}}
	symbolTable = make(map[string]interface{})
	localSymbolTable = make(map[string]interface{})
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
		SymbolTable:      symbolTable,
		LocalSymbolTable: localSymbolTable,
	}

	for _, node := range root.Children {
		evaluator.Evaluate(node)
	}

	return symbolTable
}

func Test_Tokenize(t *testing.T) {
	content, err := os.ReadFile("test_source_files/token.lark")
	assert.Equal(t, nil, err)

	expectedOrderOfTokens := []types.Token{
		{TokenType: types.ID},
		{TokenType: types.ASSIGN},
		{TokenType: types.LITERAL},
		{TokenType: types.PLUS},
		{TokenType: types.LITERAL},
		{TokenType: types.SEMICOLON},
		{TokenType: types.ID},
		{TokenType: types.ASSIGN},
		{TokenType: types.LITERAL},
		{TokenType: types.MINUS},
		{TokenType: types.LITERAL},
		{TokenType: types.SEMICOLON},
	}

	tokens := token.Tokenize(string(content))
	assert.Equal(t, 12, len(tokens.Tokens))
	assert.NotNil(t, tokens)
	for i, token := range tokens.Tokens {
		assert.Equal(t, expectedOrderOfTokens[i].TokenType, token.TokenType)
	}
}

func areMapsSame(map1, map2 map[string]interface{}) bool {
	if len(map1) != len(map2) {
		return false
	}
	for key, value1 := range map1 {
		value2, exists := map2[key]
		if !exists {
			return false
		}
		if !areValuesEqual(value1, value2) {
			return false
		}
	}
	return true
}

func areValuesEqual(v1, v2 interface{}) bool {
	switch v1 := v1.(type) {
	case nil:
		return v2 == nil
	case int:
		v2, ok := v2.(int)
		return ok && v1 == v2
	case string:
		v2, ok := v2.(string)
		return ok && v1 == v2
	case bool:
		v2, ok := v2.(bool)
		return ok && v1 == v2
	case float64:
		v2, ok := v2.(float64)
		return ok && v1 == v2
	case []interface{}:
		v2, ok := v2.([]interface{})
		if !ok || len(v1) != len(v2) {
			return false
		}
		for i := range v1 {
			if !areValuesEqual(v1[i], v2[i]) {
				return false
			}
		}
		return true
	case map[string]interface{}:
		v2, ok := v2.(map[string]interface{})
		if !ok {
			return false
		}
		return areMapsSame(v1, v2)
	default:
		return false
	}
}
func Test_Parser(t *testing.T) {
	symbolTable := evaluate(t, "test_source_files/parse.lark")
	expectedSymbolTableVars := map[string]interface{}{
		"a":   false,
		"b":   true,
		"c":   -123,
		"e":   -117,
		"tt":  true,
		"d":   "not_ok",
		"and": false,
		"or":  true,
	}

	assert.Equal(t, len(expectedSymbolTableVars), len(symbolTable))
	for key := range expectedSymbolTableVars {
		assert.Equal(t, expectedSymbolTableVars[key], symbolTable[key])
	}
}

func Test_Function(t *testing.T) {
	symbolTable := evaluate(t, "test_source_files/function.lark")
	assert.Equal(t, 12, len(symbolTable))
	assert.Equal(t, 1000, symbolTable["fna"])
	assert.Equal(t, 500, symbolTable["fnb"])
	assert.Equal(t, 1500, symbolTable["fnval"])
	assert.Equal(t, 3, symbolTable["sum"])
	assert.Equal(t, 9, symbolTable["localSum"])
	assert.Equal(t, 103, symbolTable["dynamicSum"])
	assert.Equal(t, 603, symbolTable["expressionSum"])
}

func Test_Array(t *testing.T) {
	symbolTable := evaluate(t, "test_source_files/array.lark")
	assert.Equal(t, 6, len(symbolTable))
	assert.Equal(t, 60, symbolTable["all"])
	assert.Equal(t, 300, symbolTable["sum"])
	assert.Equal(t, 7, symbolTable["twoSum"])
}
