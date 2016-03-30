package syntree

import (
	"fmt"
)

type ExpAssignNode struct {
	NodeBase
}

func NewExpAssignNode() Node {
	n := new(ExpAssignNode)
	n.position = Position{-1, -1}
	n.sibling = nil
	n.children = nil
	return n
}

func (n ExpAssignNode) String() string {
	return fmt.Sprintf("assign [%+v]", n.Pos())
}
