package syntree

import (
	"fmt"
)

type ExpParamNode struct {
	NodeBase
	name    string
	expType ExpressionType
}

func NewExpParamNode() Node {
	n := new(ExpParamNode)
	n.row = -1
	n.col = -1
	n.name = ""
	n.expType = UNK_EXPRESSION_TYPE
	n.sibling = nil
	n.children = nil
	return n
}

func (n ExpParamNode) Name() string {
	return n.name
}

func (n *ExpParamNode) SetName(name string) {
	n.name = name
}

func (n ExpParamNode) ExpType() ExpressionType {
	return n.expType
}

func (n *ExpParamNode) SetExpType(expType ExpressionType) {
	n.expType = expType
}

func (n ExpParamNode) String() string {
	row, col := n.Pos()
	return fmt.Sprintf("param %+v %+v [%+v:%+v]\n", n.Name(), n.ExpType(), row, col)
}
