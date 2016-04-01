package syntree

import (
	"fmt"
)

type ExpVarArrayNode struct {
	NodeBase
	name    string
	value   int
	expType ExpressionType
}

func NewExpVarArrayNode(row int, col int, expType ExpressionType, name string, value int) Node {
	n := new(ExpVarArrayNode)
	n.position = Position{row, col}
	n.name = name
	n.value = value
	n.expType = expType
	n.sibling = nil
	n.children = nil
	return n
}

func (n ExpVarArrayNode) Name() string {
	return n.name
}

func (n ExpVarArrayNode) Value() int {
	return n.value
}

func (n ExpVarArrayNode) ExpType() ExpressionType {
	return n.expType
}

func (n ExpVarArrayNode) Save() bool {
	return true
}

func (n ExpVarArrayNode) String() string {
	return fmt.Sprintf("var_array %+v %+v %+v [%+v]", n.Name(), n.Value(), n.ExpType(), n.Pos())
}
