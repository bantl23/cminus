package syntree

import (
	"fmt"
)

type ExpIdNode struct {
	NodeBase
	name string
}

func NewExpIdNode() Node {
	n := new(ExpIdNode)
	n.row = -1
	n.col = -1
	n.name = ""
	n.sibling = nil
	n.children = nil
	return n
}

func (n ExpIdNode) Name() string {
	return n.name
}

func (n *ExpIdNode) SetName(name string) {
	n.name = name
}

func (n ExpIdNode) String() string {
	row, col := n.Pos()
	return fmt.Sprintf("id %+v [%+v:%+v]\n", n.Name(), row, col)
}
