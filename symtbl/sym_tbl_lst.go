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

var GlbSymTblLst *SymTblLst
var CurSymTblLst *SymTblLst
var TblSep string = "$"
var GblName string = "$global"
var InnerName string = "$inner"
var PrevDeclareName = ""
var LastDeclareName = "main"
var MaxInt = 2147483647
var MinInt = -2147483648
var MaxArrayInt = MaxInt
var MinArrayInt = 0
var FoundRet = false
var RetHasChild = false

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
	log.AnalyzeLog.Printf("check %+v", node)
	if PrevDeclareName == LastDeclareName {
		log.ErrorLog.Printf(">>>> Error main function must be the last declaration [%+v]", node.Pos())
	}
	if node.(syntree.Symbol).IsDeclaration() {
		PrevDeclareName = node.(syntree.Name).Name()
	}

	if node.(syntree.Symbol).IsArray() {
		if node.(syntree.Symbol).IsParam() == false {
			if node.(syntree.Value).Value() > MaxArrayInt {
				log.ErrorLog.Printf(">>>> Error array size %d is greater than %d [%+v]", node.(syntree.Value).Value(), MaxArrayInt, node.Pos())
			} else if node.(syntree.Value).Value() < MinArrayInt {
				log.ErrorLog.Printf(">>>> Error array size %d is less than %d [%+v]", node.(syntree.Value).Value(), MinArrayInt, node.Pos())
			}
		}
	}

	if node.(syntree.Symbol).IsReturn() {
		FoundRet = true
		if node.Children() != nil {
			RetHasChild = true
		} else {
			RetHasChild = false
		}
	}
	if node.(syntree.Symbol).IsFunc() {
		if FoundRet == true {
			if node.(syntree.ExpType).ExpType() == syntree.VOID_TYPE && RetHasChild == true {
				log.ErrorLog.Printf(">>>> Error void function returns a value [%+v]", node.Pos())
			} else if node.(syntree.ExpType).ExpType() == syntree.INTEGER_TYPE && RetHasChild == false {
				log.ErrorLog.Printf(">>>> Error non-void function has and empty return statement [%+v]", node.Pos())
			}
		} else {
			if node.(syntree.ExpType).ExpType() == syntree.INTEGER_TYPE {
				log.ErrorLog.Printf(">>>> Error non-void function does not have a return statement [%+v]", node.Pos())
			}
		}
		FoundRet = false
	}
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
