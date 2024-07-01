package token

import (
	"log"
	"strconv"
	"unicode"
)

type TOKEN_TYPE string

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
	NUMBER               = "number"
	EOF                  = "EOF"
)

type Token struct {
	Position  int
	TokenType TOKEN_TYPE
	Value     interface{}
}

type Source struct {
	Content         string
	Tokens          []Token
	CurrentPosition int
	TokenSoFar      string
}

func isnumeric(c rune) bool {
	return unicode.IsNumber(c)
}

func isalpha(c rune) bool {
	return unicode.IsLetter(c)
}

func (s *Source) eatNum() {
	for s.CurrentPosition < len(s.Content) && (isnumeric(rune(s.Content[s.CurrentPosition]))) {
		s.CurrentPosition += 1
	}
	// s.CurrentPosition += 1
}

func (s *Source) eatVar() {
	for s.CurrentPosition < len(s.Content) && (isalpha(rune(s.Content[s.CurrentPosition]))) {
		s.CurrentPosition += 1
	}
	// s.CurrentPosition += 1
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
	// s.eat()
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
		// case '!':
		// if string(s.Content[s.CurrentPosition+1]) == "=" {
		// s.Tokens = append(s.Tokens, NOT_EQUAL)
		// s.eat()
		// } else {
		// s.Tokens = append(s.Tokens, NOT)
		// }
		// s.eat()
		// break
		// case '(':
		// s.Tokens = append(s.Tokens, LBRACE)
		// s.eat()
		// break
		// case ')':
		// s.Tokens = append(s.Tokens, RBRACE)
		// s.eat()
		// break
		case ';':
			s.Tokens = append(s.Tokens, Token{TokenType: SEMICOLON, Value: ';'})
			s.eat()
		case ' ':
			s.eat()
		case '\n':
			s.eat()
		default: // variable decl
			if unicode.IsNumber(rune(s.Content[s.CurrentPosition])) {
				before := s.CurrentPosition
				s.eatNum()
				after := s.CurrentPosition
				variable := s.Content[before:after]
				n, e := strconv.Atoi(variable)
				if e != nil {
					log.Fatal("cannot parse source")
				}
				s.Tokens = append(s.Tokens, Token{TokenType: NUMBER, Value: n})
			} else {
				before := s.CurrentPosition
				s.eatVar()
				after := s.CurrentPosition
				variable := s.Content[before:after]
				s.Tokens = append(s.Tokens, Token{TokenType: ID, Value: variable})
			}

		}
	}

	return &s
}
