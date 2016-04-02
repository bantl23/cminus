package syntree

import (
	"fmt"
)

type StmtReturnNode struct {
	NodeBase
}

func NewStmtReturnNode(row int, col int) Node {
	n := new(StmtReturnNode)
	n.position = Position{row, col}
	n.sibling = nil
	n.children = nil
	return n
}

func (n StmtReturnNode) IsReturn() bool {
	return true
}

func (n StmtReturnNode) String() string {
	return fmt.Sprintf("return [%+v]", n.Pos())
}
