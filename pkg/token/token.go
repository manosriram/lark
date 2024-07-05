package token

import (
	"lark/pkg/types"
	"log"
	"strconv"
	"unicode"
)

type Source struct {
	Content         string
	Tokens          []types.Token
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
	for s.CurrentPosition < len(s.Content) && (isnumeric(rune(s.Content[s.CurrentPosition])) || string(s.getCurrentToken()) == types.DOT) {
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
		Tokens:          make([]types.Token, 0),
		TokenSoFar:      "",
	}
	for s.CurrentPosition < len(s.Content) {
		charAtPosition := s.Content[s.CurrentPosition]
		switch charAtPosition {
		case '+':
			s.Tokens = append(s.Tokens, types.Token{TokenType: types.PLUS, Value: types.Literal{Value: '+', Type: types.OPERATOR}})
			s.eat()
		case '-':
			s.Tokens = append(s.Tokens, types.Token{TokenType: types.MINUS, Value: types.Literal{Value: '-', Type: types.OPERATOR}})
			s.eat()
		case '*':
			s.Tokens = append(s.Tokens, types.Token{TokenType: types.MULTIPLY, Value: types.Literal{Value: '*', Type: types.OPERATOR}})
			s.eat()
		case '/':
			s.Tokens = append(s.Tokens, types.Token{TokenType: types.DIVIDE, Value: types.Literal{Value: '/', Type: types.OPERATOR}})
			s.eat()
		case '=':
			s.eat()
			switch s.Content[s.CurrentPosition] {
			case '=':
				s.Tokens = append(s.Tokens, types.Token{TokenType: types.EQUALS, Value: types.Literal{Value: "==", Type: types.OPERATOR}})
				s.eat()
			default:
				s.Tokens = append(s.Tokens, types.Token{TokenType: types.ASSIGN, Value: types.Literal{Value: '=', Type: types.OPERATOR}})
			}
		case '(':
			s.Tokens = append(s.Tokens, types.Token{TokenType: types.LBRACE, Value: types.Literal{Value: '(', Type: types.OPERATOR}})
			s.eat()
		case ')':
			s.Tokens = append(s.Tokens, types.Token{TokenType: types.RBRACE, Value: types.Literal{Value: ')', Type: types.OPERATOR}})
			s.eat()
		case '!':
			s.eat()
			switch s.Content[s.CurrentPosition] {
			case '=':
				s.Tokens = append(s.Tokens, types.Token{TokenType: types.NOT_EQUAL, Value: types.Literal{Value: "!=", Type: types.OPERATOR}})
				s.eat()
			default:
				s.Tokens = append(s.Tokens, types.Token{TokenType: types.NOT, Value: types.Literal{Value: "!", Type: types.OPERATOR}})
			}
		case '<':
			s.eat()
			switch s.Content[s.CurrentPosition] {
			case '=':
				s.Tokens = append(s.Tokens, types.Token{TokenType: types.LESSER_OR_EQUAL, Value: types.Literal{Value: "<=", Type: types.OPERATOR}})
				s.eat()
			default:
				s.Tokens = append(s.Tokens, types.Token{TokenType: types.LESSER, Value: types.Literal{Value: "<", Type: types.OPERATOR}})
			}
		case '>':
			s.eat()
			switch s.Content[s.CurrentPosition] {
			case '=':
				s.Tokens = append(s.Tokens, types.Token{TokenType: types.GREATER_OR_EQUAL, Value: types.Literal{Value: ">=", Type: types.OPERATOR}})
				s.eat()
			default:
				s.Tokens = append(s.Tokens, types.Token{TokenType: types.GREATER, Value: types.Literal{Value: ">", Type: types.OPERATOR}})
			}
		case '"':
			s.eat()
			before := s.CurrentPosition
			for s.Content[s.CurrentPosition] != '"' {
				s.CurrentPosition += 1
			}
			s.eat()
			variable := s.Content[before : s.CurrentPosition-1]
			s.Tokens = append(s.Tokens, types.Token{TokenType: types.LITERAL, Value: types.Literal{Value: variable, Type: types.STRING}})
			break
		case ';':
			s.Tokens = append(s.Tokens, types.Token{TokenType: types.SEMICOLON, Value: types.Literal{Value: types.SEMICOLON, Type: types.OPERATOR}})
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
					s.Tokens = append(s.Tokens, types.Token{TokenType: types.LITERAL, Value: types.Literal{Value: number, Type: types.INTEGER}})
				} else if number, err := strconv.ParseFloat(variable, 64); err == nil {
					s.Tokens = append(s.Tokens, types.Token{TokenType: types.LITERAL, Value: types.Literal{Value: number, Type: types.FLOAT64}})
				} else {
					log.Fatalf("cannot parse source file\n")
				}
			} else if unicode.IsLetter(rune(s.Content[s.CurrentPosition])) {
				before := s.CurrentPosition
				s.eatVar()
				after := s.CurrentPosition
				variable := s.Content[before:after]
				s.Tokens = append(s.Tokens, types.Token{TokenType: types.ID, Value: types.Literal{Value: variable, Type: types.STRING}})
			} else {
				log.Fatalf("unsupported type %v\n", string(s.Content[s.CurrentPosition]))
			}

		}
	}

	return &s
}
