package syntree

import (
	"fmt"
)

type ExpAssignNode struct {
	NodeBase
}

func NewExpAssignNode(row int, col int) Node {
	n := new(ExpAssignNode)
	n.position = Position{row, col}
	n.sibling = nil
	n.children = nil
	return n
}

func (n ExpAssignNode) String() string {
	return fmt.Sprintf("assign [%+v]", n.Pos())
}
