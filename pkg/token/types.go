package token

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
)

const (
	PLUS      TOKEN_TYPE = "+"
	MINUS                = "-"
	MULTIPLY             = "*"
	DIVIDE               = "/"
	LBRACE               = "("
	RBRACE               = ")"
	EQUAL                = "="
	NOT                  = "!"
	NOT_EQUAL            = "!="
	SEMICOLON            = ";"
	ID                   = "id"
	LITERAL              = "literal"
	EOF                  = "EOF"
	EXPR                 = "expr"
	DOT                  = "."
)

type Token struct {
	Position    int
	TokenType   TOKEN_TYPE
	Value       interface{}
	LiteralType LITERAL_TYPE
}

type Source struct {
	Content         string
	Tokens          []Token
	CurrentPosition int
	TokenSoFar      string
}
