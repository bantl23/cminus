package main

import (
	"github.com/bantl23/cminus/log"
	"github.com/bantl23/cminus/symtbl"
	"github.com/bantl23/cminus/syntree"
)

var glbMemLoc symtbl.MemLoc

func InsertVarParamInSymTbl(table *symtbl.SymTblLst, node syntree.Node) {
	log.AnalyzeLog.Printf("insert var/param: %+v %+v", table, node)
	if value, ok := table.SymTbl()[node.Name()]; ok {
		value.AddLine(node.Pos().Row())
	} else {
		t := symtbl.UNK_SYM_TYPE
		if node.IsArray() {
			t = symtbl.ARR_SYM_TYPE
		} else if node.IsInt() {
			t = symtbl.INT_SYM_TYPE
		}
		table.SymTbl()[node.Name()] = symtbl.NewSymTblVal(glbMemLoc, t, node.Pos().Row())
		glbMemLoc.Inc()
	}
}

func InsertFuncInSymTbl(table *symtbl.SymTblLst, node syntree.Node) {
	log.AnalyzeLog.Printf("insert func: %+v %+v", table, node)
	if value, ok := table.SymTbl()[node.Name()]; ok {
		value.AddLine(node.Pos().Row())
	} else {
		table.SymTbl()[node.Name()] = symtbl.NewSymTblVal(glbMemLoc, symtbl.FUNC_SYM_TYPE, node.Pos().Row())
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
