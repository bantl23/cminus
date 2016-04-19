package symtbl

type MemLoc int

func (m *MemLoc) Inc() {
	*m++
}

func (m *MemLoc) Reset() {
	*m = 0
}

func (m MemLoc) Get() int {
	return int(m)
}

func (m *MemLoc) Set(memLoc MemLoc) {
	*m = memLoc
}
