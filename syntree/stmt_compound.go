package syntree

import (
	"fmt"
)

type StmtCompoundNode struct {
	NodeBase
}

func NewStmtCompoundNode(row int, col int) Node {
	n := new(StmtCompoundNode)
	n.position = Position{row, col}
	n.sibling = nil
	n.children = nil
	return n
}

func (n StmtCompoundNode) String() string {
	return fmt.Sprintf("compound [%+v]", n.Pos())
}
