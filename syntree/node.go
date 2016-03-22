package syntree

import (
	"fmt"
)

type ExpressionType int

const (
	VOID_TYPE ExpressionType = iota
	INTEGER_TYPE
	UNK_EXPRESSION_TYPE
)

var expressionTypes = [...]string{
	"VoidExpType",
	"IntExpType",
	"UnknownExpType",
}

func (expressionType ExpressionType) String() string {
	return expressionTypes[expressionType]
}

type TokenType int

const (
	ENDFILE TokenType = iota
	ERROR
	ELSE
	IF
	INT
	RETURN
	VOID
	WHILE
	PLUS
	MINUS
	TIMES
	OVER
	LT
	LTE
	GT
	GTE
	EQ
	NEQ
	SEMI
	COMMA
	LPAREN
	RPAREN
	LBRACKET
	RBRACKET
	LBRACE
	RBRACE
	UNK_TOKEN_TYPE
)

var tokenTypes = [...]string{
	"EndOfFile",
	"Error",
	"Else",
	"If",
	"Int",
	"Return",
	"Void",
	"While",
	"Plus",
	"Minus",
	"Times",
	"Over",
	"LessThan",
	"LessThanEqual",
	"GreaterThan",
	"GreaterThanEqual",
	"Equals",
	"NotEquals",
	"Semicolon",
	"Comma",
	"LeftParen",
	"RightParen",
	"LeftBracket",
	"RightBracket",
	"LeftBrace",
	"RightBrace",
	"UnknownTokenType",
}

func (tokenType TokenType) String() string {
	return tokenTypes[tokenType]
}

type Node interface {
	Sibling() Node
	SetSibling(Node)
	Children() []Node
	AddChild(Node)
}

type Location interface {
	Pos() (int, int)
	SetPos(int, int)
}

type Name interface {
	Name() string
	SetName(string)
}

type Value interface {
	Value() int
	SetValue(int)
}

type ExpType interface {
	ExpType() ExpressionType
	SetExpType(ExpressionType)
}

type TokType interface {
	TokType() TokenType
	SetTokType(TokenType)
}

type NodeBase struct {
	row      int
	col      int
	sibling  Node
	children []Node
}

func (n NodeBase) Pos() (int, int) {
	return n.row, n.col
}

func (n *NodeBase) SetPos(row int, col int) {
	n.row = row
	n.col = col
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

func Print(node Node, indent int) {
	indent += 2
	for node != nil {
		for i := 0; i < indent; i++ {
			fmt.Print(" ")
		}
		fmt.Printf("%+v\n", node)
		for _, v := range node.Children() {
			Print(v, indent+2)
		}
		node = node.Sibling()
	}
	indent -= 2
}

type Procedure func(class Node)

func Traverse(node Node, pre Procedure, post Procedure) {
	if node != nil {
		pre(node)
		for _, n := range node.Children() {
			Traverse(n, pre, post)
		}
		post(node)
		Traverse(node.Sibling(), pre, post)
	}
}

func Nothing(node Node) {
	return
}
