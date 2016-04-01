package symtbl

type SymbolType int

const (
	FUNCTION_TYPE SymbolType = iota
	INTEGER_TYPE
	ARRAY_TYPE
	UNK_SYMBOL_TYPE
)

var symTypes = [...]string{
	"function",
	"integer",
	"array",
	"unknown",
}

func (symType SymbolType) String() string {
	return symTypes[symType]
}

type SymType interface {
	SymType() SymbolType
}
