package main

import (
	"fmt"
	"os"
	"regexp"
	"text/scanner"
)

type Lexer struct {
	file       *os.File
	scan       scanner.Scanner
	pos        scanner.Position
	text       string
	tok        rune
	id         *regexp.Regexp
	num        *regexp.Regexp
	commentBeg string
	commentEnd string
}

func NewLexer(f *os.File) *Lexer {
	l := new(Lexer)
	l.file = f
	l.scan.Init(l.file)
	l.pos = l.scan.Pos()
	l.text = l.scan.TokenText()
	l.tok = yyEofCode
	l.id = regexp.MustCompile("[a-zA-Z]+")
	l.num = regexp.MustCompile("[0-9]+")
	l.commentBeg = "/*"
	l.commentEnd = "*/"
	return l
}

func (l *Lexer) readValue() {
	l.pos = l.scan.Pos()
	l.text = l.scan.TokenText()
	l.tok = l.scan.Scan()
}

func (l *Lexer) Lex(lval *yySymType) int {
	l.readValue()
	if l.tok != scanner.EOF {
		if l.text == l.commentBeg {
			for l.text != l.commentEnd || l.tok != scanner.EOF {
				l.readValue()
			}
			l.readValue()
		}
	}

	fmt.Println(l)

	if l.tok != scanner.EOF {
		return l.GetTok(lval, l.text)
	}
	return yyEofCode
}

func (l *Lexer) Error(e string) {
	fmt.Println("Error: ", l)
}

func (l Lexer) String() string {
	return fmt.Sprintf("%s [%+v] (%s)", l.file.Name(), l.scan.Pos(), l.scan.TokenText())
}

func (l Lexer) GetTok(lval *yySymType, str string) int {
	switch str {
	case "if":
		lval.yys = IF
		lval.str = str
		return IF
	case "else":
		lval.yys = ELSE
		lval.str = str
		return ELSE
	case "int":
		lval.yys = INT
		lval.str = str
		return INT
	case "return":
		lval.yys = RETURN
		lval.str = str
		return RETURN
	case "void":
		lval.yys = VOID
		lval.str = str
		return VOID
	case "while":
		lval.yys = WHILE
		lval.str = str
		return WHILE
	case "+":
		lval.yys = PLUS
		lval.str = str
		return PLUS
	case "-":
		lval.yys = MINUS
		lval.str = str
		return MINUS
	case "*":
		lval.yys = TIMES
		lval.str = str
		return TIMES
	case "/":
		lval.yys = OVER
		lval.str = str
		return OVER
	case "<":
		lval.yys = LT
		lval.str = str
		return LT
	case "<=":
		lval.yys = LTE
		lval.str = str
		return LTE
	case ">":
		lval.yys = GT
		lval.str = str
		return GT
	case ">=":
		lval.yys = GTE
		lval.str = str
		return GTE
	case "==":
		lval.yys = EQ
		lval.str = str
		return EQ
	case "!=":
		lval.yys = NEQ
		lval.str = str
		return NEQ
	case "=":
		lval.yys = ASSIGN
		lval.str = str
		return ASSIGN
	case ";":
		lval.yys = SEMI
		lval.str = str
		return SEMI
	case ",":
		lval.yys = COMMA
		lval.str = str
		return COMMA
	case "(":
		lval.yys = LPAREN
		lval.str = str
		return LPAREN
	case ")":
		lval.yys = RPAREN
		lval.str = str
		return RPAREN
	case "{":
		lval.yys = LBRACE
		lval.str = str
		return LBRACE
	case "}":
		lval.yys = RBRACE
		lval.str = str
		return RBRACE
	case "[":
		lval.yys = LBRACKET
		lval.str = str
		return LBRACKET
	case "]":
		lval.yys = RBRACKET
		lval.str = str
		return RBRACKET
	default:
		if l.id.MatchString(str) {
			lval.yys = ID
			lval.str = str
			return ID
		} else if l.num.MatchString(str) {
			lval.yys = NUM
			lval.str = str
			return NUM
		}
	}
	lval.yys = yyErrCode
	lval.str = "Invalid token"
	return yyErrCode
}
