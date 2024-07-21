package types

type TOKEN_TYPE string
type STATEMENT_TYPE int
type EXPRESSION_TYPE int
type LITERAL_TYPE int
type ASSIGN_TYPE int

const (
	EXPRESSION_STATEMENT STATEMENT_TYPE = iota
	ASSIGN_STATEMENT
	IF_STATEMENT
)

const (
	GLOBAL_ASSIGN ASSIGN_TYPE = iota
	LOCAL_ASSIGN
)

const (
	UNARY_OP EXPRESSION_TYPE = iota
	BIN_OP
)

const (
	FLOAT64 LITERAL_TYPE = iota
	INTEGER
	STRING
	BOOLEAN
	ARRAY
	OPERATOR
	STATEMENT
	EXPRESSION
	KEYWORD
	IDENT
	ARRAY_INDEX_POSITION
)

const (
	PLUS     TOKEN_TYPE = "+"
	MINUS               = "-"
	MULTIPLY            = "*"
	DIVIDE              = "/"
	LBRACE              = "("
	RBRACE              = ")"
	LPAREN              = "{"
	RPAREN              = "}"
	ASSIGN              = "<-"

	GREATER          = ">"
	GREATER_OR_EQUAL = ">="
	LESSER           = "<"
	LESSER_OR_EQUAL  = "<="
	NOT              = "!"
	NOT_EQUAL        = "!="
	EQUALS           = "=="

	IF   = "if"
	ELSE = "else"

	SEMICOLON = ";"
	ID        = "id"
	LITERAL   = "literal"
	EOF       = "EOF"
	EXPR      = "expr"
	LOCAL     = "local"
	DOT       = "."
	COMMENT   = "//"
	SWAP      = "<->"

	TRUE  = "true"
	FALSE = "false"

	OR          = "||"
	AND         = "&&"
	BITWISE_OR  = "|"
	BITWISE_AND = "&"

	FUNCTION                    = "fn"
	FUNCTION_ARGUMENT_OPEN      = "["
	FUNCTION_ARGUMENT_CLOSE     = "]"
	FUNCTION_RETURN             = "return"
	FUNCTION_ARGUMENT_SEPARATOR = ","
	STATEMENT_BLOCK_OPEN        = "<<"
	STATEMENT_BLOCK_CLOSE       = ">>"
	FUNCTION_CALL_OPEN          = "("
	FUNCTION_CALL_CLOSE         = ")"
	FUNCTION_CALL               = "()"
	FUNCTION_ARGUMENT           = "arg"

	ARRAY_OPEN      = "["
	ARRAY_CLOSE     = "]"
	ARRAY_SEPARATOR = ","
	ARRAY_INDEX     = "@"
)

type Token struct {
	Position    int
	TokenType   TOKEN_TYPE
	Value       Node
	LiteralType LITERAL_TYPE
	LineNumber  int
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
	return "STATEMENT"
}

func (e Expression) NodeType() string {
	return "EXPRESSION"
}

type FunctionCall struct {
	Name      string
	Arguments []Node
}

func (f FunctionCall) NodeType() string {
	return FUNCTION_CALL
}

func (f FunctionCall) String() string {
	return f.Name
}

type Function struct {
	Name             string
	Arguments        []Node
	Children         []Node
	ReturnExpression Node
	Variables        []string
}

func (f Function) NodeType() string {
	return FUNCTION
}

func (f Function) String() string {
	return f.Name
}

type IfElseStatement struct {
	Condition    Node
	IfChildren   []Node
	ElseChildren []Node
}

func (i IfElseStatement) NodeType() string {
	return "IF_ELSE"
}

func (i IfElseStatement) String() string {
	return ""
}

type UnaryOP struct {
	Left  TOKEN_TYPE
	Right Node
}

func (u UnaryOP) NodeType() string {
	return "UNARY_OP"
}

func (u UnaryOP) String() string {
	return ""
}

type BinOP struct {
	Left  Node
	Right Node
	Op    TOKEN_TYPE
}

func (b BinOP) NodeType() string {
	return "BIN_OP"
}

func (b BinOP) String() string {
	return ""
}

type Assign struct {
	Id    Node
	Value Node
	Type  ASSIGN_TYPE
}

func (a Assign) NodeType() string {
	return ASSIGN
}

func (a Assign) String() string {
	return a.Id.String()
}

type Swap struct {
	Left  Node
	Right Node
}

func (s Swap) NodeType() string {
	return SWAP
}

func (s Swap) String() string {
	return ""
}

type Array struct {
	Name  string
	Value interface{}
	Index int
}

func (a Array) NodeType() string {
	return "ARRAY"
}

func (a Array) String() string {
	return a.Name
}

type Literal struct {
	Value interface{}
	Type  LITERAL_TYPE
}

func (l Literal) NodeType() string {
	return LITERAL
}

func (l Literal) String() string {
	return ""
}

type Id struct {
	Name  string
	Value Node
	Type  LITERAL_TYPE
}

func (i Id) NodeType() string {
	return ID
}

func (i Id) String() string {
	return i.Name
}
