package syntree

import (
	"fmt"
)

type ExpIdNode struct {
	NodeBase
	name string
}

func NewExpIdNode(row int, col int, name string) Node {
	n := new(ExpIdNode)
	n.position = Position{row, col}
	n.name = name
	n.sibling = nil
	n.children = nil
	n.symbolKey = ""
	return n
}

func (n ExpIdNode) Name() string {
	return n.name
}

func (n ExpIdNode) IsId() bool {
	return true
}

func (n ExpIdNode) String() string {
	return fmt.Sprintf("id %+v [%+v]", n.Name(), n.Pos())
}
