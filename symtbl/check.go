package symtbl

import (
	"github.com/bantl23/cminus/log"
	"github.com/bantl23/cminus/syntree"
)

var PrevDeclareName = ""
var LastDeclareName = "main"
var MaxInt = 2147483647
var MinInt = -2147483648
var MaxArrayInt = MaxInt
var MinArrayInt = 0
var FoundRet = false
var RetHasChild = false

func Analyze(node syntree.Node) {
	syntree.Traverse(node, PreCheck, PostCheck)
}

func PreCheck(node syntree.Node) {
	log.AnalyzeLog.Printf("precheck %+v", node)
}

func PostCheck(node syntree.Node) {
	log.AnalyzeLog.Printf("postcheck %+v", node)
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
