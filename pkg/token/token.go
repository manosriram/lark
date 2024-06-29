package token

import "unicode"

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
	VAR                  = "var"
	EOF                  = "EOF"
)

type Token struct {
	position  int
	tokenType TOKEN_TYPE
}

type Source struct {
	Content         string
	Tokens          []TOKEN_TYPE
	CurrentPosition int
	TokenSoFar      string
}

func isalphanumeric(c rune) bool {
	if !unicode.IsLetter(c) && !unicode.IsNumber(c) {
		return false
	}
	return true
}
func (s *Source) eatVar() {
	for s.CurrentPosition < len(s.Content) && (isalphanumeric(rune(s.Content[s.CurrentPosition]))) {
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
		Tokens:          make([]TOKEN_TYPE, 0),
		TokenSoFar:      "",
	}
	// s.eat()
	for s.CurrentPosition < len(s.Content) {
		charAtPosition := s.Content[s.CurrentPosition]
		switch charAtPosition {
		case '+':
			s.Tokens = append(s.Tokens, PLUS)
			s.eat()
			break
		case '-':
			s.Tokens = append(s.Tokens, MINUS)
			s.eat()
			break
		case '*':
			s.Tokens = append(s.Tokens, MULTIPLY)
			s.eat()
			break
		case '/':
			s.Tokens = append(s.Tokens, DIVIDE)
			s.eat()
			break
		case '=':
			s.Tokens = append(s.Tokens, EQUAL)
			s.eat()
			break
		case '!':
			if string(s.Content[s.CurrentPosition+1]) == "=" {
				s.Tokens = append(s.Tokens, NOT_EQUAL)
				s.eat()
			} else {
				s.Tokens = append(s.Tokens, NOT)
			}
			s.eat()
			break
		case '(':
			s.Tokens = append(s.Tokens, LBRACE)
			s.eat()
			break
		case ')':
			s.Tokens = append(s.Tokens, RBRACE)
			s.eat()
			break
		case ';':
			s.Tokens = append(s.Tokens, SEMICOLON)
			s.eat()
			break
		case ' ':
			s.eat()
			break
		case '\n':
			s.eat()
			break
		default: // variable decl
			// bef := s.CurrentPosition
			s.eatVar()
			// variable := s.Content[bef:aft]
			// s.eat()
			s.Tokens = append(s.Tokens, VAR)

		}
	}

	return &s
}
