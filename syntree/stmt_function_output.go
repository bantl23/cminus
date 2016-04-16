package syntree

import (
	"fmt"
)

type StmtFunctionOutputNode struct {
	NodeBase
	name    string
	expType ExpressionType
}

func NewStmtFunctionOutputNode() Node {
	n := new(StmtFunctionNode)
	n.position = Position{-1, -1}
	n.name = "output"
	n.expType = VOID_EXP_TYPE
	n.sibling = nil
	n.children = nil
	n.symbolKey = ""
	return n
}

func (n StmtFunctionOutputNode) Name() string {
	return n.name
}

func (n StmtFunctionOutputNode) ExpType() ExpressionType {
	return n.expType
}

func (n StmtFunctionOutputNode) IsStmt() bool {
	return true
}

func (n StmtFunctionOutputNode) IsFunc() bool {
	return true
}

func (n StmtFunctionOutputNode) String() string {
	return fmt.Sprintf("function_output %+v %+v [%+v]", n.Name(), n.ExpType(), n.Pos())
}
