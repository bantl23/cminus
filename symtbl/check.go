package symtbl

import (
	"github.com/bantl23/cminus/log"
	"github.com/bantl23/cminus/syntree"
)

func Analyze(node syntree.Node) {
	CurSymTblLst = GlbSymTblLst
	syntree.Traverse(node, PreCheck, PostCheck)
}

var PrevDeclareName = ""
var LastDeclareName = "main"

func CheckMainLastDeclaration(node syntree.Node) {
	if PrevDeclareName == LastDeclareName {
		log.ErrorLog.Printf(">>>> Error main function must be the last declaration[%+v]", node.Pos())
	}
	if node.(syntree.Symbol).IsDeclaration() {
		PrevDeclareName = node.(syntree.Name).Name()
	}
}

var MaxArrayInt = 2147483647
var MinArrayInt = 0

func CheckArrayIndexSize(node syntree.Node) {
	if node.(syntree.Symbol).IsArray() {
		if node.(syntree.Symbol).IsParam() == false {
			if node.(syntree.Value).Value() > MaxArrayInt {
				log.ErrorLog.Printf(">>>> Error array size %d is greater than %d [%+v]", node.(syntree.Value).Value(), MaxArrayInt, node.Pos())
			} else if node.(syntree.Value).Value() < MinArrayInt {
				log.ErrorLog.Printf(">>>> Error array size %d is less than %d [%+v]", node.(syntree.Value).Value(), MinArrayInt, node.Pos())
			}
		}
	}
}

var FoundReturn = false
var ReturnHasChild = false

func CheckReturnValue(node syntree.Node) {
	if node.(syntree.Symbol).IsReturn() {
		FoundReturn = true
		if node.Children() != nil {
			ReturnHasChild = true
		} else {
			ReturnHasChild = false
		}
	}
	if node.(syntree.Symbol).IsFunc() {
		if FoundReturn == true {
			if node.(syntree.ExpType).ExpType() == syntree.VOID_TYPE && ReturnHasChild == true {
				log.ErrorLog.Printf(">>>> Error void function returns a value [%+v]", node.Pos())
			} else if node.(syntree.ExpType).ExpType() == syntree.INTEGER_TYPE && ReturnHasChild == false {
				log.ErrorLog.Printf(">>>> Error non-void function has empty return statement [%+v]", node.Pos())
			}
		} else {
			if node.(syntree.ExpType).ExpType() == syntree.INTEGER_TYPE {
				log.ErrorLog.Printf(">>>> Error non-void function does not have a return statement [%+v]", node.Pos())
			}
		}
		FoundReturn = false
	}
}

var SibCnt int = 0

func CheckFuncArgs(node syntree.Node) {
	if node.(syntree.Symbol).IsCall() {
		SibCnt = 0
		if node.Children()[0] != nil {
			SibCnt = 1
			sibling := node.Children()[0].Sibling()
			log.AnalyzeLog.Printf("child %+v sib %+v", node.Children()[0], sibling)
			for sibling != nil {
				sibling = sibling.Sibling()
				SibCnt++
			}
		}
		glbArgs := GlbSymTblLst.SymTbl[node.(syntree.Name).Name()].Args
		if len(glbArgs) != SibCnt {
			log.ErrorLog.Printf(">>>> Error calling %s with %d arguments but takes %d arguments [%+v]", node.(syntree.Name).Name(), SibCnt, len(glbArgs), node.Pos())
		}
	}
}

var GlbNextCnt int = 0

func PreCheck(node syntree.Node) {
	log.AnalyzeLog.Printf("precheck %+v", node)
	if node.(syntree.Symbol).AddScope() {
		if len(CurSymTblLst.Next) > 1 {
			CurSymTblLst = CurSymTblLst.Next[GlbNextCnt]
			GlbNextCnt++
		} else {
			CurSymTblLst = CurSymTblLst.Next[0]
		}
	}
}

func PostCheck(node syntree.Node) {
	log.AnalyzeLog.Printf("postcheck %+v", node)
	CheckMainLastDeclaration(node)
	CheckArrayIndexSize(node)
	CheckReturnValue(node)
	CheckFuncArgs(node)
	if node.(syntree.Symbol).AddScope() {
		CurSymTblLst = CurSymTblLst.Prev
	}
}
