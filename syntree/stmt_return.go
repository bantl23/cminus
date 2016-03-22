package syntree

import (
	"fmt"
)

type StmtReturnNode struct {
	NodeBase
}

func NewStmtReturnNode() Node {
	n := new(StmtReturnNode)
	n.row = -1
	n.col = -1
	n.sibling = nil
	n.children = nil
	return n
}

func (n StmtReturnNode) String() string {
	row, col := n.Pos()
	return fmt.Sprintf("return [%+v:%+v]\n", row, col)
}
