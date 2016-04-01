package symtbl

import (
	"fmt"
	"github.com/bantl23/cminus/log"
	"github.com/bantl23/cminus/syntree"
)

type Value struct {
	MemLoc MemLoc
	Lines  []int
}

type SymTbl map[string]*Value

type SymTblLst struct {
	SymTbl SymTbl
	Prev   *SymTblLst
	Next   []*SymTblLst
}

var CurSymTblLst *SymTblLst
var GlbSymTblLst *SymTblLst

func NewSymTblLst(prev *SymTblLst) *SymTblLst {
	s := new(SymTblLst)
	s.SymTbl = make(SymTbl)
	s.Next = nil
	s.Prev = prev
	prev.Next = append(prev.Next, s)
	return s
}

func NewGlbSymTblLst() {
	GlbSymTblLst = new(SymTblLst)
	GlbSymTblLst.SymTbl = make(SymTbl)
	GlbSymTblLst.Prev = nil
	GlbSymTblLst.Next = nil
	GlbSymTblLst.SymTbl.Insert("input", -1)
	GlbSymTblLst.SymTbl.Insert("output", -1)
	CurSymTblLst = GlbSymTblLst
}

func PrintSymTblLst() {
	PrintTableList(GlbSymTblLst)
}

var depth = -1

func PrintTableList(lst *SymTblLst) {
	depth++
	if lst != nil {
		tbl := lst.SymTbl
		fmt.Printf("    Scope level %+v\n", depth)
		tbl.PrintTable()
		for _, t := range lst.Next {
			PrintTableList(t)
		}
	}
	depth--
}

func Build(node syntree.Node) {
	syntree.Traverse(node, Insert, Popout)
}

func Analyze(node syntree.Node) {
	syntree.Traverse(node, syntree.Nothing, Check)
}

func Insert(node syntree.Node) {
	if node.(syntree.Symbol).Save() == true {
		CurSymTblLst.SymTbl.Insert(node.(syntree.Name).Name(), node.Pos().Row())
		log.AnalyzeLog.Printf("inserted %+v", node)
	}
	if node.(syntree.Symbol).AddScope() == true {
		n := NewSymTblLst(CurSymTblLst)
		CurSymTblLst = n
		log.AnalyzeLog.Printf("added scoped symbol table")
	}
}

func Popout(node syntree.Node) {
	if node.(syntree.Symbol).AddScope() == true {
		CurSymTblLst = CurSymTblLst.Prev
		log.AnalyzeLog.Printf("returned from scoped symbol table")
	}
}

func Check(node syntree.Node) {
}

func (s *SymTbl) PrintTable() {
	fmt.Printf("    Variable Name Type Memory Location Lines\n")
	fmt.Printf("    ============= ==== =============== =====\n")
	for i, e := range *s {
		fmt.Printf("    %+v\t 0x%08x\t %+v\n", i, e.MemLoc, e.Lines)
	}
	fmt.Printf("\n")
}

func (s *SymTbl) Insert(variable string, line int) {
	table := *s
	_, ok := table[variable]
	if ok == true {
		table[variable].Lines = append(table[variable].Lines, line)
	} else {
		table[variable] = new(Value)
		if line != -1 {
			table[variable].Lines = append(table[variable].Lines, line)
		}
		table[variable].MemLoc = glbMemLoc
		glbMemLoc.Inc()
	}
}

func (s *SymTbl) Obtain(scope string, variable string) MemLoc {
	table := *s
	_, ok := table[variable]
	if ok == true {
		return table[variable].MemLoc
	}
	return -1
}
