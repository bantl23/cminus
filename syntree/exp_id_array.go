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
	n.row = -1
	n.col = -1
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
	row, col := n.Pos()
	return fmt.Sprintf("id_array %+v [%+v:%+v]", n.Name(), row, col)
}
