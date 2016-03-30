package syntree

import (
	"fmt"
)

type StmtIterationNode struct {
	NodeBase
}

func NewStmtIterationNode() Node {
	n := new(StmtIterationNode)
	n.position = Position{-1, -1}
	n.sibling = nil
	n.children = nil
	return n
}

func (n StmtIterationNode) String() string {
	return fmt.Sprintf("iteration [%+v]", n.Pos())
}
