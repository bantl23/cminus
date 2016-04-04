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
	n.symbolKey = ""
	return n
}

func (n StmtFunctionNode) Name() string {
	return n.name
}

func (n StmtFunctionNode) ExpType() ExpressionType {
	return n.expType
}

func (n StmtFunctionNode) IsFunc() bool {
	return true
}

func (n StmtFunctionNode) IsDecl() bool {
	return true
}

func (n StmtFunctionNode) String() string {
	return fmt.Sprintf("function %+v %+v [%+v]", n.Name(), n.ExpType(), n.Pos())
}
