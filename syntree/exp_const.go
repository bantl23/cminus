package syntree

import (
	"fmt"
)

type ExpConstNode struct {
	NodeBase
	value int
}

func NewExpConstNode() Node {
	n := new(ExpConstNode)
	n.position = Position{-1, -1}
	n.value = -1
	n.sibling = nil
	n.children = nil
	return n
}

func (n ExpConstNode) Value() int {
	return n.value
}

func (n *ExpConstNode) SetValue(value int) {
	n.value = value
}

func (n ExpConstNode) String() string {
	return fmt.Sprintf("constant %+v [%+v]", n.Value(), n.Pos())
}
