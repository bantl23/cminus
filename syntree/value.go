package syntree

type Value interface {
	Value() int
	SetValue(int)
}
