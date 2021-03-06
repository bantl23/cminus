package syntree

import (
	"fmt"
)

type ExpIdArrayNode struct {
	NodeBase
	name string
}

func NewExpIdArrayNode(row int, col int, name string) Node {
	n := new(ExpIdArrayNode)
	n.position = Position{row, col}
	n.name = name
	n.sibling = nil
	n.children = nil
	n.symbolKey = ""
	return n
}

func (n ExpIdArrayNode) Name() string {
	return n.name
}

func (n ExpIdArrayNode) IsExp() bool {
	return true
}

func (n ExpIdArrayNode) IsId() bool {
	return true
}

func (n ExpIdArrayNode) IsArray() bool {
	return true
}

func (n ExpIdArrayNode) String() string {
	return fmt.Sprintf("id_array %+v [%+v]", n.Name(), n.Pos())
}
