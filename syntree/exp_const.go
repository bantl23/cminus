package syntree

import (
	"fmt"
)

type ExpConstNode struct {
	NodeBase
	value int
}

func NewExpConstNode(row int, col int, value int) Node {
	n := new(ExpConstNode)
	n.position = Position{row, col}
	n.value = value
	n.sibling = nil
	n.children = nil
	return n
}

func (n ExpConstNode) Value() int {
	return n.value
}

func (n ExpConstNode) String() string {
	return fmt.Sprintf("constant %+v [%+v]", n.Value(), n.Pos())
}
