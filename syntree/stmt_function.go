package syntree

import (
	"fmt"
)

type StmtFunctionNode struct {
	NodeBase
	name    string
	expType ExpressionType
}

func NewStmtFunctionNode() Node {
	n := new(StmtFunctionNode)
	n.row = -1
	n.col = -1
	n.name = ""
	n.expType = UNK_EXPRESSION_TYPE
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

func (n StmtFunctionNode) String() string {
	row, col := n.Pos()
	return fmt.Sprintf("function %+v %+v [%+v:%+v]\n", n.Name(), n.ExpType(), row, col)
}
