package syntree

type Symbol interface {
	Save() bool
	AddScope() bool
	IsFunc() bool
	IsArray() bool
	IsInt() bool
	IsDeclaration() bool
	IsReturn() bool
	IsParam() bool
	IsCall() bool
}
