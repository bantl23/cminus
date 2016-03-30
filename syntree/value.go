package syntree

type ValueGet interface {
	Value() int
}

type ValueSet interface {
	SetValue(int)
}
