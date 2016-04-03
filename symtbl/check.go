package symtbl

import (
	"github.com/bantl23/cminus/log"
	"github.com/bantl23/cminus/syntree"
)

func Analyze(node syntree.Node) {
	syntree.Traverse(node, PreCheck, PostCheck)
}

var PrevDeclareName = ""
var LastDeclareName = "main"

func CheckMainLastDeclaration(node syntree.Node) {
	if PrevDeclareName == LastDeclareName {
		log.ErrorLog.Printf(">>>> Error main function must be the last declaration[%+v]", node.Pos())
	}
	if node.IsDecl() {
		PrevDeclareName = node.Name()
	}
}

var MaxArrayInt = 2147483647
var MinArrayInt = 0

func CheckArrayIndexSize(node syntree.Node) {
	if node.IsArray() {
		if node.IsParam() == false {
			if node.Value() > MaxArrayInt {
				log.ErrorLog.Printf(">>>> Error array size %d is greater than %d [%+v]", node.Value(), MaxArrayInt, node.Pos())
			} else if node.Value() < MinArrayInt {
				log.ErrorLog.Printf(">>>> Error array size %d is less than %d [%+v]", node.Value(), MinArrayInt, node.Pos())
			}
		}
	}
}

var FoundReturn = false
var ReturnHasChild = false

func CheckReturnValue(node syntree.Node) {
	if node.IsReturn() {
		FoundReturn = true
		if node.Children() != nil {
			ReturnHasChild = true
		} else {
			ReturnHasChild = false
		}
	}
	if node.IsFunc() {
		if FoundReturn == true {
			if node.ExpType() == syntree.VOID_EXP_TYPE && ReturnHasChild == true {
				log.ErrorLog.Printf(">>>> Error void function returns a value [%+v]", node.Pos())
			} else if node.ExpType() == syntree.INT_EXP_TYPE && ReturnHasChild == false {
				log.ErrorLog.Printf(">>>> Error non-void function has empty return statement [%+v]", node.Pos())
			}
		} else {
			if node.ExpType() == syntree.INT_EXP_TYPE {
				log.ErrorLog.Printf(">>>> Error non-void function does not have a return statement [%+v]", node.Pos())
			}
		}
		FoundReturn = false
	}
}

var SibCnt int = 0

func CheckFuncArgs(node syntree.Node) {
	if node.IsCall() {
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
		glbArgs := GlbSymTblLst.SymTbl[node.Name()].Args
		if len(glbArgs) != SibCnt {
			log.ErrorLog.Printf(">>>> Error calling %s with %d arguments but takes %d arguments [%+v]", node.Name(), SibCnt, len(glbArgs), node.Pos())
		}
	}
}

func PreCheck(node syntree.Node) {
	log.AnalyzeLog.Printf("precheck %+v", node)
}

func PostCheck(node syntree.Node) {
	log.AnalyzeLog.Printf("postcheck %+v", node)
	CheckMainLastDeclaration(node)
	CheckArrayIndexSize(node)
	CheckReturnValue(node)
	CheckFuncArgs(node)
}
