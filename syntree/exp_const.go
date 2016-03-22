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
	n.row = -1
	n.col = -1
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
	row, col := n.Pos()
	return fmt.Sprintf("constant %+v [%+v:%+v]\n", n.Value(), row, col)
}
