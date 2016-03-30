package syntree

import (
	"fmt"
)

type ExpCallNode struct {
	NodeBase
	name string
}

func NewExpCallNode() Node {
	n := new(ExpCallNode)
	n.position = Position{-1, -1}
	n.name = ""
	n.sibling = nil
	n.children = nil
	return n
}

func (n ExpCallNode) Name() string {
	return n.name
}

func (n *ExpCallNode) SetName(name string) {
	n.name = name
}

func (n ExpCallNode) String() string {
	return fmt.Sprintf("call %+v [%+v]", n.Name(), n.Pos())
}
