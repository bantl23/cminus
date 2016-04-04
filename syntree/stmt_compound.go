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
	n.symbolKey = ""
	return n
}

func (n StmtCompoundNode) IsCompound() bool {
	return true
}

func (n StmtCompoundNode) String() string {
	return fmt.Sprintf("compound [%+v]", n.Pos())
}
