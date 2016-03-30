package syntree

import (
	"fmt"
)

type ExpCallOutputNode struct {
	NodeBase
	name string
}

func NewExpCallOutputNode() Node {
	n := new(ExpCallOutputNode)
	n.position = Position{-1, -1}
	n.name = "output"
	n.sibling = nil
	n.children = nil
	return n
}

func (n ExpCallOutputNode) Name() string {
	return n.name
}

func (n *ExpCallOutputNode) SetName(name string) {
	n.name = name
}

func (n ExpCallOutputNode) String() string {
	return fmt.Sprintf("call_output %+v [%+v]", n.Name(), n.Pos())
}
