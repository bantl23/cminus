package symtbl

import (
	"github.com/bantl23/cminus/log"
	"github.com/bantl23/cminus/syntree"
)

var GlbSymTblLst *SymTblLst
var CurSymTblLst *SymTblLst
var TblSep string = "$"
var GblName string = "$global"
var InnerName string = "$inner"

func NewGlbSymTblLst() {
	GlbSymTblLst = new(SymTblLst)
	GlbSymTblLst.SymTbl = make(SymTbl)
	GlbSymTblLst.Name = GblName
	GlbSymTblLst.Prev = nil
	GlbSymTblLst.Next = nil
	input := syntree.NewStmtFunctionInputNode()
	output := syntree.NewStmtFunctionOutputNode()
	GlbSymTblLst.Insert(input)
	GlbSymTblLst.Insert(output)
	GlbSymTblLst.SymTbl[output.Name()].Args = append(GlbSymTblLst.SymTbl[output.Name()].Args, INT_SYM_TYPE)
	CurSymTblLst = GlbSymTblLst
}

func PrintGlbSymTblLst() {
	PrintTableList(GlbSymTblLst)
}

func Build(node syntree.Node) {
	syntree.Traverse(node, Pushin, Popout)
}

func Pushin(node syntree.Node) {
	if node.Save() == true {
		ok := CurSymTblLst.Insert(node)
		if ok == true {
			log.AnalyzeLog.Printf("inserted %+v into %+v", node, CurSymTblLst.Name)
		}
	} else {
		log.AnalyzeLog.Printf("received %+v", node)
	}
	if node.AddScope() == true {
		name := CurSymTblLst.Name
		n := NewSymTblLst(CurSymTblLst)
		CurSymTblLst = n
		if node.Name() != "" {
			CurSymTblLst.Name = name + TblSep + node.Name()
		} else {
			CurSymTblLst.Name = name + InnerName
		}
		log.AnalyzeLog.Printf("created %+v", CurSymTblLst.Name)
	}
}

func Popout(node syntree.Node) {
	if node.AddScope() == true {
		log.AnalyzeLog.Printf("returned %+v", CurSymTblLst.Name)
		CurSymTblLst = CurSymTblLst.Prev
	}
}
