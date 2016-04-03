package syntree

import (
	"fmt"
)

type ExpParamNode struct {
	NodeBase
	name    string
	expType ExpressionType
}

func NewExpParamNode(row int, col int, expType ExpressionType, name string) Node {
	n := new(ExpParamNode)
	n.position = Position{row, col}
	n.name = name
	n.expType = expType
	n.sibling = nil
	n.children = nil
	return n
}

func (n ExpParamNode) Name() string {
	return n.name
}

func (n ExpParamNode) ExpType() ExpressionType {
	return n.expType
}

func (n ExpParamNode) Save() bool {
	if n.expType == INT_EXP_TYPE {
		return true
	}
	return false
}

func (n ExpParamNode) IsInt() bool {
	if n.expType == INT_EXP_TYPE {
		return true
	}
	return false
}

func (n ExpParamNode) IsDecl() bool {
	return true
}

func (n ExpParamNode) IsParam() bool {
	return true
}

func (n ExpParamNode) String() string {
	return fmt.Sprintf("param %+v %+v [%+v]", n.Name(), n.ExpType(), n.Pos())
}
