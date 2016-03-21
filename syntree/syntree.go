package syntree

import (
	"fmt"
)

type NodeKind int

const (
	UNK_NODE_KIND NodeKind = iota
	EXPRESSION_KIND
	STATEMENT_KIND
)

var nodeKinds = [...]string{
	"UnknownNodeKind",
	"ExpressionNodeKind",
	"StatementNodeKind",
}

func (nodeKind NodeKind) String() string {
	return nodeKinds[nodeKind]
}

type StatementKind int

const (
	COMPOUND_KIND StatementKind = iota
	FUNCTION_KIND
	ITERATION_KIND
	RETURN_KIND
	SELECTION_KIND
	UNK_STATEMENT_KIND
)

var statementKinds = [...]string{
	"CompondStmtKind",
	"FunctionStmtKind",
	"IterationStmtKind",
	"ReturnStmtKind",
	"SelectionStmtKind",
	"UnknownStmtKind",
}

func (statementKind StatementKind) String() string {
	return statementKinds[statementKind]
}

type ExpressionKind int

const (
	ASSIGN_KIND ExpressionKind = iota
	CALL_KIND
	CONST_KIND
	ID_KIND
	ID_ARRAY_KIND
	OP_KIND
	PARAM_KIND
	PARAM_ARRAY_KIND
	VAR_KIND
	VAR_ARRAY_KIND
	UNK_EXPRESSION_KIND
)

var expressionKinds = [...]string{
	"AssignExpKind",
	"CallExpKind",
	"ConstExpKind",
	"IdExpKind",
	"IdArrayExpKind",
	"OpExpKind",
	"ParamExpKind",
	"ParamArrayExpKind",
	"VarExpKind",
	"VarArrayExpKind",
	"UnknownExpKind",
}

func (expressionKind ExpressionKind) String() string {
	return expressionKinds[expressionKind]
}

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

type Node struct {
	NodeKind  NodeKind
	StmtKind  StatementKind
	ExpKind   ExpressionKind
	ExpType   ExpressionType
	TokenType TokenType
	Name      string
	Value     int
	Row       int
	Col       int
	Sibling   *Node
	Children  []*Node
}

func NewNode() *Node {
	n := new(Node)
	n.NodeKind = UNK_NODE_KIND
	n.StmtKind = UNK_STATEMENT_KIND
	n.ExpKind = UNK_EXPRESSION_KIND
	n.TokenType = UNK_TOKEN_TYPE
	n.Name = ""
	n.Value = -1
	n.Row = -1
	n.Col = -1
	n.Sibling = nil
	n.Children = nil
	return n
}

func Print(node *Node, indent int) {
	indent += 2
	for node != nil {
		for i := 0; i < indent; i++ {
			fmt.Print(" ")
		}
		if node.NodeKind == STATEMENT_KIND {
			switch node.StmtKind {
			case SELECTION_KIND:
				fmt.Printf("selection [%+v%+v]\n", node.Row, node.Col)
			case ITERATION_KIND:
				fmt.Printf("iteration [%+v:%+v]\n", node.Row, node.Col)
			case COMPOUND_KIND:
				fmt.Printf("compound [%+v:%+v]\n", node.Row, node.Col)
			case FUNCTION_KIND:
				fmt.Printf("function %+v %+v [%+v:%+v]\n", node.Name, node.ExpType, node.Row, node.Col)
			case RETURN_KIND:
				fmt.Printf("return [%+v:%+v]\n", node.Row, node.Col)
			default:
				fmt.Printf("unknown statement [%+v:%+v]\n", node.Row, node.Col)
			}
		} else if node.NodeKind == EXPRESSION_KIND {
			switch node.ExpKind {
			case ASSIGN_KIND:
				fmt.Printf("assign [%+v:%+v]\n", node.Row, node.Col)
			case CALL_KIND:
				fmt.Printf("call %+v [%+v:%+v]\n", node.Name, node.Row, node.Col)
			case CONST_KIND:
				fmt.Printf("constant %+v [%+v:%+v]\n", node.Value, node.Row, node.Col)
			case ID_KIND:
				fmt.Printf("id %+v [%+v:%+v]\n", node.Name, node.Row, node.Col)
			case ID_ARRAY_KIND:
				fmt.Printf("id_array %+v [%+v:%+v]\n", node.Name, node.Row, node.Col)
			case OP_KIND:
				fmt.Printf("operator %+v [%+v:%+v]\n", node.TokenType, node.Row, node.Col)
			case PARAM_KIND:
				fmt.Printf("param %+v %+v [%+v:%+v]\n", node.Name, node.ExpType, node.Row, node.Col)
			case PARAM_ARRAY_KIND:
				fmt.Printf("param_array %+v %+v [%+v:%+v]\n", node.Name, node.ExpType, node.Row, node.Col)
			case VAR_KIND:
				fmt.Printf("var %+v %+v [%+v:%+v]\n", node.Name, node.ExpType, node.Row, node.Col)
			case VAR_ARRAY_KIND:
				fmt.Printf("var_array %+v %+v %+v [%+v:%+v]\n", node.Name, node.Value, node.ExpType, node.Row, node.Col)
			default:
				fmt.Printf("unknown expression [%+v:%+v]\n", node.Row, node.Col)
			}
		} else {
			fmt.Printf("unknown node [%+v:%+v]\n", node.Row, node.Col)
		}
		for _, v := range node.Children {
			Print(v, indent+2)
		}
		node = node.Sibling
	}
	indent -= 2
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
