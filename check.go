package main

import (
	"github.com/bantl23/cminus/log"
	"github.com/bantl23/cminus/syntree"
)

func Analyze(node syntree.Node) {
	syntree.TraverseNode(node, precheck, postcheck)
}

var PrevDeclareName = ""
var LastDeclareName = "main"

func CheckMainLastDeclare(node syntree.Node) {
	if PrevDeclareName == LastDeclareName {
		log.ErrorLog.Printf(">>>>> Error main function must be the last declared function [%+v]", node.Pos())
	}
	if node.IsVar() || node.IsParam() || node.IsFunc() {
		PrevDeclareName = node.Name()
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
	if node.IsFunc() || node.IsCompound() {
		curSymTblLst = curSymTblLst.Parent()
	}
}
