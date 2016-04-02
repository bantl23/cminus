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

func PreCheck(node syntree.Node) {
	log.AnalyzeLog.Printf("precheck %+v", node)
}

func PostCheck(node syntree.Node) {
	log.AnalyzeLog.Printf("postcheck %+v", node)
	CheckMainLastDeclaration(node)
	CheckArrayIndexSize(node)
	CheckReturnValue(node)
}
