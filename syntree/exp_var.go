package syntree

import (
	"fmt"
)

type ExpVarNode struct {
	NodeBase
	name    string
	expType ExpressionType
}

func NewExpVarNode() Node {
	n := new(ExpVarNode)
	n.position = Position{-1, -1}
	n.name = ""
	n.expType = UNK_EXPRESSION_TYPE
	n.sibling = nil
	n.children = nil
	return n
}

func (n ExpVarNode) Name() string {
	return n.name
}

func (n *ExpVarNode) SetName(name string) {
	n.name = name
}

func (n ExpVarNode) ExpType() ExpressionType {
	return n.expType
}

func (n *ExpVarNode) SetExpType(expType ExpressionType) {
	n.expType = expType
}

func (n ExpVarNode) String() string {
	return fmt.Sprintf("var %+v %+v [%+v]", n.Name(), n.ExpType(), n.Pos())
}
