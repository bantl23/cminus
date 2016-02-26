package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
)

var regex *regexp.Regexp = regexp.MustCompile("(?P<BEG_COMMENT>\\/\\*)|" +
	"(?P<END_COMMENT>\\*\\/)|" +
	"(?P<IF>if)|" +
	"(?P<ELSE>else)`|" +
	"(?P<INT>int)|" +
	"(?P<VOID>void)|" +
	"(?P<WHILE>while)|" +
	"(?P<NEQ>!=)|" +
	"(?P<EQ>==)|" +
	"(?P<LTE><=)|" +
	"(?P<GTE>>=)|" +
	"(?P<ASSIGN>=)|" +
	"(?P<LT><)|" +
	"(?P<GT>>)|" +
	"(?P<PLUS>\\+)|" +
	"(?P<MINUS>-)|" +
	"(?P<TIMES>\\*)|" +
	"(?P<OVER>\\/)|" +
	"(?P<LPAREN>\\()|" +
	"(?P<RPAREN>\\))|" +
	"(?P<LBRACKET>\\[)|" +
	"(?P<RBRACKET>\\])|" +
	"(?P<LBRACE>\\{)|" +
	"(?P<RBRACE>\\})|" +
	"(?P<COMMA>,)|" +
	"(?P<SEMI>;)|" +
	"(?P<NUM>[0-9]+)|" +
	"(?P<ID>[a-zA-Z]+)|" +
	"(?P<WHITESPACE>[ \\t]+)" +
	"(?P<NEWLINE>\\n)")

var names []string = regex.SubexpNames()

type Lexer struct {
	file   *os.File
	reader *bufio.Reader
	row    int
	col    int
	txt    string
	name   string
	line   string
}

func NewLexer(f *os.File) *Lexer {
	l := new(Lexer)
	l.file = f
	l.reader = bufio.NewReader(l.file)
	l.row = 0
	l.col = 0
	l.txt = ""
	l.name = ""
	l.line = ""
	return l
}

func (l *Lexer) String() string {
	return fmt.Sprintf("%s[%d:%d] %s %s", l.FileName(), l.Row(), l.Col(), l.Text(), l.Name())
}

func (l *Lexer) GetToken(value string) int {
	switch value {
	case "IF":
		return IF
	case "ELSE":
		return ELSE
	case "INT":
		return INT
	case "VOID":
		return VOID
	case "WHILE":
		return WHILE
	case "NEQ":
		return NEQ
	case "EQ":
		return EQ
	case "LTE":
		return LTE
	case "GTE":
		return GTE
	case "ASSIGN":
		return ASSIGN
	case "LT":
		return LT
	case "GT":
		return GT
	case "PLUS":
		return PLUS
	case "MINUS":
		return MINUS
	case "TIMES":
		return TIMES
	case "OVER":
		return OVER
	case "LPAREN":
		return LPAREN
	case "RPAREN":
		return RPAREN
	case "LBRACKET":
		return LBRACKET
	case "RBRACKET":
		return RBRACKET
	case "LBRACE":
		return LBRACE
	case "RBRACE":
		return RBRACE
	case "COMMA":
		return COMMA
	case "SEMI":
		return SEMI
	case "NUM":
		return NUM
	case "ID":
		return ID
	default:
		fmt.Printf("Error unknown token")
		return 0
	}
	return 0
}

func (l *Lexer) Lex(lval *yySymType) int {

	in_comment := false
	keepReading := true
	for keepReading == true {
		var err error
		l.line, err = l.reader.ReadString('\n')
		l.row++
		if err == nil {
			keepLining := true
			for keepLining == true {
				matches := regex.FindStringSubmatchIndex(l.line)
				if matches != nil {
					for i := 2; i < len(matches); i = i + 2 {
						if matches[i] != -1 {
							if names[i/2] == "BEG_COMMENT" {
								in_comment = true
							}
							if in_comment == false {
								l.col = matches[i] + 1
								l.txt = l.line[matches[i]:matches[i+1]]
								l.line = l.line[matches[i+1]:len(l.line)]
								l.name = names[i/2]
								fmt.Printf("token: %+v\n", l)
							} else {
								l.line = l.line[matches[i+1]:len(l.line)]
							}
							if names[i/2] == "END_COMMENT" {
								in_comment = false
							}
							break
						}
					}
				} else {
					keepLining = false
				}
			}
		} else {
			if err != io.EOF {
				fmt.Printf("error reading from file")
			}
			keepReading = false
		}
	}

	return l.GetToken(l.name)
}

func (l *Lexer) Error(e string) {
	fmt.Printf("Error: %s [%+v]\n", e, l)
}

func (l *Lexer) FileName() string {
	return l.file.Name()
}

func (l *Lexer) Row() int {
	return l.row
}

func (l *Lexer) Col() int {
	return l.col
}

func (l *Lexer) Text() string {
	return l.txt
}

func (l *Lexer) Name() string {
	return l.name
}
