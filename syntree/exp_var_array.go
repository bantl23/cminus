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
	n.symbolKey = ""
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

func (n ExpVarArrayNode) IsExp() bool {
	return true
}

func (n ExpVarArrayNode) IsVar() bool {
	return true
}

func (n ExpVarArrayNode) IsArray() bool {
	return true
}

func (n ExpVarArrayNode) IsInt() bool {
	if n.expType == INT_EXP_TYPE {
		return true
	}
	return false
}

func (n ExpVarArrayNode) String() string {
	return fmt.Sprintf("var_array %+v %+v %+v [%+v]", n.Name(), n.Value(), n.ExpType(), n.Pos())
}
