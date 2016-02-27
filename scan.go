package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
)

type Lexer struct {
	file     *os.File
	reader   *bufio.Reader
	row      int
	column   int
	text     string
	tokName  string
	curLine  string
	readLine bool
	regex    *regexp.Regexp
	tokNames []string
}

func NewLexer(f *os.File) *Lexer {
	l := new(Lexer)
	l.file = f
	l.reader = bufio.NewReader(l.file)
	l.row = 0
	l.column = 0
	l.text = ""
	l.tokName = ""
	l.curLine = ""
	l.readLine = true
	l.regex = regexp.MustCompile("(?P<BEG_COMMENT>\\/\\*)|" +
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
	l.tokNames = l.regex.SubexpNames()
	return l
}

func (l *Lexer) String() string {
	return fmt.Sprintf("%s[%d:%d] %s %s", l.FileName(), l.Row(), l.Col(), l.Text(), l.Name())
}

func (l *Lexer) GetToken(value string, lval *yySymType) int {
	lval.yys = yyErrCode
	lval.str = l.Text()
	fmt.Printf("%s\n", lval.str)
	switch value {
	case "IF":
		lval.yys = IF
		return IF
	case "ELSE":
		lval.yys = ELSE
		return ELSE
	case "INT":
		lval.yys = INT
		return INT
	case "VOID":
		lval.yys = VOID
		return VOID
	case "WHILE":
		lval.yys = WHILE
		return WHILE
	case "NEQ":
		lval.yys = NEQ
		return NEQ
	case "EQ":
		lval.yys = EQ
		return EQ
	case "LTE":
		lval.yys = LTE
		return LTE
	case "GTE":
		lval.yys = GTE
		return GTE
	case "ASSIGN":
		lval.yys = ASSIGN
		return ASSIGN
	case "LT":
		lval.yys = LT
		return LT
	case "GT":
		lval.yys = GT
		return GT
	case "PLUS":
		lval.yys = PLUS
		return PLUS
	case "MINUS":
		lval.yys = MINUS
		return MINUS
	case "TIMES":
		lval.yys = TIMES
		return TIMES
	case "OVER":
		lval.yys = OVER
		return OVER
	case "LPAREN":
		lval.yys = LPAREN
		return LPAREN
	case "RPAREN":
		lval.yys = RPAREN
		return RPAREN
	case "LBRACKET":
		lval.yys = LBRACKET
		return LBRACKET
	case "RBRACKET":
		lval.yys = RBRACKET
		return RBRACKET
	case "LBRACE":
		lval.yys = LBRACE
		return LBRACE
	case "RBRACE":
		lval.yys = RBRACE
		return RBRACE
	case "COMMA":
		lval.yys = COMMA
		return COMMA
	case "SEMI":
		lval.yys = SEMI
		return SEMI
	case "NUM":
		lval.yys = NUM
		lval.str = l.Text()
		return NUM
	case "ID":
		lval.yys = ID
		lval.str = l.Text()
		return ID
	case "EOF":
		lval.yys = 0
		return 0
	default:
		lval.yys = yyErrCode
		fmt.Printf("Error unknown token")
		return 0
	}
	lval.yys = yyErrCode
	fmt.Printf("Error unknown")
	return 0
}

func (l *Lexer) Lex(lval *yySymType) int {
	in_comment := false
	keepProcessing := true
	for keepProcessing == true {
		var err error = nil
		if l.readLine == true {
			l.curLine, err = l.reader.ReadString('\n')
			if err == nil {
				fmt.Printf("%d: %s\n", l.row, strings.TrimSpace(l.curLine))
				l.curLine = strings.TrimSpace(l.curLine)
				l.row++
				l.readLine = false
			}
		}
		if err == nil {
			keepScanning := true
			for keepScanning == true {
				matches := l.regex.FindStringSubmatchIndex(l.curLine)
				if matches != nil {
					for i := 2; i < len(matches); i = i + 2 {
						if matches[i] != -1 {
							if l.tokNames[i/2] == "BEG_COMMENT" {
								in_comment = true
							}
							if in_comment == false {
								l.column = matches[i] + 1
								l.text = l.curLine[matches[i]:matches[i+1]]
								l.curLine = l.curLine[matches[i+1]:len(l.curLine)]
								l.tokName = l.tokNames[i/2]
								keepScanning = false
								keepProcessing = false
							} else {
								l.curLine = l.curLine[matches[i+1]:len(l.curLine)]
							}
							if l.tokNames[i/2] == "END_COMMENT" {
								in_comment = false
							}
							break
						}
					}
				} else {
					keepScanning = false
					l.readLine = true
				}
			}
		} else {
			if err == io.EOF {
				l.tokName = "EOF"
			} else {
				fmt.Printf("error reading from file")
			}
			keepProcessing = false
		}
	}
	return l.GetToken(l.tokName, lval)
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
	return l.column
}

func (l *Lexer) Text() string {
	return l.text
}

func (l *Lexer) Name() string {
	return l.tokName
}
