package syntree

type Symbol interface {
	Save() bool
	AddScope() bool
}
