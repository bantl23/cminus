package syntree

type ExpressionType int

const (
	VOID_EXP_TYPE ExpressionType = iota
	INT_EXP_TYPE
	UNK_EXP_TYPE
)

var expressionTypes = [...]string{
	"VoidExp",
	"IntExp",
	"UnknownExp",
}

func (expressionType ExpressionType) String() string {
	return expressionTypes[expressionType]
}
