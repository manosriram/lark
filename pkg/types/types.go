package types

import "fmt"

type TOKEN_TYPE string
type STATEMENT_TYPE int
type EXPRESSION_TYPE int
type LITERAL_TYPE int

const (
	EXPRESSION_STATEMENT STATEMENT_TYPE = iota
	ASSIGN_STATEMENT
)

const (
	UNARY_OP EXPRESSION_TYPE = iota
	BIN_OP
)

const (
	FLOAT64 LITERAL_TYPE = iota
	INTEGER
	STRING
	OPERATOR
)

const (
	PLUS     TOKEN_TYPE = "+"
	MINUS               = "-"
	MULTIPLY            = "*"
	DIVIDE              = "/"
	LBRACE              = "("
	RBRACE              = ")"
	ASSIGN              = "="

	GREATER          = ">"
	GREATER_OR_EQUAL = ">="
	LESSER           = "<"
	LESSER_OR_EQUAL  = "<="
	NOT              = "!"
	NOT_EQUAL        = "!="
	EQUALS           = "=="

	SEMICOLON = ";"
	ID        = "id"
	LITERAL   = "literal"
	EOF       = "EOF"
	EXPR      = "expr"
	DOT       = "."
)

type Token struct {
	Position    int
	TokenType   TOKEN_TYPE
	Value       Node
	LiteralType LITERAL_TYPE
}

type StatementType int

const (
	_ StatementType = iota
	AssignType
)

type ExpressionType int

const (
	_ ExpressionType = iota
	UnaryOpType
	BinOpType
	IdType
	LiteralType
)

type Node interface {
	NodeType() string
	String() string
}

type Compound struct {
	Children []Node
}

func (s Compound) NodeType() string {
	return "compound"
}

func (e Compound) String() string {
	return "compound"
}

type Statement struct {
	Node
	StatementType
}

type Expression struct {
	Node
	ExpressionType
}

func (s Statement) NodeType() string {
	return "statement"
}

func (e Expression) NodeType() string {
	return "expression"
}

type UnaryOP struct {
	Left  TOKEN_TYPE
	Right Node
}

func (u UnaryOP) NodeType() string {
	return "unary"
}

func (u UnaryOP) String() string {
	return string(u.Left)
}

type BinOP struct {
	Left  Node
	Right Node
	Op    TOKEN_TYPE
}

func (b BinOP) NodeType() string {
	return "bin"
}

func (b BinOP) String() string {
	return string(b.Op)
}

type Assign struct {
	Id    Node
	Value Node
}

func (a Assign) String() string {
	return a.Id.String()
}

func (a Assign) NodeType() string {
	return "assign"
}

type Literal struct {
	Value interface{}
	Type  LITERAL_TYPE
}

func (l Literal) NodeType() string {
	return "literal"
}

func (l Literal) String() string {
	return fmt.Sprintf("%v", l.Type)
}

type Id struct {
	Name  string
	Value Node
}

func (i Id) NodeType() string {
	return "id"
}

func (i Id) String() string {
	return i.Name
}