package syntree

import (
	"fmt"
)

type ExpVarNode struct {
	NodeBase
	name    string
	expType ExpressionType
}

func NewExpVarNode(row int, col int, expType ExpressionType, name string) Node {
	n := new(ExpVarNode)
	n.position = Position{row, col}
	n.name = name
	n.expType = expType
	n.sibling = nil
	n.children = nil
	return n
}

func (n ExpVarNode) Name() string {
	return n.name
}

func (n ExpVarNode) ExpType() ExpressionType {
	return n.expType
}

func (n ExpVarNode) Save() bool {
	if n.expType == INT_EXP_TYPE {
		return true
	}
	return false
}

func (n ExpVarNode) IsInt() bool {
	if n.expType == INT_EXP_TYPE {
		return true
	}
	return false
}

func (n ExpVarNode) IsDecl() bool {
	return true
}

func (n ExpVarNode) String() string {
	return fmt.Sprintf("var %+v %+v [%+v]", n.Name(), n.ExpType(), n.Pos())
}
