package syntree

import (
	"fmt"
)

type StmtFunctionNode struct {
	NodeBase
	name    string
	expType ExpressionType
}

func NewStmtFunctionNode(row int, col int, expType ExpressionType, name string) Node {
	n := new(StmtFunctionNode)
	n.position = Position{row, col}
	n.name = name
	n.expType = expType
	n.sibling = nil
	n.children = nil
	return n
}

func (n StmtFunctionNode) Name() string {
	return n.name
}

func (n *StmtFunctionNode) SetName(name string) {
	n.name = name
}

func (n StmtFunctionNode) ExpType() ExpressionType {
	return n.expType
}

func (n *StmtFunctionNode) SetExpType(expType ExpressionType) {
	n.expType = expType
}

func (n StmtFunctionNode) Save() bool {
	return true
}

func (n StmtFunctionNode) AddScope() bool {
	return true
}

func (n StmtFunctionNode) String() string {
	return fmt.Sprintf("function %+v %+v [%+v]", n.Name(), n.ExpType(), n.Pos())
}
