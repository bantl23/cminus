package opt

import (
	"github.com/bantl23/cminus/log"
	"github.com/bantl23/cminus/syntree"
)

func OptParseTree(node syntree.Node) {
	log.OptLog.Printf("=> Optimize Parse Tree\n")
	ConstantFolding(node)
}
