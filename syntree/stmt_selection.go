package syntree

import (
	"fmt"
)

type StmtSelectionNode struct {
	NodeBase
}

func NewStmtSelectionNode() Node {
	n := new(StmtSelectionNode)
	n.position = Position{-1, -1}
	n.sibling = nil
	n.children = nil
	return n
}

func (n StmtSelectionNode) String() string {
	return fmt.Sprintf("selection [%+v]", n.Pos())
}
