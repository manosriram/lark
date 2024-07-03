package token

import (
	"log"
	"strconv"
	"unicode"
)

type TOKEN_TYPE string
type STATEMENT_TYPE int
type LITERAL_TYPE int

const (
	EXPRESSION_STATEMENT STATEMENT_TYPE = iota
	ASSIGN_STATEMENT
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

func (s *Source) getCurrentToken() byte {
	return s.Content[s.CurrentPosition]
}

func isnumeric(c rune) bool {
	return unicode.IsDigit(c)
}

func isalpha(c rune) bool {
	return unicode.IsLetter(c)
}

func (s *Source) eatNum() {
	for s.CurrentPosition < len(s.Content) && (isnumeric(rune(s.Content[s.CurrentPosition])) || string(s.getCurrentToken()) == DOT) {
		s.CurrentPosition += 1
	}
}

func (s *Source) eatVar() {
	for s.CurrentPosition < len(s.Content) && (isalpha(rune(s.Content[s.CurrentPosition]))) {
		s.CurrentPosition += 1
	}
}
func (s *Source) eat() {
	done := false
	for s.CurrentPosition < len(s.Content) && (s.Content[s.CurrentPosition] == ' ' || s.Content[s.CurrentPosition] == '\n') {
		s.CurrentPosition += 1
		done = true
	}
	if !done {
		s.CurrentPosition += 1
	}
}

func Tokenize(source string) *Source {
	s := Source{
		Content:         source,
		CurrentPosition: 0,
		Tokens:          make([]Token, 0),
		TokenSoFar:      "",
	}
	for s.CurrentPosition < len(s.Content) {
		charAtPosition := s.Content[s.CurrentPosition]
		switch charAtPosition {
		case '+':
			s.Tokens = append(s.Tokens, Token{TokenType: PLUS, Value: '+'})
			s.eat()
		case '-':
			s.Tokens = append(s.Tokens, Token{TokenType: MINUS, Value: '-'})
			s.eat()
		case '*':
			s.Tokens = append(s.Tokens, Token{TokenType: MULTIPLY, Value: '*'})
			s.eat()
		case '/':
			s.Tokens = append(s.Tokens, Token{TokenType: DIVIDE, Value: '/'})
			s.eat()
		case '=':
			s.Tokens = append(s.Tokens, Token{TokenType: EQUAL, Value: '='})
			s.eat()
		case '(':
			s.Tokens = append(s.Tokens, Token{TokenType: LBRACE, Value: '('})
			s.eat()
		case ')':
			s.Tokens = append(s.Tokens, Token{TokenType: RBRACE, Value: ')'})
			s.eat()
		case '"':
			s.eat()
			before := s.CurrentPosition
			for s.Content[s.CurrentPosition] != '"' {
				s.CurrentPosition += 1
			}
			s.eat()
			variable := s.Content[before : s.CurrentPosition-1]
			s.Tokens = append(s.Tokens, Token{TokenType: LITERAL, Value: variable, LiteralType: STRING})
			break
		case ';':
			s.Tokens = append(s.Tokens, Token{TokenType: SEMICOLON, Value: ';'})
			s.eat()
			break
		case ' ':
			s.eat()
			break
		case '\n':
			s.eat()
			break
		default:
			if unicode.IsNumber(rune(s.Content[s.CurrentPosition])) {
				before := s.CurrentPosition
				s.eatNum()
				after := s.CurrentPosition
				variable := s.Content[before:after]
				if number, err := strconv.Atoi(variable); err == nil {
					s.Tokens = append(s.Tokens, Token{TokenType: LITERAL, Value: number, LiteralType: INTEGER})
				} else if number, err := strconv.ParseFloat(variable, 64); err == nil {
					s.Tokens = append(s.Tokens, Token{TokenType: LITERAL, Value: number, LiteralType: FLOAT64})
				} else {
					log.Fatalf("cannot parse source file\n")
				}
			} else if unicode.IsLetter(rune(s.Content[s.CurrentPosition])) {
				before := s.CurrentPosition
				s.eatVar()
				after := s.CurrentPosition
				variable := s.Content[before:after]
				s.Tokens = append(s.Tokens, Token{TokenType: ID, Value: variable, LiteralType: STRING})
			} else {
				log.Fatalf("unsupported type %v\n", string(s.Content[s.CurrentPosition]))
			}

		}
	}

	return &s
}
