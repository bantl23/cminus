package symtbl

import (
	"fmt"
	"github.com/bantl23/cminus/log"
	"github.com/bantl23/cminus/syntree"
)

type SymTblLst struct {
	SymTbl SymTbl
	Name   string
	Prev   *SymTblLst
	Next   []*SymTblLst
}

func NewSymTblLst(prev *SymTblLst) *SymTblLst {
	s := new(SymTblLst)
	s.SymTbl = make(SymTbl)
	s.Next = nil
	s.Prev = prev
	prev.Next = append(prev.Next, s)
	return s
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
			log.ErrorLog.Printf(">>>> ERROR %+v already declared [%+v]", variable, node.Pos())
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
				log.ErrorLog.Printf(">>>> ERROR: %+v not declared [%+v]", variable, node.Pos())
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
