package syntree

import (
	"fmt"
)

type ExpAssignNode struct {
	NodeBase
}

func NewExpAssignNode() Node {
	n := new(ExpAssignNode)
	n.row = -1
	n.col = -1
	n.sibling = nil
	n.children = nil
	return n
}

func (n ExpAssignNode) String() string {
	row, col := n.Pos()
	return fmt.Sprintf("assign [%+v:%+v]\n", row, col)
}
