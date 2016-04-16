package main

import (
	"github.com/bantl23/cminus/log"
	"github.com/bantl23/cminus/symtbl"
	"github.com/bantl23/cminus/syntree"
)

var glbMemLoc symtbl.MemLoc

func InsertFuncInSymTbl(table *symtbl.SymTblLst, node syntree.Node) {
	log.AnalyzeLog.Printf("insert func: %+v %+v", table.Scope(), node)
	if value, ok := table.SymTbl()[node.Name()]; ok {
		value.AddLine(node.Pos().Row())
	} else {
		r := symtbl.UNK_RET_TYPE
		if node.ExpType() == syntree.VOID_EXP_TYPE {
			r = symtbl.VOID_RET_TYPE
		} else if node.ExpType() == syntree.INT_EXP_TYPE {
			r = symtbl.INT_RET_TYPE
		}
		table.SymTbl()[node.Name()] = symtbl.NewSymTblVal(glbMemLoc, symtbl.FUNC_SYM_TYPE, r, node.Pos().Row())
		glbMemLoc.Inc()

		if len(node.Children()) > 0 {
			n := node.Children()[0]
			for n != nil {
				if n.ExpType() != syntree.VOID_EXP_TYPE {
					if n.IsArray() {
						table.SymTbl()[node.Name()].AddArg(symtbl.ARR_SYM_TYPE)
					} else if n.IsInt() {
						table.SymTbl()[node.Name()].AddArg(symtbl.INT_SYM_TYPE)
					}
				}
				n = n.Sibling()
			}
		}
	}
}

func InsertVarParamInSymTbl(table *symtbl.SymTblLst, node syntree.Node) {
	log.AnalyzeLog.Printf("insert var/param: %+v %+v", table.Scope(), node)
	if value, ok := table.SymTbl()[node.Name()]; ok {
		value.AddLine(node.Pos().Row())
	} else {
		t := symtbl.UNK_SYM_TYPE
		if node.IsArray() {
			t = symtbl.ARR_SYM_TYPE
		} else if node.IsInt() {
			t = symtbl.INT_SYM_TYPE
		}
		table.SymTbl()[node.Name()] = symtbl.NewSymTblVal(glbMemLoc, t, symtbl.UNK_RET_TYPE, node.Pos().Row())
		glbMemLoc.Inc()
	}
}

func InsertCallIdInSymTbl(table *symtbl.SymTblLst, node syntree.Node) {
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
