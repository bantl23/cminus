package syntree

type ExpressionType int

const (
	VOID_TYPE ExpressionType = iota
	INTEGER_TYPE
	UNK_EXPRESSION_TYPE
)

var expressionTypes = [...]string{
	"VoidExpType",
	"IntExpType",
	"UnknownExpType",
}

func (expressionType ExpressionType) String() string {
	return expressionTypes[expressionType]
}

type ExpType interface {
	ExpType() ExpressionType
}