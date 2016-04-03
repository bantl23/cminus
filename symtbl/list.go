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
		tbl.Print()
		for _, t := range lst.Next {
			PrintTableList(t)
		}
	}
	depth--
}

var PrevFunc = ""

func (s *SymTblLst) Insert(node syntree.Node) bool {
	inserted := false

	table := (*s).SymTbl
	variable := node.Name()
	line := node.Pos().Row()
	symType := UNK_SYM_TYPE

	if node.IsFunc() {
		symType = FUNC_SYM_TYPE
		PrevFunc = node.Name()
	} else if node.IsArray() {
		symType = ARR_SYM_TYPE
	} else if node.IsInt() {
		symType = INT_SYM_TYPE
	}

	_, ok := table[variable]
	if node.IsDecl() {
		if ok == true {
			log.ErrorLog.Printf(">>>> ERROR %+v already declared [%+v]", variable, node.Pos())
		} else {
			table[variable] = new(Value)
			if line != -1 {
				table[variable].Lines = append(table[variable].Lines, line)
			}
			table[variable].SymType = symType
			table[variable].MemLoc = glbMemLoc
			if node.IsParam() {
				GlbSymTblLst.SymTbl[PrevFunc].Args = append(GlbSymTblLst.SymTbl[PrevFunc].Args, symType)
			}
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

func (s *SymTblLst) ObtainMemLog(variable string) MemLoc {
	lst := s
	for lst != nil {
		table := lst.SymTbl
		_, ok := table[variable]
		if ok == true {
			return table[variable].MemLoc
		} else {
			lst = lst.Prev
		}
	}
	return -1
}

func (s *SymTblLst) ObtainSymType(variable string) SymbolType {
	lst := s
	for lst != nil {
		table := lst.SymTbl
		_, ok := table[variable]
		if ok == true {
			return table[variable].SymType
		} else {
			lst = lst.Prev
		}
	}
	return UNK_SYM_TYPE
}
