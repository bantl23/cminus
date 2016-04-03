package symtbl

type SymbolType int

const (
	FUNC_SYM_TYPE SymbolType = iota
	INT_SYM_TYPE
	ARR_SYM_TYPE
	UNK_SYM_TYPE
)

var symTypes = [...]string{
	"fun",
	"int",
	"arr",
	"unkSym",
}

func (symType SymbolType) String() string {
	return symTypes[symType]
}

type SymType interface {
	SymType() SymbolType
}
