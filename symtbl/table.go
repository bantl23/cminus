package symtbl

import (
	"fmt"
	"strconv"
)

type SymTblVal struct {
	memLoc  MemLoc
	symType SymbolType
	size    int
	args    []SymbolType
	retType ReturnType
	lines   []int
}

func NewSymTblVal(memLoc MemLoc, symType SymbolType, size int, retType ReturnType, line int) *SymTblVal {
	s := new(SymTblVal)
	s.memLoc = memLoc
	s.symType = symType
	s.size = size
	s.retType = retType
	if line != -1 {
		s.lines = append(s.lines, line)
	}
	return s
}

type SymTbl map[string]*SymTblVal

func NewSymTbl() *SymTbl {
	s := make(SymTbl)
	return &s
}

func (s SymTblVal) MemLoc() MemLoc {
	return s.memLoc
}

func (s SymTblVal) SymType() SymbolType {
	return s.symType
}

func (s SymTblVal) Size() int {
	return s.size
}

func (s SymTblVal) Args() []SymbolType {
	return s.args
}

func (s *SymTblVal) AddArg(arg SymbolType) {
	s.args = append(s.args, arg)
}

func (s SymTblVal) RetType() ReturnType {
	return s.retType
}

func (s SymTblVal) Lines() []int {
	return s.lines
}

func (s *SymTblVal) AddLine(line int) {
	s.lines = append(s.lines, line)
}

func PrintTable(s *SymTbl, indent int) {
	if s != nil {
		for i := 0; i < indent; i++ {
			fmt.Print(" ")
		}
		fmt.Printf("Variable Name Memory Location Type Size RType Args Lines\n")
		for i := 0; i < indent; i++ {
			fmt.Print(" ")
		}
		fmt.Printf("============= =============== ==== ==== ===== ==== =====\n")
		for k, v := range *s {
			for i := 0; i < indent; i++ {
				fmt.Print(" ")
			}
			if v != nil {
				fmt.Printf("%13s 0x%013x %4s %4d %5s %+v %+v\n", k, v.MemLoc(), v.SymType(), v.Size(), v.RetType(), v.Args(), v.Lines())
			} else {
				fmt.Printf("%13s\n", k)
			}
		}
		fmt.Printf("\n")
	}
}

var SCOPE_SEPARATOR string = "$"
var ROOT_SCOPE string = "global"
var INNER_SCOPE string = "inner"
var INNER_COUNT int = 0
var ROOT_KEY string = SCOPE_SEPARATOR + ROOT_SCOPE

type SymTblLst struct {
	scope      string
	baseMemLoc MemLoc
	symTbl     SymTbl
	parent     *SymTblLst
	children   []*SymTblLst
}

type SymTblLstMap map[string]*SymTblLst

func PrintTableMap(s SymTblLstMap, indent int) {
	for k, v := range s {
		for i := 0; i < indent; i++ {
			fmt.Print(" ")
		}
		scope := "nil"
		if v.Parent() != nil {
			scope = v.Parent().Scope()
		}
		fmt.Printf("%s => %s\n", k, scope)
	}
}

func NewSymTblLst(scope string, parent *SymTblLst) *SymTblLst {
	s := new(SymTblLst)
	if scope == "" {
		scope = INNER_SCOPE + strconv.FormatInt(int64(INNER_COUNT), 10)
		INNER_COUNT++
	}
	s.scope = SCOPE_SEPARATOR + scope
	s.baseMemLoc = MemLoc(0)
	s.symTbl = *NewSymTbl()
	s.parent = parent
	s.children = nil
	if parent != nil {
		parent.AddChildren(s)
		s.scope = parent.Scope() + SCOPE_SEPARATOR + scope
	}
	return s
}

func (s SymTblLst) Scope() string {
	return s.scope
}

func (s SymTblLst) BaseMemLoc() MemLoc {
	return s.baseMemLoc
}

func (s *SymTblLst) IncBaseMemLoc(size int) {
	s.baseMemLoc = MemLoc(int(s.baseMemLoc) + size)
}

func (s SymTblLst) SymTbl() SymTbl {
	return s.symTbl
}

func (s SymTblLst) Parent() *SymTblLst {
	return s.parent
}

func (s *SymTblLst) SetParent(p *SymTblLst) {
	s.parent = p
}

func (s SymTblLst) Children() []*SymTblLst {
	return s.children
}

func (s *SymTblLst) AddChildren(c *SymTblLst) {
	s.children = append(s.children, c)
}

func (s *SymTblLst) HasId(variable string) bool {
	has := false
	lst := s
	for lst != nil {
		if _, ok := lst.symTbl[variable]; ok {
			has = true
			break
		}
		lst = lst.parent
	}
	return has
}

func (s *SymTblLst) GetMemLoc(variable string) MemLoc {
	memLoc := MemLoc(-1)
	lst := s
	for lst != nil {
		if val, ok := lst.symTbl[variable]; ok {
			memLoc = val.MemLoc()
			break
		}
		lst = lst.parent
	}
	return memLoc
}

func (s *SymTblLst) GetSize(variable string) int {
	size := -1
	lst := s
	for lst != nil {
		if val, ok := lst.symTbl[variable]; ok {
			size = val.Size()
			break
		}
		lst = lst.parent
	}
	return size
}

func (s *SymTblLst) GetIdType(variable string) SymbolType {
	typ := UNK_SYM_TYPE
	lst := s
	for lst != nil {
		if val, ok := lst.symTbl[variable]; ok {
			typ = val.SymType()
			break
		}
		lst = lst.parent
	}
	return typ
}

func (s *SymTblLst) GetIdArgs(variable string) []SymbolType {
	var typ []SymbolType
	lst := s
	for lst != nil {
		if val, ok := lst.symTbl[variable]; ok {
			typ = val.Args()
			break
		}
		lst = lst.parent
	}
	return typ
}

func PrintTableList(s *SymTblLst, indent int) {
	indent += 4
	if s != nil {
		for i := 0; i < indent; i++ {
			fmt.Printf(" ")
		}
		fmt.Printf("scope %s\n", s.scope)
		t := s.SymTbl()
		PrintTable(&t, indent)
		for _, v := range s.Children() {
			PrintTableList(v, indent)
		}
	}
	indent -= 4
}
