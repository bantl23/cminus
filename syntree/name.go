package syntree

type NameGet interface {
	Name() string
}

type NameSet interface {
	SetName(string)
}
