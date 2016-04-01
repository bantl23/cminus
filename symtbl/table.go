package symtbl

import (
	"fmt"
	"github.com/bantl23/cminus/log"
	"github.com/bantl23/cminus/syntree"
)

type Value struct {
	MemLoc  MemLoc
	SymType SymbolType
	Lines   []int
}

type SymTbl map[string]*Value

type SymTblLst struct {
	SymTbl SymTbl
	Name   string
	Prev   *SymTblLst
	Next   []*SymTblLst
}

var GlbSymTblLst *SymTblLst
var CurSymTblLst *SymTblLst
var TblSep string = "$"
var GblName string = "$global"
var InnerName string = "$inner"

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
	GlbSymTblLst.Name = GblName
	GlbSymTblLst.Prev = nil
	GlbSymTblLst.Next = nil
	input := syntree.NewStmtFunctionInputNode()
	output := syntree.NewStmtFunctionOutputNode()
	GlbSymTblLst.Insert(input)
	GlbSymTblLst.Insert(output)
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
		fmt.Printf("    Scope %+v level %+v\n", lst.Name, depth)
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
		ok := CurSymTblLst.Insert(node)
		if ok == true {
			log.AnalyzeLog.Printf("inserted %+v into %+v", node, CurSymTblLst.Name)
		}
	} else {
		log.AnalyzeLog.Printf("received %+v", node)
	}
	if node.(syntree.Symbol).AddScope() == true {
		name := CurSymTblLst.Name
		n := NewSymTblLst(CurSymTblLst)
		CurSymTblLst = n
		if _, ok := node.(syntree.Name); ok {
			CurSymTblLst.Name = name + TblSep + node.(syntree.Name).Name()
		} else {
			CurSymTblLst.Name = name + InnerName
		}
		log.AnalyzeLog.Printf("created %+v", CurSymTblLst.Name)
	}
}

func Popout(node syntree.Node) {
	if node.(syntree.Symbol).AddScope() == true {
		log.AnalyzeLog.Printf("returned %+v", CurSymTblLst.Name)
		CurSymTblLst = CurSymTblLst.Prev
	}
}

func Check(node syntree.Node) {
}

func (s *SymTbl) PrintTable() {
	fmt.Printf("    Variable Name Type Memory Location Lines\n")
	fmt.Printf("    ============= ==== =============== =====\n")
	for i, e := range *s {
		fmt.Printf("    %+v\t %+v 0x%08x\t %+v\n", i, e.SymType, e.MemLoc, e.Lines)
	}
	fmt.Printf("\n")
}

func (s *SymTblLst) Insert(node syntree.Node) bool {
	inserted := false

	table := (*s).SymTbl
	variable := node.(syntree.Name).Name()
	line := node.Pos().Row()
	symType := UNK_SYMBOL_TYPE

	if node.(syntree.Symbol).IsFunc() {
		symType = FUNCTION_TYPE
	} else if node.(syntree.Symbol).IsArray() {
		symType = ARRAY_TYPE
	} else if node.(syntree.Symbol).IsInt() {
		symType = INTEGER_TYPE
	}

	_, ok := table[variable]
	if node.(syntree.Symbol).IsDeclaration() {
		if ok == true {
			log.ErrorLog.Printf(">>>> ERROR %+v already declared", variable)
		} else {
			table[variable] = new(Value)
			if line != -1 {
				table[variable].Lines = append(table[variable].Lines, line)
			}
			table[variable].SymType = symType
			table[variable].MemLoc = glbMemLoc
			glbMemLoc.Inc()
			inserted = true
		}
	} else {
		if ok == true {
			table[variable].Lines = append(table[variable].Lines, line)
			inserted = true
		} else {
			lst := s
			for lst != nil {
				tbl := lst.SymTbl
				_, pOk := tbl[variable]
				if pOk == true {
					tbl[variable].Lines = append(tbl[variable].Lines, line)
					inserted = true
					break
				}
				log.AnalyzeLog.Printf("%+v", lst.Name)
				lst = lst.Prev
			}
			if inserted == false {
				log.ErrorLog.Printf(">>>> ERROR: %+v not declared", variable)
			}
		}
	}

	return inserted
}

func (s *SymTblLst) Obtain(scope string, variable string) MemLoc {
	table := (*s).SymTbl
	_, ok := table[variable]
	if ok == true {
		return table[variable].MemLoc
	}
	return -1
}
