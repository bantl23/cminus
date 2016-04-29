package syntree

import (
	"fmt"
)

type Node interface {
	Parent() Node
	SetParent(Node)
	Sibling() Node
	SetSibling(Node)
	Children() []Node
	AddChild(Node)
	ClearChildren()
	Pos() Position
	SetPos(int, int)
	Name() string
	Value() int
	IsExp() bool
	IsStmt() bool
	IsFunc() bool
	IsCompound() bool
	IsArray() bool
	IsInt() bool
	IsReturn() bool
	IsId() bool
	IsParam() bool
	IsCall() bool
	IsVar() bool
	IsOp() bool
	IsConst() bool
	IsAssign() bool
	IsSelection() bool
	IsIteration() bool
	IsTail() bool
	SetTail(tail bool)
	ExpType() ExpressionType
	TokType() TokenType
	SymKey() string
	SetSymKey(string)
}

type NodeBase struct {
	position  Position
	parent    Node
	sibling   Node
	children  []Node
	symbolKey string
	tail      bool
}

func (n NodeBase) Pos() Position {
	return n.position
}

func (n *NodeBase) SetPos(row int, col int) {
	n.position = Position{row, col}
}

func (n NodeBase) Parent() Node {
	return n.parent
}

func (n *NodeBase) SetParent(parent Node) {
	n.parent = parent
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

func (n *NodeBase) ClearChildren() {
	n.children = nil
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

func (n NodeBase) IsExp() bool {
	return false
}

func (n NodeBase) IsStmt() bool {
	return false
}

func (n NodeBase) IsFunc() bool {
	return false
}

func (n NodeBase) IsCompound() bool {
	return false
}

func (n NodeBase) IsArray() bool {
	return false
}

func (n NodeBase) IsInt() bool {
	return false
}

func (n NodeBase) IsReturn() bool {
	return false
}

func (n NodeBase) IsId() bool {
	return false
}

func (n NodeBase) IsParam() bool {
	return false
}

func (n NodeBase) IsCall() bool {
	return false
}

func (n NodeBase) IsVar() bool {
	return false
}

func (n NodeBase) IsOp() bool {
	return false
}

func (n NodeBase) IsConst() bool {
	return false
}

func (n NodeBase) IsAssign() bool {
	return false
}

func (n NodeBase) IsSelection() bool {
	return false
}

func (n NodeBase) IsIteration() bool {
	return false
}

func (n NodeBase) IsTail() bool {
	return n.tail
}

func (n *NodeBase) SetTail(tail bool) {
	n.tail = tail
}

func (n NodeBase) SymKey() string {
	return n.symbolKey
}

func (n *NodeBase) SetSymKey(s string) {
	n.symbolKey = s
}

func PrintNode(node Node, indent int) {
	PrintNodePre(node, indent)
}

func PrintNodePre(node Node, indent int) {
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

func PrintNodeIn(node Node, indent int) {
	indent += 4
	for node != nil {
		if node.Children() != nil && len(node.Children()) >= 1 {
			PrintNode(node.Children()[0], indent)
		}
		for i := 0; i < indent; i++ {
			fmt.Print(" ")
		}
		fmt.Printf("%+v\n", node)
		if node.Children() != nil && len(node.Children()) >= 2 {
			PrintNode(node.Children()[1], indent)
		}
		node = node.Sibling()
	}
	indent -= 4
}

func PrintNodePost(node Node, indent int) {
	indent += 4
	for node != nil {
		for _, v := range node.Children() {
			PrintNode(v, indent)
		}
		for i := 0; i < indent; i++ {
			fmt.Print(" ")
		}
		fmt.Printf("%+v\n", node)
		node = node.Sibling()
	}
	indent -= 4
}

func PrintNodeWithSymKey(node Node, indent int) {
	indent += 4
	for node != nil {
		for i := 0; i < indent; i++ {
			fmt.Print(" ")
		}
		fmt.Printf("%+v (%+v)\n", node, node.SymKey())
		for _, v := range node.Children() {
			PrintNodeWithSymKey(v, indent)
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
