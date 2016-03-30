package syntree

import (
	"fmt"
)

type StmtSelectionNode struct {
	NodeBase
}

func NewStmtSelectionNode(row int, col int) Node {
	n := new(StmtSelectionNode)
	n.position = Position{row, col}
	n.sibling = nil
	n.children = nil
	return n
}

func (n StmtSelectionNode) String() string {
	return fmt.Sprintf("selection [%+v]", n.Pos())
}
