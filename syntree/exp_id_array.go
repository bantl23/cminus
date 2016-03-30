package syntree

import (
	"fmt"
)

type ExpIdArrayNode struct {
	NodeBase
	name string
}

func NewExpIdArrayNode() Node {
	n := new(ExpIdArrayNode)
	n.position = Position{-1, -1}
	n.name = ""
	n.sibling = nil
	n.children = nil
	return n
}

func (n ExpIdArrayNode) Name() string {
	return n.name
}

func (n *ExpIdArrayNode) SetName(name string) {
	n.name = name
}

func (n ExpIdArrayNode) String() string {
	return fmt.Sprintf("id_array %+v [%+v]", n.Name(), n.Pos())
}
