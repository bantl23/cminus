package syntree

import (
	"fmt"
)

type ExpCallNode struct {
	NodeBase
	name string
}

func NewExpCallNode(row int, col int, name string) Node {
	n := new(ExpCallNode)
	n.position = Position{row, col}
	n.name = name
	n.sibling = nil
	n.children = nil
	return n
}

func (n ExpCallNode) Name() string {
	return n.name
}

func (n ExpCallNode) IsCall() bool {
	return true
}

func (n ExpCallNode) String() string {
	return fmt.Sprintf("call %+v [%+v]", n.Name(), n.Pos())
}
