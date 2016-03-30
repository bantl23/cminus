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
	n.position = Position{-1, -1}
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
	return fmt.Sprintf("id %+v [%+v]", n.Name(), n.Pos())
}
