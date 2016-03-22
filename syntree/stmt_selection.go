package syntree

import (
	"fmt"
)

type StmtSelectionNode struct {
	NodeBase
}

func NewStmtSelectionNode() Node {
	n := new(StmtSelectionNode)
	n.row = -1
	n.col = -1
	n.sibling = nil
	n.children = nil
	return n
}

func (n StmtSelectionNode) String() string {
	row, col := n.Pos()
	return fmt.Sprintf("selection [%+v:%+v]\n", row, col)
}
