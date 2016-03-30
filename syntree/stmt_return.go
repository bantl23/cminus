package syntree

import (
	"fmt"
)

type StmtReturnNode struct {
	NodeBase
}

func NewStmtReturnNode() Node {
	n := new(StmtReturnNode)
	n.position = Position{-1, -1}
	n.sibling = nil
	n.children = nil
	return n
}

func (n StmtReturnNode) String() string {
	return fmt.Sprintf("return [%+v]", n.Pos())
}
