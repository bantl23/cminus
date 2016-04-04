package syntree

import (
	"fmt"
)

type StmtFunctionInputNode struct {
	NodeBase
	name    string
	expType ExpressionType
}

func NewStmtFunctionInputNode() Node {
	n := new(StmtFunctionInputNode)
	n.position = Position{-1, -1}
	n.name = "input"
	n.expType = INT_EXP_TYPE
	n.sibling = nil
	n.children = nil
	n.symbolKey = ""
	return n
}

func (n StmtFunctionInputNode) Name() string {
	return n.name
}

func (n StmtFunctionInputNode) ExpType() ExpressionType {
	return n.expType
}

func (n StmtFunctionInputNode) IsFunc() bool {
	return true
}

func (n StmtFunctionInputNode) IsDecl() bool {
	return true
}

func (n StmtFunctionInputNode) String() string {
	return fmt.Sprintf("function_input %+v %+v [%+v]", n.Name(), n.ExpType(), n.Pos())
}
