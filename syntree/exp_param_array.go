package syntree

import (
	"fmt"
)

type ExpParamArrayNode struct {
	NodeBase
	name    string
	expType ExpressionType
}

func NewExpParamArrayNode(row int, col int, expType ExpressionType, name string) Node {
	n := new(ExpParamArrayNode)
	n.position = Position{row, col}
	n.name = name
	n.expType = expType
	n.sibling = nil
	n.children = nil
	n.symbolKey = ""
	return n
}

func (n ExpParamArrayNode) Name() string {
	return n.name
}

func (n ExpParamArrayNode) ExpType() ExpressionType {
	return n.expType
}

func (n ExpParamArrayNode) IsExp() bool {
	return true
}

func (n ExpParamArrayNode) IsParam() bool {
	return true
}

func (n ExpParamArrayNode) IsArray() bool {
	return true
}

func (n ExpParamArrayNode) IsInt() bool {
	if n.expType == INT_EXP_TYPE {
		return true
	}
	return false
}

func (n ExpParamArrayNode) String() string {
	return fmt.Sprintf("param_array %+v %+v [%+v]", n.Name(), n.ExpType(), n.Pos())
}
