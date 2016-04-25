package symtbl

import (
	"github.com/bantl23/cminus/log"
	"github.com/bantl23/cminus/syntree"
)

var GlbSymTblLst *SymTblLst
var curSymTblLst *SymTblLst
var GlbSymTblMap SymTblLstMap = make(SymTblLstMap)

func GlbPrintTableList() {
	PrintTableList(GlbSymTblLst, 0)
}

func GlbPrintTableMap() {
	PrintTableMap(GlbSymTblMap, 4)
}

func NewGlbSymTblLst() {
	GlbSymTblLst = NewSymTblLst(ROOT_SCOPE, nil)
	GlbSymTblMap[GlbSymTblLst.Scope()] = GlbSymTblLst
	curSymTblLst = GlbSymTblLst
	input := syntree.NewStmtFunctionInputNode()
	output := syntree.NewStmtFunctionOutputNode()
	InsertFuncInSymTbl(curSymTblLst, input)
	InsertFuncInSymTbl(curSymTblLst, output)
	GlbSymTblLst.SymTbl()[output.Name()].AddArg(INT_SYM_TYPE)
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
		curSymTblLst = NewSymTblLst(node.Name(), curSymTblLst)
		GlbSymTblMap[curSymTblLst.Scope()] = curSymTblLst
	}
	node.SetSymKey(curSymTblLst.Scope())
}

func postbuild(node syntree.Node) {
	log.AnalyzeLog.Printf("postbuild %+v", node)
	if node.IsFunc() || node.IsCompound() {
		curSymTblLst = curSymTblLst.Parent()
	}
}
