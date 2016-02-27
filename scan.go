package main

import (
	"bufio"
	"fmt"
	"github.com/bantl23/cminus/log"
	"io"
	"os"
	"regexp"
	"strings"
)

type Lexer struct {
	file       *os.File
	reader     *bufio.Reader
	row        int
	column     int
	text       string
	tokName    string
	curLine    string
	prevColumn int
	readLine   bool
	regex      *regexp.Regexp
	tokNames   []string
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
	l.prevColumn = 0
	l.readLine = true
	l.regex = regexp.MustCompile("(?P<BEG_COMMENT>\\/\\*)|" +
		"(?P<END_COMMENT>\\*\\/)|" +
		"(?P<IF>if)|" +
		"(?P<ELSE>else)|" +
		"(?P<INT>int)|" +
		"(?P<VOID>void)|" +
		"(?P<WHILE>while)|" +
		"(?P<RETURN>return)|" +
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
	switch value {
	case "IF":
		lval.yys = IF
	case "ELSE":
		lval.yys = ELSE
	case "INT":
		lval.yys = INT
	case "VOID":
		lval.yys = VOID
	case "WHILE":
		lval.yys = WHILE
	case "RETURN":
		lval.yys = RETURN
	case "NEQ":
		lval.yys = NEQ
	case "EQ":
		lval.yys = EQ
	case "LTE":
		lval.yys = LTE
	case "GTE":
		lval.yys = GTE
	case "ASSIGN":
		lval.yys = ASSIGN
	case "LT":
		lval.yys = LT
	case "GT":
		lval.yys = GT
	case "PLUS":
		lval.yys = PLUS
	case "MINUS":
		lval.yys = MINUS
	case "TIMES":
		lval.yys = TIMES
	case "OVER":
		lval.yys = OVER
	case "LPAREN":
		lval.yys = LPAREN
	case "RPAREN":
		lval.yys = RPAREN
	case "LBRACKET":
		lval.yys = LBRACKET
	case "RBRACKET":
		lval.yys = RBRACKET
	case "LBRACE":
		lval.yys = LBRACE
	case "RBRACE":
		lval.yys = RBRACE
	case "COMMA":
		lval.yys = COMMA
	case "SEMI":
		lval.yys = SEMI
	case "NUM":
		lval.yys = NUM
		lval.str = l.Text()
	case "ID":
		lval.yys = ID
		lval.str = l.Text()
	case "EOF":
		lval.yys = 0
	default:
		lval.yys = yyErrCode
		log.Error.Printf("Unknown token\n")
	}
	return lval.yys
}

func (l *Lexer) Lex(lval *yySymType) int {
	in_comment := false
	keepProcessing := true
	for keepProcessing == true {
		var err error = nil
		if l.readLine == true {
			l.curLine, err = l.reader.ReadString('\n')
			if err == nil {
				log.Echo.Printf("%d: %s\n", l.row, strings.TrimSpace(l.curLine))
				l.curLine = strings.TrimSpace(l.curLine)
				l.row++
				l.column = 1
				l.prevColumn = 0
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
								l.column = l.column + matches[i] + l.prevColumn
								l.prevColumn = matches[i+1] - matches[i]
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
				log.Error.Printf("File read %s\n", err)
			}
			keepProcessing = false
		}
	}
	return l.GetToken(l.tokName, lval)
}

func (l *Lexer) Error(e string) {
	log.Error.Printf("%s [%+v]\n", e, l)
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
