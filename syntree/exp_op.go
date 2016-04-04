package syntree

import (
	"fmt"
)

type ExpOpNode struct {
	NodeBase
	tokType TokenType
}

func NewExpOpNode(row int, col int, tokType TokenType) Node {
	n := new(ExpOpNode)
	n.position = Position{row, col}
	n.tokType = tokType
	n.sibling = nil
	n.children = nil
	n.symbolKey = ""
	return n
}

func (n ExpOpNode) TokType() TokenType {
	return n.tokType
}

func (n ExpOpNode) String() string {
	return fmt.Sprintf("op %+v [%+v]", n.TokType(), n.Pos())
}
