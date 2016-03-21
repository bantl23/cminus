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
	UNK_STATEMENT_KIND StatementKind = iota
	COMPOUND_KIND
	FUNCTION_KIND
	ITERATION_KIND
	RETURN_KIND
	SELECTION_KIND
)

var statementKinds = [...]string{
	"UnknownStmtKind",
	"CompondStmtKind",
	"FunctionStmtKind",
	"IterationStmtKind",
	"ReturnStmtKind",
	"SelectionStmtKind",
}

func (statementKind StatementKind) String() string {
	return statementKinds[statementKind]
}

type ExpressionKind int

const (
	UNK_EXPRESSION_KIND ExpressionKind = iota
	ASSIGN_KIND
	CALL_KIND
	CONST_KIND
	ID_KIND
	ID_ARRAY_KIND
	OP_KIND
	PARAM_KIND
	PARAM_ARRAY_KIND
	VAR_KIND
	VAR_ARRAY_KIND
)

var expressionKinds = [...]string{
	"UnknownExpKind",
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
}

func (expressionKind ExpressionKind) String() string {
	return expressionKinds[expressionKind]
}

type ExpressionType int

const (
	UNK_EXPRESSION_TYPE ExpressionType = iota
	VOID_TYPE
	INTEGER_TYPE
)

var expressionTypes = [...]string{
	"UnknownExpType",
	"VoidExpType",
	"IntExpType",
}

func (expressionType ExpressionType) String() string {
	return expressionTypes[expressionType]
}

type TokenType int

const (
	UNK_TOKEN_TYPE TokenType = iota
	ENDFILE
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

var tokenTypes = [...]string{
	"UnknownTokenType",
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
}

func (tokenType TokenType) String() string {
	return tokenTypes[tokenType]
}

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
	indent += 2
	for node != nil {
		for i := 0; i < indent; i++ {
			fmt.Print(" ")
		}
		if node.NodeKind == STATEMENT_KIND {
			switch node.StmtKind {
			case SELECTION_KIND:
				fmt.Println("selection")
			case ITERATION_KIND:
				fmt.Println("iteration")
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
			case ASSIGN_KIND:
				fmt.Println("assign", node.Name)
			case CALL_KIND:
				fmt.Println("call", node.Name)
			case CONST_KIND:
				fmt.Println("constant", node.Value)
			case ID_KIND:
				fmt.Println("id", node.Name)
			case ID_ARRAY_KIND:
				fmt.Println("id_array", node.Name)
			case OP_KIND:
				fmt.Println("operator", node.TokenType)
			case PARAM_KIND:
				fmt.Println("param", node.Name)
			case PARAM_ARRAY_KIND:
				fmt.Println("param_array", node.Name)
			case VAR_KIND:
				fmt.Println("var", node.Name)
			case VAR_ARRAY_KIND:
				fmt.Println("var_array", node.Name)
			default:
				fmt.Println("unknown expression")
			}
		} else {
			fmt.Println("unknown node")
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
