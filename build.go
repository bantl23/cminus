package main

import (
	"github.com/bantl23/cminus/log"
	"github.com/bantl23/cminus/symtbl"
	"github.com/bantl23/cminus/syntree"
)

var GlbSymTblLst *symtbl.SymTblLst
var curSymTblLst *symtbl.SymTblLst

func PrintTableList() {
	symtbl.PrintTableList(GlbSymTblLst, 0)
}

func NewGlbSymTblLst() {
	GlbSymTblLst = symtbl.NewSymTblLst(symtbl.ROOT_SCOPE, nil)
	symtbl.SymTblLstMap[GlbSymTblLst.Scope()] = GlbSymTblLst
	curSymTblLst = GlbSymTblLst
	input := syntree.NewStmtFunctionInputNode()
	output := syntree.NewStmtFunctionOutputNode()
	InsertFuncInSymTbl(curSymTblLst, input)
	InsertFuncInSymTbl(curSymTblLst, output)
	GlbSymTblLst.SymTbl()[output.Name()].AddArg(symtbl.INT_SYM_TYPE)
}

func BuildTableList(node syntree.Node) {
	NewGlbSymTblLst()
	syntree.TraverseNode(node, prebuild, postbuild)
}

func prebuild(node syntree.Node) {
	log.AnalyzeLog.Printf("prebuild %+v", node)
	if node.IsFunc() {
		InsertFuncInSymTbl(curSymTblLst, node)
	} else if node.IsCall() {
		InsertCallIdInSymTbl(curSymTblLst, node)
	} else if node.IsId() {
		InsertCallIdInSymTbl(curSymTblLst, node)
	} else if node.IsVar() && node.IsInt() {
		InsertVarParamInSymTbl(curSymTblLst, node)
	} else if node.IsParam() && node.IsInt() {
		InsertVarParamInSymTbl(curSymTblLst, node)
	}
	if node.IsFunc() || node.IsCompound() {
		curSymTblLst = symtbl.NewSymTblLst(node.Name(), curSymTblLst)
		symtbl.SymTblLstMap[curSymTblLst.Scope()] = curSymTblLst
	}
}

func postbuild(node syntree.Node) {
	log.AnalyzeLog.Printf("postbuild %+v", node)
	if node.IsFunc() || node.IsCompound() {
		curSymTblLst = curSymTblLst.Parent()
	}
}
