package syntree

import (
	"fmt"
)

type Node interface {
	Sibling() Node
	SetSibling(Node)
	Children() []Node
	AddChild(Node)
	Pos() Position
	SetPos(int, int)
	Name() string
	Value() int
	Save() bool
	AddScope() bool
	IsFunc() bool
	IsArray() bool
	IsInt() bool
	IsDecl() bool
	IsReturn() bool
	IsParam() bool
	IsCall() bool
	ExpType() ExpressionType
	TokType() TokenType
}

type NodeBase struct {
	position Position
	sibling  Node
	children []Node
}

func (n NodeBase) Pos() Position {
	return n.position
}

func (n *NodeBase) SetPos(row int, col int) {
	n.position = Position{row, col}
}

func (n NodeBase) Sibling() Node {
	return n.sibling
}

func (n *NodeBase) SetSibling(sibling Node) {
	n.sibling = sibling
}

func (n NodeBase) Children() []Node {
	return n.children
}

func (n *NodeBase) AddChild(child Node) {
	n.children = append(n.children, child)
}

func (n NodeBase) ExpType() ExpressionType {
	return UNK_EXP_TYPE
}

func (n NodeBase) TokType() TokenType {
	return UNK_TOKEN_TYPE
}

func (n NodeBase) Name() string {
	return ""
}

func (n NodeBase) Value() int {
	return 0
}

func (n NodeBase) Save() bool {
	return false
}

func (n NodeBase) AddScope() bool {
	return false
}

func (n NodeBase) IsFunc() bool {
	return false
}

func (n NodeBase) IsArray() bool {
	return false
}

func (n NodeBase) IsInt() bool {
	return false
}

func (n NodeBase) IsDecl() bool {
	return false
}

func (n NodeBase) IsReturn() bool {
	return false
}

func (n NodeBase) IsParam() bool {
	return false
}

func (n NodeBase) IsCall() bool {
	return false
}

func PrintNode(node Node, indent int) {
	indent += 4
	for node != nil {
		for i := 0; i < indent; i++ {
			fmt.Print(" ")
		}
		fmt.Printf("%+v\n", node)
		for _, v := range node.Children() {
			PrintNode(v, indent)
		}
		node = node.Sibling()
	}
	indent -= 4
}

type Procedure func(class Node)

func TraverseNode(node Node, pre Procedure, post Procedure) {
	if node != nil {
		pre(node)
		for _, n := range node.Children() {
			TraverseNode(n, pre, post)
		}
		post(node)
		TraverseNode(node.Sibling(), pre, post)
	}
}

func Nothing(node Node) {
	return
}
