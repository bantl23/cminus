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

func NewExpVarArrayNode() Node {
	n := new(ExpVarArrayNode)
	n.row = -1
	n.col = -1
	n.name = ""
	n.expType = UNK_EXPRESSION_TYPE
	n.sibling = nil
	n.children = nil
	return n
}

func (n ExpVarArrayNode) Name() string {
	return n.name
}

func (n *ExpVarArrayNode) SetName(name string) {
	n.name = name
}

func (n ExpVarArrayNode) Value() int {
	return n.value
}

func (n *ExpVarArrayNode) SetValue(value int) {
	n.value = value
}

func (n ExpVarArrayNode) ExpType() ExpressionType {
	return n.expType
}

func (n *ExpVarArrayNode) SetExpType(expType ExpressionType) {
	n.expType = expType
}

func (n ExpVarArrayNode) String() string {
	row, col := n.Pos()
	return fmt.Sprintf("var_array %+v %+v %+v [%+v:%+v]\n", n.Name(), n.Value(), n.ExpType(), row, col)
}
