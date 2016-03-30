package syntree

import (
	"fmt"
)

type StmtIterationNode struct {
	NodeBase
}

func NewStmtIterationNode(row int, col int) Node {
	n := new(StmtIterationNode)
	n.position = Position{row, col}
	n.sibling = nil
	n.children = nil
	return n
}

func (n StmtIterationNode) String() string {
	return fmt.Sprintf("iteration [%+v]", n.Pos())
}
