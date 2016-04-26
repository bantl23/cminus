package symtbl

type IdentifierType int

const (
	VAR_ID_TYPE IdentifierType = iota
	PARAM_ID_TYPE
	UNK_ID_TYPE
)

var idTypes = [...]string{
	"var",
	"par",
	"id?",
}

func (idType IdentifierType) String() string {
	return idTypes[idType]
}

type IdType interface {
	IdType() IdentifierType
}
