package symtbl

import (
	"github.com/bantl23/cminus/log"
	"github.com/bantl23/cminus/syntree"
)

func InsertFuncInSymTbl(table *SymTblLst, node syntree.Node) {
	log.AnalyzeLog.Printf("insert func: %+v %+v", table.Scope(), node)
	if value, ok := table.SymTbl()[node.Name()]; ok {
		value.AddLine(node.Pos().Row())
	} else {
		r := UNK_RET_TYPE
		if node.ExpType() == syntree.VOID_EXP_TYPE {
			r = VOID_RET_TYPE
		} else if node.ExpType() == syntree.INT_EXP_TYPE {
			r = INT_RET_TYPE
		}
		memLoc := table.BaseMemLoc()
		table.SymTbl()[node.Name()] = NewSymTblVal(memLoc, FUNC_SYM_TYPE, 1, r, node.Pos().Row())
		table.IncBaseMemLoc()

		if len(node.Children()) > 0 {
			n := node.Children()[0]
			for n != nil {
				if n.ExpType() != syntree.VOID_EXP_TYPE {
					if n.IsArray() {
						table.SymTbl()[node.Name()].AddArg(ARR_SYM_TYPE)
					} else if n.IsInt() {
						table.SymTbl()[node.Name()].AddArg(INT_SYM_TYPE)
					}
				}
				n = n.Sibling()
			}
		}
	}
}

func InsertVarParamInSymTbl(table *SymTblLst, node syntree.Node) {
	log.AnalyzeLog.Printf("insert var/param: %+v %+v", table.Scope(), node)
	if value, ok := table.SymTbl()[node.Name()]; ok {
		value.AddLine(node.Pos().Row())
	} else {
		t := UNK_SYM_TYPE
		size := 1
		if node.IsArray() {
			t = ARR_SYM_TYPE
			size = node.Value()
		} else if node.IsInt() {
			t = INT_SYM_TYPE
			size = 1
		}
		memLoc := table.BaseMemLoc()
		table.SymTbl()[node.Name()] = NewSymTblVal(memLoc, t, size, UNK_RET_TYPE, node.Pos().Row())
		table.IncBaseMemLoc()
	}
}

func InsertCallIdInSymTbl(table *SymTblLst, node syntree.Node) {
	log.AnalyzeLog.Printf("insert call/id: %+v %+v", table.Scope(), node)
	t := table
	for t != nil {
		if value, ok := t.SymTbl()[node.Name()]; ok {
			value.AddLine(node.Pos().Row())
			break
		}
		t = t.Parent()
	}
}
