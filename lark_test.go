package main

import (
	token "lark/pkg/token"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Tokenize(t *testing.T) {
	content, err := os.ReadFile("test_source_files/token.lark")
	assert.Equal(t, nil, err)

	tokens := token.Tokenize(string(content))
	assert.NotNil(t, tokens)
	assert.Equal(t, 8, len(tokens.Tokens))
}
