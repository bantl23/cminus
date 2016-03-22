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
	n.row = -1
	n.col = -1
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
	row, col := n.Pos()
	return fmt.Sprintf("op %+v [%+v:%+v]", n.TokType(), row, col)
}
