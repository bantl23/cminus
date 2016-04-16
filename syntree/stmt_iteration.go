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
	n.symbolKey = ""
	return n
}

func (n StmtIterationNode) IsStmt() bool {
	return true
}

func (n StmtIterationNode) IsIteration() bool {
	return true
}

func (n StmtIterationNode) String() string {
	return fmt.Sprintf("iteration [%+v]", n.Pos())
}
