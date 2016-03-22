package syntree

import (
	"fmt"
)

type ExpParamArrayNode struct {
	NodeBase
	name    string
	expType ExpressionType
}

func NewExpParamArrayNode() Node {
	n := new(ExpParamArrayNode)
	n.row = -1
	n.col = -1
	n.name = ""
	n.expType = UNK_EXPRESSION_TYPE
	n.sibling = nil
	n.children = nil
	return n
}

func (n ExpParamArrayNode) Name() string {
	return n.name
}

func (n *ExpParamArrayNode) SetName(name string) {
	n.name = name
}

func (n ExpParamArrayNode) ExpType() ExpressionType {
	return n.expType
}

func (n *ExpParamArrayNode) SetExpType(expType ExpressionType) {
	n.expType = expType
}

func (n ExpParamArrayNode) String() string {
	row, col := n.Pos()
	return fmt.Sprintf("param_array %+v %+v [%+v:%+v]\n", n.Name(), n.ExpType(), row, col)
}
