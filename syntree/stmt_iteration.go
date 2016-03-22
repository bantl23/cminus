package syntree

import (
	"fmt"
)

type StmtIterationNode struct {
	NodeBase
}

func NewStmtIterationNode() Node {
	n := new(StmtIterationNode)
	n.row = -1
	n.col = -1
	n.sibling = nil
	n.children = nil
	return n
}

func (n StmtIterationNode) String() string {
	row, col := n.Pos()
	return fmt.Sprintf("iteration [%+v:%+v]\n", row, col)
}
