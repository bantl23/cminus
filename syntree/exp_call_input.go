package syntree

import (
	"fmt"
)

type ExpCallInputNode struct {
	NodeBase
	name string
}

func NewExpCallInputNode() Node {
	n := new(ExpCallInputNode)
	n.position = Position{-1, -1}
	n.name = "input"
	n.sibling = nil
	n.children = nil
	return n
}

func (n ExpCallInputNode) Name() string {
	return n.name
}

func (n *ExpCallInputNode) SetName(name string) {
	n.name = name
}

func (n ExpCallInputNode) String() string {
	return fmt.Sprintf("call_input %+v [%+v]", n.Name(), n.Pos())
}
