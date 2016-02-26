package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"unicode"
)

type LexerState int

const (
	LEX_START LexerState = iota
	LEX_COMMENT_FIRST
	LEX_COMMENT
	LEX_COMMENT_STAR
	LEX_NUM
	LEX_ID
	LEX_NEQ
	LEX_EQ
	LEX_LTE
	LEX_GTE
	LEX_DONE
)

type Lexer struct {
	file   *os.File
	reader *bufio.Reader
	row    int
	col    int
	txt    string
	tok    rune
	buf    []rune
	state  LexerState
}

func NewLexer(f *os.File) *Lexer {
	l := new(Lexer)
	l.file = f
	l.reader = bufio.NewReader(l.file)
	l.row = 1
	l.col = 0
	l.txt = ""
	l.tok = 0
	l.buf = make([]rune, 0)
	l.state = LEX_START
	return l
}

func (l *Lexer) String() string {
	return fmt.Sprintf("Lexer: %s, [%d:%d] (%s)", l.Name(), l.Row(), l.Col(), l.Text())
}

func (l *Lexer) Lex(lval *yySymType) int {
	var value int = 0
	l.state = LEX_START
	for l.state != LEX_DONE {
		var keep bool = true
		var err error
		l.tok, err = l.Read()
		if err == nil {

			if l.tok == '\n' {
				l.row++
				l.col = 0
			}
			l.col++

			switch l.state {
			case LEX_START:
				l.buf = make([]rune, 0)
				if unicode.IsDigit(l.tok) {
					l.state = LEX_NUM
				} else if unicode.IsLetter(l.tok) {
					l.state = LEX_ID
				} else if l.tok == '!' {
					l.state = LEX_NEQ
				} else if l.tok == '=' {
					l.state = LEX_EQ
				} else if l.tok == '<' {
					l.state = LEX_LTE
				} else if l.tok == '>' {
					l.state = LEX_GTE
				} else if l.tok == '/' {
					l.state = LEX_COMMENT_FIRST
				} else if unicode.IsSpace(l.tok) {
					keep = false
				} else {
					l.state = LEX_DONE
					switch l.tok {
					case '+':
						value = PLUS
					case '-':
						value = MINUS
					case '*':
						value = TIMES
					case ';':
						value = SEMI
					case '(':
						value = LPAREN
					case ')':
						value = RPAREN
					case '[':
						value = LBRACKET
					case ']':
						value = RBRACKET
					case '{':
						value = LBRACE
					case '}':
						value = RBRACE
					case ',':
						value = COMMA
					default:
						keep = false
						fmt.Printf("Error: unknown token [%+v]\n", l)
					}
				}
			case LEX_NUM:
				if unicode.IsDigit(l.tok) == false {
					l.Unread()
					l.state = LEX_DONE
					value = NUM
					keep = false
				}
			case LEX_ID:
				if unicode.IsLetter(l.tok) == false {
					l.Unread()
					l.state = LEX_DONE
					keep = false
					switch string(l.buf) {
					case "if":
						value = IF
					case "else":
						value = ELSE
					case "int":
						value = INT
					case "void":
						value = VOID
					case "while":
						value = WHILE
					default:
						value = ID
					}
				}
			case LEX_NEQ:
				l.state = LEX_DONE
				if l.tok == '=' {
					value = NEQ
				} else {
					keep = false
					fmt.Printf("Error: unknown token !%c\n", l.tok)
				}
			case LEX_EQ:
				l.state = LEX_DONE
				if l.tok == '=' {
					value = EQ
				} else {
					value = ASSIGN
				}
			case LEX_LTE:
				l.state = LEX_DONE
				if l.tok == '=' {
					value = LTE
				} else {
					value = LT
				}
			case LEX_GTE:
				l.state = LEX_DONE
				if l.tok == '=' {
					value = GTE
				} else {
					value = GT
				}
			case LEX_COMMENT_FIRST:
				if l.tok == '*' {
					l.state = LEX_COMMENT_STAR
				} else {
					l.state = LEX_DONE
					value = OVER
				}
			case LEX_COMMENT:
				if l.tok == '*' {
					l.state = LEX_COMMENT_STAR
				}
			case LEX_COMMENT_STAR:
				if l.tok == '/' {
					l.state = LEX_START
				} else if l.tok == '*' {
					l.state = LEX_COMMENT_STAR
				} else {
					l.state = LEX_COMMENT
				}
			case LEX_DONE:
			default:
				l.state = LEX_DONE
				keep = false
				fmt.Printf("Error: internal error, unknown state\n")
			}
		} else {
			l.state = LEX_DONE
			keep = false
			if err != io.EOF {
				fmt.Printf("Error: file read %s\n", err)
			}
		}

		if keep == true {
			l.buf = append(l.buf, l.tok)
		}

		if l.state == LEX_DONE {
			lval.yys = value
			lval.str = string(l.buf)
		}
	}
	return value
}

func (l *Lexer) Error(e string) {
	fmt.Printf("Error: %s [%+v]\n", e, l)
}

func (l *Lexer) Read() (rune, error) {
	r, _, err := l.reader.ReadRune()
	return r, err
}

func (l *Lexer) Unread() error {
	return l.reader.UnreadRune()
}

func (l *Lexer) Name() string {
	return l.file.Name()
}

func (l *Lexer) Row() int {
	return l.row
}

func (l *Lexer) Col() int {
	return l.col
}

func (l *Lexer) Text() string {
	return string(l.buf)
}
