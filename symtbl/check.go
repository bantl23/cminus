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

var ArgLst []SymbolType = nil

func CheckFuncArgsPre(node syntree.Node) {
	if node.(syntree.Symbol).IsCall() {
		ArgLst = nil
	}
}

func CheckFuncArgsPost(node syntree.Node) {
	if CurSymTblLst.ObtainSymType(node.(syntree.Name).Name()) == ARRAY_TYPE {
		ArgLst = append(ArgLst, ARRAY_TYPE)
	} else if CurSymTblLst.ObtainSymType(node.(syntree.Name).Name()) == INTEGER_TYPE {
		ArgLst = append(ArgLst, INTEGER_TYPE)
	}
	if node.(syntree.Symbol).IsCall() {
		glbArgs := GlbSymTblLst.SymTbl[node.(syntree.Name).Name()].Args
		if ArgLst == nil && glbArgs != nil {
			log.ErrorLog.Printf(">>>> Error calling %s with no arguments %d required [%+v]", node.(syntree.Name).Name(), len(glbArgs), node.Pos())
		} else if ArgLst != nil && glbArgs == nil {
			log.ErrorLog.Printf(">>>> Error calling %s with arguments but takes no arguments [%+v]", node.(syntree.Name).Name(), node.Pos())
		} else if len(ArgLst) != len(glbArgs) {
			log.ErrorLog.Printf(">>>> Error calling %s with %d arguments but takes %d arguments [%+v]", node.(syntree.Name).Name(), len(ArgLst), len(glbArgs), node.Pos())
		} else {
			for i := range ArgLst {
				if ArgLst[i] != glbArgs[i] {
					log.ErrorLog.Printf(">>>> Error calling %s argument %d receiving %s but expects %s [%+v]", node.(syntree.Name).Name(), i+1, ArgLst[i], glbArgs[i], node.Pos())
				}
			}
		}
	}
}

var Count int = 0

func PreCheck(node syntree.Node) {
	log.AnalyzeLog.Printf("precheck %+v", node)
	if node.(syntree.Symbol).AddScope() {
		if len(CurSymTblLst.Next) > 1 {
			CurSymTblLst = CurSymTblLst.Next[Count]
			Count++
		} else {
			CurSymTblLst = CurSymTblLst.Next[0]
		}
	}
	CheckFuncArgsPre(node)
}

func PostCheck(node syntree.Node) {
	log.AnalyzeLog.Printf("postcheck %+v", node)
	CheckMainLastDeclaration(node)
	CheckArrayIndexSize(node)
	CheckReturnValue(node)
	CheckFuncArgsPost(node)
	if node.(syntree.Symbol).AddScope() {
		CurSymTblLst = CurSymTblLst.Prev
	}
}
