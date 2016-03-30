package symtbl

import (
	"fmt"
	"github.com/bantl23/cminus/log"
	"github.com/bantl23/cminus/syntree"
)

type MemoryLocation int

type Key struct {
	Scope    string
	Variable string
}

type Value struct {
	Lines  []int
	MemLoc MemoryLocation
	Next   *Key
}

type SymbolTable map[Key]*Value

var GlbMemLoc MemoryLocation = 0

func (m *MemoryLocation) Inc() {
	*m++
}

func (m *MemoryLocation) Reset() {
	*m = 0
}

func (m *MemoryLocation) Get() int {
	return int(*m)
}

func NewSymbolTable() *SymbolTable {
	s := make(SymbolTable)
	return &s
}

func (s *SymbolTable) Build(node syntree.Node) {
	syntree.Traverse(node, s.InsertNode, syntree.Nothing)
}

func (s *SymbolTable) Analyze(node syntree.Node) {
	syntree.Traverse(node, syntree.Nothing, s.CheckNode)
}

func (s *SymbolTable) PrintTable() {
	table := *s
	fmt.Printf("    Var Loc LineNumbers\n")
	fmt.Printf("    ===================\n")
	for i, e := range table {
		fmt.Printf("    %+v %+v %+v\n", i, e.MemLoc, e.Lines)
	}
}

func (s *SymbolTable) Insert(scope string, variable string, pos syntree.Position) {
	table := *s
	_, ok := table[Key{scope, variable}]
	if ok == true {
		table[Key{scope, variable}].Lines = append(table[Key{scope, variable}].Lines, pos.Row())
	} else {
		table[Key{scope, variable}] = new(Value)
		table[Key{scope, variable}].Lines = append(table[Key{scope, variable}].Lines, pos.Row())
		table[Key{scope, variable}].MemLoc = GlbMemLoc
		GlbMemLoc.Inc()
	}
}

func (s *SymbolTable) Obtain(scope string, variable string) MemoryLocation {
	table := *s
	_, ok := table[Key{scope, variable}]
	if ok == true {
		return table[Key{scope, variable}].MemLoc
	}
	return -1
}

func (s *SymbolTable) InsertNode(node syntree.Node) {
	log.AnalyzeLog.Printf("insert: %+v", node)
}

func (s *SymbolTable) CheckNode(node syntree.Node) {
}
