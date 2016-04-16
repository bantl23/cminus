package symtbl

type ReturnType int

const (
	VOID_RET_TYPE ReturnType = iota
	INT_RET_TYPE
	UNK_RET_TYPE
)

var retTypes = [...]string{
	"void",
	"int",
	"ret?",
}

func (retType ReturnType) String() string {
	return retTypes[retType]
}

type RetType interface {
	RetType() ReturnType
}
