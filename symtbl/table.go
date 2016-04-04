package symtbl

import (
	"fmt"
	"strconv"
)

type SymTblVal struct {
	memLoc  MemLoc
	symType SymbolType
	args    []SymbolType
	lines   []int
}

func NewSymTblVal(memLoc MemLoc, symType SymbolType, line int) *SymTblVal {
	s := new(SymTblVal)
	s.memLoc = memLoc
	s.symType = symType
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

func (s SymTblVal) Args() []SymbolType {
	return s.args
}

func (s *SymTblVal) AddArg(arg SymbolType) {
	s.args = append(s.args, arg)
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
		fmt.Printf("Variable Name Memory Location Type Args Lines\n")
		for i := 0; i < indent; i++ {
			fmt.Print(" ")
		}
		fmt.Printf("============= =============== ==== ==== =====\n")
		for k, v := range *s {
			for i := 0; i < indent; i++ {
				fmt.Print(" ")
			}
			if v != nil {
				fmt.Printf("%13s 0x%013x %4s %+v %+v\n", k, v.MemLoc(), v.SymType(), v.Args(), v.Lines())
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

type SymTblLst struct {
	scope    string
	symTbl   SymTbl
	parent   *SymTblLst
	children []*SymTblLst
}

func NewSymTblLst(scope string, parent *SymTblLst) *SymTblLst {
	s := new(SymTblLst)
	if scope == "" {
		scope = INNER_SCOPE + strconv.FormatInt(int64(INNER_COUNT), 10)
		INNER_COUNT++
	}
	s.scope = SCOPE_SEPARATOR + scope
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
