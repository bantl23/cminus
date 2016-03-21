package syntree

import (
	"fmt"
)

type NodeKind int

const (
	STATEMENT_KIND NodeKind = iota
	EXPRESSION_KIND
)

type StatementKind int

const (
	SELECTION_KIND StatementKind = iota
	ITERATION_KIND
	COMPOUND_KIND
	FUNCTION_KIND
	RETURN_KIND
)

type ExpressionKind int

const (
	OP_KIND ExpressionKind = iota
	CONST_KIND
	ID_KIND
	ID_ARRAY_KIND
	ASSIGN_KIND
	CALL_KIND
	PARAM_KIND
	PARAM_ARRAY_KIND
	VAR_KIND
	VAR_ARRAY_KIND
)

type ExpressionType int

const (
	VOID_TYPE ExpressionType = iota
	INTEGER_TYPE
)

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
)

type Node struct {
	NodeKind   NodeKind
	StmtKind   StatementKind
	ExpKind    ExpressionKind
	ExpType    ExpressionType
	TokenType  TokenType
	Name       string
	Value      int
	LineNumber int
	Sibling    *Node
	Children   []*Node
}

func NewNode() *Node {
	n := new(Node)
	return n
}

func Print(node *Node, indent int) {
	for i := 0; i < indent; i++ {
		fmt.Print(" ")
	}
	for node != nil {
		if node.NodeKind == STATEMENT_KIND {
			switch node.StmtKind {
			case SELECTION_KIND:
				fmt.Println("if")
			case ITERATION_KIND:
				fmt.Println("while")
			case COMPOUND_KIND:
				fmt.Println("compound")
			case FUNCTION_KIND:
				fmt.Println("function")
			case RETURN_KIND:
				fmt.Println("return")
			default:
				fmt.Println("unknown statement")
			}
		} else if node.NodeKind == EXPRESSION_KIND {
			switch node.ExpKind {
			case OP_KIND:
				fmt.Println("operator=", node.TokenType)
			case CONST_KIND:
				fmt.Println("constant=", node.Value)
			case ID_KIND:
				fmt.Println("id=", node.Name)
			case ID_ARRAY_KIND:
				fmt.Println("id_array=", node.Name)
			case ASSIGN_KIND:
				fmt.Println("assign=", node.Name)
			case CALL_KIND:
				fmt.Println("call=", node.Name)
			case PARAM_KIND:
				fmt.Println("param=", node.Name)
			case PARAM_ARRAY_KIND:
				fmt.Println("param_array=", node.Name)
			case VAR_KIND:
				fmt.Println("var=", node.Name)
			case VAR_ARRAY_KIND:
				fmt.Println("var_array=", node.Name)
			default:
				fmt.Println("unknown expression kind")
			}
		} else {
			fmt.Println("unknown node kind")
		}
		for _, v := range node.Children {
			Print(v, indent+2)
		}
		node = node.Sibling
	}
}

type Procedure func(class *Node)

func Traverse(node *Node, pre Procedure, post Procedure) {
	if node != nil {
		pre(node)
		for _, n := range node.Children {
			Traverse(n, pre, post)
		}
		post(node)
		Traverse(node.Sibling, pre, post)
	}
}

func Nothing(node *Node) {
	return
}
