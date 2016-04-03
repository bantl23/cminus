package syntree

type TokenType int

const (
	ENDFILE TokenType = iota
	ERROR
	ELSE
	IF
	INT
	RETURN
	VOID
	WHILE
	PLUS
	MINUS
	TIMES
	OVER
	LT
	LTE
	GT
	GTE
	EQ
	NEQ
	SEMI
	COMMA
	LPAREN
	RPAREN
	LBRACKET
	RBRACKET
	LBRACE
	RBRACE
	UNK_TOKEN_TYPE
)

var tokenTypes = [...]string{
	"EOF",
	"ERROR",
	"else",
	"if",
	"int",
	"return",
	"void",
	"while",
	"+",
	"-",
	"*",
	"/",
	"<",
	"<=",
	">",
	">=",
	"==",
	"!=",
	";",
	",",
	"(",
	")",
	"[",
	"]",
	"{",
	"}",
	"unkTok",
}

func (tokenType TokenType) String() string {
	return tokenTypes[tokenType]
}
