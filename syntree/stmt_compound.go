package syntree

import (
	"fmt"
)

type StmtCompoundNode struct {
	NodeBase
}

func NewStmtCompoundNode() Node {
	n := new(StmtCompoundNode)
	n.row = -1
	n.col = -1
	n.sibling = nil
	n.children = nil
	return n
}

func (n StmtCompoundNode) String() string {
	row, col := n.Pos()
	return fmt.Sprintf("compound [%+v:%+v]", row, col)
}
