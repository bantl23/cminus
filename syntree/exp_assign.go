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
	n.symbolKey = ""
	return n
}

func (n ExpAssignNode) IsExp() bool {
	return true
}

func (n ExpAssignNode) IsAssign() bool {
	return true
}

func (n ExpAssignNode) String() string {
	return fmt.Sprintf("assign [%+v]", n.Pos())
}
