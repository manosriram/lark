package token

import (
	"lark/pkg/types"
	"log"
	"strconv"
	"unicode"
)

type Source struct {
	Content           string
	Tokens            []types.Token
	CurrentPosition   int
	TokenSoFar        string
	CurrentLineNumber int
}

func isnumeric(c rune) bool {
	return unicode.IsDigit(c)
}

func isalpha(c rune) bool {
	return unicode.IsLetter(c)
}

func (s *Source) getCurrentToken() byte {
	return s.Content[s.CurrentPosition]
}

func (s *Source) eatVar() string {
	var result string
	for s.CurrentPosition < len(s.Content) && (isalpha(rune(s.Content[s.CurrentPosition]))) {
		result += string(s.getCurrentToken())
		s.eatN(1)
	}
	return result
}

func (s *Source) eatNum() string {
	var result string
	for s.CurrentPosition < len(s.Content) && (isnumeric(rune(s.Content[s.CurrentPosition])) || string(s.getCurrentToken()) == types.DOT) {
		result += string(s.getCurrentToken())
		s.eatN(1)
	}
	return result
}
func (s *Source) closedUntil(untilLiteral byte) string {
	// var result string = string(s.getCurrentToken())
	var result string
	for s.CurrentPosition < len(s.Content) && s.getCurrentToken() != untilLiteral {
		result += string(s.getCurrentToken())
		s.eatN(1)
	}
	s.eatN(1)
	return result
}

func (s *Source) openUntil(untilLiteral byte) string {
	// var result string = string(s.getCurrentToken())
	var result string
	for s.CurrentPosition < len(s.Content) && s.getCurrentToken() != untilLiteral {
		result += string(s.getCurrentToken())
		s.eatN(1)
	}
	return result
}

func (s *Source) expect(expectedLiteral string) bool {
	current := s.CurrentPosition
	i := 0
	if len(s.Content)-current < len(expectedLiteral) {
		return false
	}
	for current < len(s.Content) && i < len(expectedLiteral) {
		if s.Content[current] != expectedLiteral[i] {
			return false
		}
		i++
		current++
	}
	return true
}

func (s *Source) peek(forward int) rune {
	for s.CurrentPosition < len(s.Content) {
		return rune(s.Content[s.CurrentPosition+forward])
	}
	return rune(0)
}

func (s *Source) eatN(n int) {
	done := false
	for s.CurrentPosition < len(s.Content) && (s.Content[s.CurrentPosition] == ' ' || s.Content[s.CurrentPosition] == '\n') {
		s.CurrentPosition += n
		done = true
	}
	if !done {
		s.CurrentPosition += n
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
		Content:           source,
		CurrentPosition:   0,
		Tokens:            make([]types.Token, 0),
		TokenSoFar:        "",
		CurrentLineNumber: 1,
	}
	for s.CurrentPosition < len(s.Content) {
		charAtPosition := s.getCurrentToken()
		switch charAtPosition {
		case '+':
			s.Tokens = append(s.Tokens, types.Token{TokenType: types.PLUS, Value: types.Literal{Value: '+', Type: types.OPERATOR}, LineNumber: s.CurrentLineNumber})
			s.eat()
		case '-':
			s.Tokens = append(s.Tokens, types.Token{TokenType: types.MINUS, Value: types.Literal{Value: '-', Type: types.OPERATOR}, LineNumber: s.CurrentLineNumber})
			s.eat()
		case '*':
			s.Tokens = append(s.Tokens, types.Token{TokenType: types.MULTIPLY, Value: types.Literal{Value: '*', Type: types.OPERATOR}, LineNumber: s.CurrentLineNumber})
			s.eat()
		case '/':
			if s.peek(1) == '/' {
				s.Tokens = append(s.Tokens, types.Token{TokenType: types.COMMENT, Value: types.Literal{Value: "//", Type: types.OPERATOR}, LineNumber: s.CurrentLineNumber})
				s.openUntil('\n')
			} else if s.peek(1) == '*' {
				s.Tokens = append(s.Tokens, types.Token{TokenType: types.COMMENT, Value: types.Literal{Value: "//", Type: types.OPERATOR}, LineNumber: s.CurrentLineNumber})

				s.closedUntil('*')
				s.closedUntil('/')
			} else {
				s.Tokens = append(s.Tokens, types.Token{TokenType: types.DIVIDE, Value: types.Literal{Value: '/', Type: types.OPERATOR}, LineNumber: s.CurrentLineNumber})
				s.eat()
			}
		case '=':
			s.eat()
			switch s.getCurrentToken() {
			case '=':
				s.Tokens = append(s.Tokens, types.Token{TokenType: types.EQUALS, Value: types.Literal{Value: "==", Type: types.OPERATOR}, LineNumber: s.CurrentLineNumber})
				s.eat()
			}
		case '(':
			s.Tokens = append(s.Tokens, types.Token{TokenType: types.LBRACE, Value: types.Literal{Value: '(', Type: types.OPERATOR}, LineNumber: s.CurrentLineNumber})
			s.eat()
		case ')':
			s.Tokens = append(s.Tokens, types.Token{TokenType: types.RBRACE, Value: types.Literal{Value: ')', Type: types.OPERATOR}, LineNumber: s.CurrentLineNumber})
			s.eat()
		case '{':
			s.Tokens = append(s.Tokens, types.Token{TokenType: types.LPAREN, Value: types.Literal{Value: '{', Type: types.OPERATOR}, LineNumber: s.CurrentLineNumber})
			s.eat()
		case '}':
			s.Tokens = append(s.Tokens, types.Token{TokenType: types.RPAREN, Value: types.Literal{Value: '}', Type: types.OPERATOR}, LineNumber: s.CurrentLineNumber})
			s.eat()
		case '!':
			s.eat()
			switch s.getCurrentToken() {
			case '=':
				s.Tokens = append(s.Tokens, types.Token{TokenType: types.NOT_EQUAL, Value: types.Literal{Value: "!=", Type: types.OPERATOR}, LineNumber: s.CurrentLineNumber})
				s.eat()
			default:
				s.Tokens = append(s.Tokens, types.Token{TokenType: types.NOT, Value: types.Literal{Value: "!", Type: types.OPERATOR}, LineNumber: s.CurrentLineNumber})
			}
		case '[':
			s.Tokens = append(s.Tokens, types.Token{TokenType: types.FUNCTION_ARGUMENT_OPEN, Value: types.Literal{Value: "[", Type: types.OPERATOR}, LineNumber: s.CurrentLineNumber})
			s.eat()
		case ']':
			s.Tokens = append(s.Tokens, types.Token{TokenType: types.FUNCTION_ARGUMENT_CLOSE, Value: types.Literal{Value: "]", Type: types.OPERATOR}, LineNumber: s.CurrentLineNumber})
			s.eat()
		case '<':
			s.eat()
			switch s.getCurrentToken() {
			case '<':
				s.Tokens = append(s.Tokens, types.Token{TokenType: types.FUNCTION_OPEN, Value: types.Literal{Value: "<<", Type: types.OPERATOR}, LineNumber: s.CurrentLineNumber})
				s.eat()
			case '=':
				s.Tokens = append(s.Tokens, types.Token{TokenType: types.LESSER_OR_EQUAL, Value: types.Literal{Value: "<=", Type: types.OPERATOR}, LineNumber: s.CurrentLineNumber})
				s.eat()
			case '-':
				if s.peek(1) == '>' {
					s.Tokens = append(s.Tokens, types.Token{TokenType: types.SWAP, Value: types.Literal{Value: "<->", Type: types.OPERATOR}, LineNumber: s.CurrentLineNumber})
				} else {
					s.Tokens = append(s.Tokens, types.Token{TokenType: types.ASSIGN, Value: types.Literal{Value: '=', Type: types.OPERATOR}, LineNumber: s.CurrentLineNumber})
				}
				s.eat()
			default:
				s.Tokens = append(s.Tokens, types.Token{TokenType: types.LESSER, Value: types.Literal{Value: "<", Type: types.OPERATOR}, LineNumber: s.CurrentLineNumber})
			}
		case '>':
			s.eat()
			switch s.getCurrentToken() {
			case '>':
				s.Tokens = append(s.Tokens, types.Token{TokenType: types.FUNCTION_CLOSE, Value: types.Literal{Value: ">>", Type: types.OPERATOR}, LineNumber: s.CurrentLineNumber})
				s.eat()
			case '=':
				s.Tokens = append(s.Tokens, types.Token{TokenType: types.GREATER_OR_EQUAL, Value: types.Literal{Value: ">=", Type: types.OPERATOR}, LineNumber: s.CurrentLineNumber})
				s.eat()
			default:
				s.Tokens = append(s.Tokens, types.Token{TokenType: types.GREATER, Value: types.Literal{Value: ">", Type: types.OPERATOR}, LineNumber: s.CurrentLineNumber})
			}
		case '"':
			s.eat()
			variable := s.openUntil(byte('"'))
			s.eat()
			s.Tokens = append(s.Tokens, types.Token{TokenType: types.LITERAL, Value: types.Literal{Value: variable, Type: types.STRING}, LineNumber: s.CurrentLineNumber})
			break
		case ';':
			s.Tokens = append(s.Tokens, types.Token{TokenType: types.SEMICOLON, Value: types.Literal{Value: types.SEMICOLON, Type: types.OPERATOR}, LineNumber: s.CurrentLineNumber})
			s.eat()
			break
		case ' ', '\t':
			s.eat()
			break
		case '\n':
			s.CurrentLineNumber++
			s.eat()
			break
		default:
			if unicode.IsNumber(rune(s.getCurrentToken())) {
				variable := s.eatNum()
				if number, err := strconv.Atoi(variable); err == nil {
					s.Tokens = append(s.Tokens, types.Token{TokenType: types.LITERAL, Value: types.Literal{Value: number, Type: types.INTEGER}, LineNumber: s.CurrentLineNumber})
				} else if number, err := strconv.ParseFloat(variable, 64); err == nil {
					s.Tokens = append(s.Tokens, types.Token{TokenType: types.LITERAL, Value: types.Literal{Value: number, Type: types.FLOAT64}, LineNumber: s.CurrentLineNumber})
				} else {
					log.Fatalf("syntax error at line %d\n", s.CurrentLineNumber)
				}
			} else if unicode.IsLetter(rune(s.getCurrentToken())) {
				variable := s.eatVar()
				switch variable {
				case "true":
					s.Tokens = append(s.Tokens, types.Token{TokenType: types.LITERAL, Value: types.Literal{Value: true, Type: types.BOOLEAN}, LineNumber: s.CurrentLineNumber})
				case "false":
					s.Tokens = append(s.Tokens, types.Token{TokenType: types.LITERAL, Value: types.Literal{Value: false, Type: types.BOOLEAN}, LineNumber: s.CurrentLineNumber})
				case "else":
					s.Tokens = append(s.Tokens, types.Token{TokenType: types.ELSE, Value: types.Literal{Value: "if", Type: types.STATEMENT}, LineNumber: s.CurrentLineNumber})
				case "if":
					s.Tokens = append(s.Tokens, types.Token{TokenType: types.IF, Value: types.Literal{Value: "if", Type: types.STATEMENT}, LineNumber: s.CurrentLineNumber})
				case "fn":
					s.Tokens = append(s.Tokens, types.Token{TokenType: types.FUNCTION, Value: types.Literal{Value: "fn", Type: types.STATEMENT}, LineNumber: s.CurrentLineNumber})
				case "ret":
					s.Tokens = append(s.Tokens, types.Token{TokenType: types.FUNCTION_RETURN, Value: types.Literal{Value: "return", Type: types.STATEMENT}, LineNumber: s.CurrentLineNumber})
				default:
					switch s.getCurrentToken() {
					case '(':
						s.eat()
						if s.getCurrentToken() == ')' {
							s.eat()
							s.Tokens = append(s.Tokens, types.Token{TokenType: types.FUNCTION_CALL, Value: types.Literal{Value: variable, Type: types.STATEMENT}, LineNumber: s.CurrentLineNumber})
							// s.eat()
						} else {
							log.Fatalf("expected ')'")
						}
						break
					default:
						s.Tokens = append(s.Tokens, types.Token{TokenType: types.ID, Value: types.Literal{Value: variable, Type: types.STRING}, LineNumber: s.CurrentLineNumber})
						break
					}
				}
			} else {
				log.Fatalf("unsupported type %v at line %d\n", string(s.getCurrentToken()), s.CurrentLineNumber)
			}
		}
	}

	return &s
}
