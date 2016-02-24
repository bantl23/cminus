package scan

import (
	"io"
	"text/scanner"
)

type Lexer struct {
	scanner scanner.Scanner
}

func NewLexer(r io.Reader) *Lexer {
	l := new(Lexer)
	l.scanner.Init(r)
	return l
}

func (l *Lexer) Lex(lval *yySymType) int {
	var tok rune
	for tok != scanner.EOF {
		tok = l.scanner.Scan()
		fmt.Println("At", l.scanner.Pos(), ":", l.scanner.TokenText())
	}
	return NUM
}

func (l *Lexer) Error(e string) {
}
