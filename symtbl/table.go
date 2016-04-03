package symtbl

import (
	"fmt"
)

type SymTblVal struct {
	memLoc  MemLoc
	symType SymbolType
	args    []SymbolType
	lines   []int
}

type SymTbl map[string]*SymTblVal

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
	for s != nil {
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
			fmt.Printf("%13s 0x%013x %4s %+v %+v\n", k, v.MemLoc(), v.SymType(), v.Args(), v.Lines())
		}
	}
}

type SymTblLst struct {
	scope    string
	symTbl   SymTbl
	parent   *SymTblLst
	children []*SymTblLst
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
	for s != nil {
		for i := 0; i < indent; i++ {
			fmt.Printf(" ")
		}
		fmt.Printf("Scope %s\n", s.scope)
		t := s.SymTbl()
		PrintTable(&t, indent)
		for _, v := range s.Children() {
			PrintTableList(v, indent)
		}
	}
	indent -= 4
}
