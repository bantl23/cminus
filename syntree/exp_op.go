package syntree

import (
	"fmt"
)

type ExpOpNode struct {
	NodeBase
	tokType TokenType
}

func NewExpOpNode() Node {
	n := new(ExpOpNode)
	n.position = Position{-1, -1}
	n.tokType = UNK_TOKEN_TYPE
	n.sibling = nil
	n.children = nil
	return n
}

func (n ExpOpNode) TokType() TokenType {
	return n.tokType
}

func (n *ExpOpNode) SetTokType(tokType TokenType) {
	n.tokType = tokType
}

func (n ExpOpNode) String() string {
	return fmt.Sprintf("op %+v [%+v]", n.TokType(), n.Pos())
}
