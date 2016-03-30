package syntree

import (
	"fmt"
)

type StmtCompoundNode struct {
	NodeBase
}

func NewStmtCompoundNode() Node {
	n := new(StmtCompoundNode)
	n.position = Position{-1, -1}
	n.sibling = nil
	n.children = nil
	return n
}

func (n StmtCompoundNode) String() string {
	return fmt.Sprintf("compound [%+v]", n.Pos())
}
