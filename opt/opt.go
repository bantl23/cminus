package opt

import (
	"github.com/bantl23/cminus/log"
	"github.com/bantl23/cminus/syntree"
)

func OptConstantFolding(node syntree.Node) {
	log.OptLog.Printf("=> Optimize constant folding")
	ConstantFolding(node)
	log.OptLog.Printf("<= Optimize constant folding")
}

func OptConstantPropagation(node syntree.Node) {
	log.OptLog.Printf("=> Optimize constant propagation")
	ConstantPropagation(node)
	log.OptLog.Printf("<= Optimize constant propagation")
}

func OptConstantFoldingAndConstantPropagation(node syntree.Node) {
	CONST_FOLDED = true
	CONST_PROPAGATED = true
	for CONST_FOLDED == true || CONST_PROPAGATED == true {
		CONST_FOLDED = false
		ConstantFolding(node)
		CONST_PROPAGATED = false
		ConstantPropagation(node)
	}
}
