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
	"EndOfFile",
	"Error",
	"Else",
	"If",
	"Int",
	"Return",
	"Void",
	"While",
	"Plus",
	"Minus",
	"Times",
	"Over",
	"LessThan",
	"LessThanEqual",
	"GreaterThan",
	"GreaterThanEqual",
	"Equals",
	"NotEquals",
	"Semicolon",
	"Comma",
	"LeftParen",
	"RightParen",
	"LeftBracket",
	"RightBracket",
	"LeftBrace",
	"RightBrace",
	"UnknownTokenType",
}

func (tokenType TokenType) String() string {
	return tokenTypes[tokenType]
}

type TokType interface {
	TokType() TokenType
	SetTokType(TokenType)
}
