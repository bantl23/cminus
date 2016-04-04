package main

import (
	"github.com/bantl23/cminus/log"
	"github.com/bantl23/cminus/symtbl"
	"github.com/bantl23/cminus/syntree"
	"os"
)

var CheckErr bool = false

func Analyze(node syntree.Node) {
	syntree.TraverseNode(node, precheck, postcheck)
	if CheckErr == true {
		os.Exit(1)
	}
}

var PrevDeclareName = ""
var LastDeclareName = "main"

func CheckMainLastDeclare(node syntree.Node) {
	if PrevDeclareName == LastDeclareName {
		log.ErrorLog.Printf(">>>>> Error main function must be the last declared function [%+v]", node.Pos())
		CheckErr = true
	}
	if node.IsVar() || node.IsParam() || node.IsFunc() {
		PrevDeclareName = node.Name()
	}
}

var MaxArrayInt = 2147483647
var MinArrayInt = 0

func CheckArrayIndexSize(node syntree.Node) {
	if node.IsArray() && node.IsParam() == false {
		if node.Value() > MaxArrayInt {
			log.ErrorLog.Printf(">>>>> Error array size %d is greater than %d [%+v]", node.Value(), MaxArrayInt, node.Pos())
			CheckErr = true
		} else if node.Value() < MinArrayInt {
			log.ErrorLog.Printf(">>>>> Error array size %d is less than %d [%+v]", node.Value(), MinArrayInt, node.Pos())
			CheckErr = true
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
				log.ErrorLog.Printf(">>>>> Error void func returns a value [%+v]", node.Pos())
				CheckErr = true
			} else if node.ExpType() == syntree.INT_EXP_TYPE && ReturnHasChild == false {
				log.ErrorLog.Printf(">>>>> Error non-void func has empty return statement [%+v]", node.Pos())
				CheckErr = true
			}
		} else {
			if node.ExpType() == syntree.INT_EXP_TYPE {
				log.ErrorLog.Printf(">>>>> Error non-void func does not have a return statement [%+v]", node.Pos())
				CheckErr = true
			}
		}
		FoundReturn = false
	}
}

func CheckFuncTypes(node syntree.Node) {
	if node.IsCall() {
		if node.Children() != nil {
			sib := node.Children()[0]
			var sibSlice []symtbl.SymbolType
			for sib != nil {
				sym := GlbSymTblMap[sib.SymKey()].GetIdType(sib.Name())
				sibSlice = append(sibSlice, sym)
				sib = sib.Sibling()
			}
			funSlice := GlbSymTblMap[node.SymKey()].GetIdArgs(node.Name())
			if sibSlice == nil && funSlice != nil {
				log.ErrorLog.Printf(">>>>> Error calling %s with 0 arguments but takes %d arguments [%+v]", node.Name(), len(funSlice), node.Pos())
				CheckErr = true
			} else if sibSlice != nil && funSlice == nil {
				log.ErrorLog.Printf(">>>>> Error calling %s with %d arguments but takes 0 arguments [%+v]", node.Name(), len(sibSlice), node.Pos())
				CheckErr = true
			} else if len(sibSlice) != len(funSlice) {
				log.ErrorLog.Printf(">>>>> Error calling %s with %d arguments but takes %d arguments [%+v]", node.Name(), len(sibSlice), len(funSlice), node.Pos())
				CheckErr = true
			} else {
				for i := range sibSlice {
					if sibSlice[i] != funSlice[i] {
						log.ErrorLog.Printf(">>>>> Error calling %s with args %+v but expecting %+v [%+v]", node.Name(), sibSlice, funSlice, node.Pos())
						CheckErr = true
						break
					}
				}
			}
		}
	}
}

func precheck(node syntree.Node) {
	log.AnalyzeLog.Printf("precheck %+v %+v", curSymTblLst.Scope(), node)
	if node.IsFunc() || node.IsCompound() {
		curSymTblLst = curSymTblLst.Children()[0]
	}
}

func postcheck(node syntree.Node) {
	log.AnalyzeLog.Printf("postcheck %+v %+v", curSymTblLst.Scope(), node)

	CheckMainLastDeclare(node)
	CheckArrayIndexSize(node)
	CheckReturnValue(node)
	CheckFuncTypes(node)

	if node.IsFunc() || node.IsCompound() {
		curSymTblLst = curSymTblLst.Parent()
	}
}
