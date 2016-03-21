%{
package main

import (
 "github.com/bantl23/cminus/log"
 "github.com/bantl23/cminus/syntree"
 "strconv"
)

var root *syntree.Node
var savedName string
%}

%union {
  node *syntree.Node
  str string
}

%token <str> IF ELSE INT RETURN VOID WHILE PLUS MINUS TIMES
%token <str> OVER LT LTE GT GTE EQ NEQ ASSIGN SEMI COMMA
%token <str> LPAREN RPAREN LBRACE RBRACE LBRACKET RBRACKET
%token <str> ID NUM

%type <node> program declaration_list declaration var_declaration
%type <node> type_specifier fun_declaration params param_list param
%type <node> compound_stmt local_declarations statement_list statement
%type <node> expression_stmt selection_stmt iteration_stmt return_stmt
%type <node> expression var simple_expression relop additive_expression addop
%type <node> term mulop factor call args args_list empty

%nonassoc THEN
%nonassoc ELSE

%%

program             : declaration_list              {
                                                      log.ParseLog.Printf("program0: %+v %v\n", $1, yylex)
                                                      yylex.(*Lexer).Text()
                                                      root = $1
                                                    }
                    ;

declaration_list    : declaration_list declaration  {
                                                      log.ParseLog.Printf("declaration_list0: %+v %+v %v\n", $1, $2, yylex)
                                                      t := $1
                                                      if t != nil {
                                                        for t.Sibling != nil {
                                                          t = t.Sibling
                                                        }
                                                        t.Sibling = $2
                                                        $$ = $1
                                                      } else {
                                                        $$ = $2
                                                      }
                                                    }
                    | declaration                   {
                                                      log.ParseLog.Printf("declaration_list1: %+v %v\n", $1, yylex)
                                                      $$ = $1
                                                    }
                    ;

declaration         : var_declaration               {
                                                      log.ParseLog.Printf("declaration0: %+v %v\n", $1, yylex)
                                                      $$ = $1
                                                    }
                    | fun_declaration               {
                                                      log.ParseLog.Printf("declaration1: %+v %v\n", $1, yylex)
                                                      $$ = $1
                                                    }
                    ;

var_declaration     : type_specifier ID SEMI        {
                                                      log.ParseLog.Printf("var_declaration0: %+v %+v %+v %v\n", $1, $2, $3, yylex)
                                                      $$ = $1
                                                      $$.NodeKind = syntree.EXPRESSION_KIND
                                                      $$.ExpKind = syntree.VAR_KIND
                                                      $$.Name = "TODO"
                                                    }
                    | type_specifier ID LBRACKET NUM RBRACKET SEMI
                                                    {
                                                      log.ParseLog.Printf("var_declaration1: %+v %+v %+v %+v %+v %+v %+v\n", $1, $2, $3, $4, $5, $6, yylex)
                                                      $$ = $1
                                                      $$.NodeKind = syntree.EXPRESSION_KIND
                                                      $$.ExpKind = syntree.VAR_ARRAY_KIND
                                                      $$.Name = "TODO"
                                                    }
                    ;

type_specifier      : INT                           {
                                                      log.ParseLog.Printf("type_specifier0: %+v %+v\n", $1, yylex)
                                                      $$ = syntree.NewNode()
                                                      $$.ExpType = syntree.INTEGER_TYPE
                                                    }
                    | VOID                          {
                                                      log.ParseLog.Printf("type_specifier1: %+v %+v\n", $1, yylex)
                                                      $$ = syntree.NewNode()
                                                      $$.ExpType = syntree.VOID_TYPE
                                                    }
                    ;

fun_declaration     : type_specifier ID LPAREN params RPAREN compound_stmt
                                                    {
                                                      log.ParseLog.Printf("fun_declaration0: %+v %+v %+v %+v %+v %+v %+v\n", $1, $2, $3, $4, $5, $6, yylex)
                                                      $$ = $1
                                                      $$.NodeKind = syntree.STATEMENT_KIND
                                                      $$.StmtKind = syntree.FUNCTION_KIND
                                                      $$.Name = "TODO"
                                                      $$.Children = append($$.Children, $4)
                                                      $$.Children = append($$.Children, $6)
                                                    }
                    ;

params              : param_list                    {
                                                      log.ParseLog.Printf("params0: %+v %+v\n", $1, yylex)
                                                      $$ = $1
																										}
                    | VOID                          {
                                                      log.ParseLog.Printf("params1: %+v %+v\n", $1, yylex)
                                                      $$ = syntree.NewNode()
                                                      $$.ExpType = syntree.VOID_TYPE
																										}
                    ;

param_list          : param_list COMMA param        {
                                                      log.ParseLog.Printf("param_list0: %+v %+v %+v %+v\n", $1, $2, $3, yylex)
                                                      t := $1
                                                      if t != nil {
                                                        for t.Sibling != nil {
                                                          t = t.Sibling
                                                        }
                                                        t.Sibling = $3
                                                        $$ = $1
                                                      } else {
                                                        $$ = $3
                                                      }
																										}
                    | param                         {
                                                      log.ParseLog.Printf("param_list1: %+v %+v\n", $1, yylex)
                                                      $$ = $1
																										}
                    ;

param               : type_specifier ID             {
                                                      log.ParseLog.Printf("param0: %+v %+v %+v\n", $1, $2, yylex)
                                                      $$ = $1
                                                      $$.NodeKind = syntree.EXPRESSION_KIND
                                                      $$.ExpKind = syntree.PARAM_KIND
                                                      $$.Name = "TODO"
																										}
                    | type_specifier ID LBRACKET RBRACKET
                                                    {
                                                      log.ParseLog.Printf("param1: %+v %+v %+v %+v %+v\n", $1, $2, $3, $3, yylex)
                                                      $$ = $1
                                                      $$.NodeKind = syntree.EXPRESSION_KIND
                                                      $$.ExpKind = syntree.PARAM_ARRAY_KIND
                                                      $$.Name = "TODO"
																										}
                    ;

compound_stmt       : LBRACE local_declarations statement_list RBRACE
                                                    {
                                                      log.ParseLog.Printf("compound_stmt0: %+v %+v %+v %+v %+v\n", $1, $2, $3, $4, yylex)
                                                      $$ = syntree.NewNode()
                                                      $$.NodeKind = syntree.STATEMENT_KIND
                                                      $$.StmtKind = syntree.COMPOUND_KIND
                                                      $$.Children = append($$.Children, $2)
                                                      $$.Children = append($$.Children, $3)
																										}
                    ;

local_declarations  : local_declarations var_declaration
                                                    {
                                                      log.ParseLog.Printf("local_declarations0: %+v %+v %+v\n", $1, $2, yylex)
                                                      t := $1
                                                      if t != nil {
                                                        for t.Sibling != nil {
                                                          t = t.Sibling
                                                        }
                                                        t.Sibling = $2
                                                        $$ = $1
                                                      } else {
                                                        $$ = $2
                                                      }
																										}
                    | empty                         {
                                                      log.ParseLog.Printf("local_declarations1: %+v %+v\n", $1, yylex)
                                                      $$ = $1
																										}
                    ;

statement_list      : statement_list statement      {
                                                      log.ParseLog.Printf("statement_list0: %+v %+v %+v\n", $1, $2, yylex)
                                                      t := $1
                                                      if t != nil {
                                                        for t.Sibling != nil {
                                                          t = t.Sibling
                                                        }
                                                        t.Sibling = $2
                                                        $$ = $1
                                                      } else {
                                                        $$ = $2
                                                      }
																										}
                    | empty                         {
                                                      log.ParseLog.Printf("statement_list1: %+v %+v\n", $1, yylex)
                                                      $$ = $1
                                                    }
                    ;

statement           : expression_stmt               {
                                                      log.ParseLog.Printf("statement0: %+v %+v\n", $1, yylex)
                                                      $$ = $1
																										}
                    | compound_stmt                 {
                                                      log.ParseLog.Printf("statement1: %+v %+v\n", $1, yylex)
                                                      $$ = $1
																										}
                    | selection_stmt                {
                                                      log.ParseLog.Printf("statement2: %+v %+v\n", $1, yylex)
                                                      $$ = $1
																										}
                    | iteration_stmt                {
                                                      log.ParseLog.Printf("statement3: %+v %+v\n", $1, yylex)
                                                      $$ = $1
																										}
                    | return_stmt                   {
                                                      log.ParseLog.Printf("statement4: %+v %+v\n", $1, yylex)
                                                      $$ = $1
																										}
                    ;

expression_stmt     : expression SEMI               {
                                                      log.ParseLog.Printf("expression_stmt0: %+v %+v %+v\n", $1, $2, yylex)
                                                      $$ = $1
																										}
                    | SEMI                          {
                                                      log.ParseLog.Printf("expression_stmt1: %+v %+v\n", $1, yylex)
                                                      $$ = nil
																										}
                    ;

selection_stmt      : IF LPAREN expression RPAREN statement %prec THEN
                                                    {
                                                      log.ParseLog.Printf("selection_stmt0: %+v %+v %+v %+v %+v %+v\n", $1, $2, $3, $4, $5, yylex)
                                                      $$ = syntree.NewNode()
                                                      $$.NodeKind = syntree.STATEMENT_KIND
                                                      $$.StmtKind = syntree.SELECTION_KIND
                                                      $$.Children = append($$.Children, $3)
                                                      $$.Children = append($$.Children, $5)
																										}
                    | IF LPAREN expression RPAREN statement ELSE statement
                                                    {
                                                      log.ParseLog.Printf("selection_stmt1: %+v %+v %+v %+v %+v %+v %+v %+v\n", $1, $2, $3, $4, $5, $6, $7, yylex)
                                                      $$ = syntree.NewNode()
                                                      $$.NodeKind = syntree.STATEMENT_KIND
                                                      $$.StmtKind = syntree.SELECTION_KIND
                                                      $$.Children = append($$.Children, $3)
                                                      $$.Children = append($$.Children, $5)
                                                      $$.Children = append($$.Children, $7)
																										}
                    ;

iteration_stmt      : WHILE LPAREN expression RPAREN statement
                                                    {
                                                      log.ParseLog.Printf("iteration_stmt0: %+v %+v %+v %+v %+v %+v\n", $1, $2, $3, $4, $5, yylex)
                                                      $$ = syntree.NewNode()
                                                      $$.NodeKind = syntree.STATEMENT_KIND
                                                      $$.StmtKind = syntree.ITERATION_KIND
                                                      $$.Children = append($$.Children, $3)
                                                      $$.Children = append($$.Children, $5)
																										}
                    ;

return_stmt         : RETURN SEMI                   {
                                                      log.ParseLog.Printf("return_stmt0: %+v %+v %+v\n", $1, $2, yylex)
                                                      $$ = syntree.NewNode()
                                                      $$.NodeKind = syntree.STATEMENT_KIND
                                                      $$.StmtKind = syntree.RETURN_KIND
																										}
                    | RETURN expression SEMI        {
                                                      log.ParseLog.Printf("return_stmt1: %+v %+v %+v %+v\n", $1, $2, $3, yylex)
                                                      $$ = syntree.NewNode()
                                                      $$.NodeKind = syntree.STATEMENT_KIND
                                                      $$.StmtKind = syntree.RETURN_KIND
                                                      $$.Children = append($$.Children, $2)
                                                    }
                    ;

expression          : var ASSIGN expression         {
                                                      log.ParseLog.Printf("expression0: %+v %+v %+v %+v\n", $1, $2, $3, yylex)
                                                      $$ = syntree.NewNode()
                                                      $$.NodeKind = syntree.EXPRESSION_KIND
                                                      $$.ExpKind = syntree.ASSIGN_KIND
                                                      $$.Children = append($$.Children, $1)
                                                      $$.Children = append($$.Children, $3)
                                                    }
                    | simple_expression             {
                                                      log.ParseLog.Printf("expression1: %+v %+v\n", $1, yylex)
                                                      $$ = $1
                                                    }
                    ;

var                 : ID                            {
                                                      log.ParseLog.Printf("var0: %+v %+v\n", $1, yylex)
                                                      $$ = syntree.NewNode()
                                                      $$.NodeKind = syntree.EXPRESSION_KIND
                                                      $$.ExpKind = syntree.ID_KIND
                                                      $$.Name = "TODO"
                                                    }
                    | ID LBRACKET expression RBRACKET
                                                    {
                                                      log.ParseLog.Printf("var1: %+v %+v %+v %+v %+v\n", $1, $2, $3, $4, yylex)
                                                      $$ = syntree.NewNode()
                                                      $$.NodeKind = syntree.EXPRESSION_KIND
                                                      $$.ExpKind = syntree.ID_ARRAY_KIND
                                                      $$.Name = "TODO"
                                                      $$.Children = append($$.Children, $3)
                                                    }
                    ;

simple_expression   : additive_expression relop additive_expression
                                                    {
                                                      log.ParseLog.Printf("simple_expression0: %+v %+v %+v %+v\n", $1, $2, $3, yylex)
                                                      $$ = $2
                                                      $$.Children = append($$.Children, $1)
                                                      $$.Children = append($$.Children, $3)
                                                    }
                    | additive_expression           {
                                                      log.ParseLog.Printf("simple_expression6: %+v %+v\n", $1, yylex)
                                                      $$ = $1
                                                    }
                    ;

relop               : LT                            {
                                                      log.ParseLog.Printf("relop0: %+v %+v\n", $1, yylex)
                                                      $$ = syntree.NewNode()
                                                      $$.NodeKind = syntree.EXPRESSION_KIND
                                                      $$.ExpKind = syntree.OP_KIND
                                                      $$.TokenType = syntree.LT
                                                    }
                    | LTE                           {
                                                      log.ParseLog.Printf("relop1: %+v %+v\n", $1, yylex)
                                                      $$ = syntree.NewNode()
                                                      $$.NodeKind = syntree.EXPRESSION_KIND
                                                      $$.ExpKind = syntree.OP_KIND
                                                      $$.TokenType = syntree.LTE
                                                    }
                    | GT                            {
                                                      log.ParseLog.Printf("relop2: %+v %+v\n", $1, yylex)
                                                      $$ = syntree.NewNode()
                                                      $$.NodeKind = syntree.EXPRESSION_KIND
                                                      $$.ExpKind = syntree.OP_KIND
                                                      $$.TokenType = syntree.GT
                                                    }
                    | GTE                           {
                                                      log.ParseLog.Printf("relop3: %+v %+v\n", $1, yylex)
                                                      $$ = syntree.NewNode()
                                                      $$.NodeKind = syntree.EXPRESSION_KIND
                                                      $$.ExpKind = syntree.OP_KIND
                                                      $$.TokenType = syntree.GTE
                                                    }
                    | EQ                            {
                                                      log.ParseLog.Printf("relop4: %+v %+v\n", $1, yylex)
                                                      $$ = syntree.NewNode()
                                                      $$.NodeKind = syntree.EXPRESSION_KIND
                                                      $$.ExpKind = syntree.OP_KIND
                                                      $$.TokenType = syntree.EQ
                                                    }
                    | NEQ                           {
                                                      log.ParseLog.Printf("relop5: %+v %+v\n", $1, yylex)
                                                      $$ = syntree.NewNode()
                                                      $$.NodeKind = syntree.EXPRESSION_KIND
                                                      $$.ExpKind = syntree.OP_KIND
                                                      $$.TokenType = syntree.NEQ
                                                    }

additive_expression : additive_expression addop term
                                                    {
																											log.ParseLog.Printf("additive_expression0: %+v %+v %+v %+v\n", $1, $2, $3, yylex)
                                                      $$ = $2
                                                      $$.Children = append($$.Children, $1)
                                                      $$.Children = append($$.Children, $3)
																										}
                    | term                          {
																											log.ParseLog.Printf("additive_expression1: %+v %+v\n", $1, yylex)
                                                      $$ = $1
																										}
                    ;

addop               : PLUS                          {
                                                      log.ParseLog.Printf("addop0: %+v %+v\n", $1, yylex)
                                                      $$ = syntree.NewNode()
                                                      $$.NodeKind = syntree.EXPRESSION_KIND
                                                      $$.ExpKind = syntree.OP_KIND
                                                      $$.TokenType = syntree.PLUS
                                                    }
                    | MINUS                         {
                                                      log.ParseLog.Printf("addop1: %+v %+v\n", $1, yylex)
                                                      $$ = syntree.NewNode()
                                                      $$.NodeKind = syntree.EXPRESSION_KIND
                                                      $$.ExpKind = syntree.OP_KIND
                                                      $$.TokenType = syntree.MINUS
                                                    }

term                : term mulop factor             {
																											log.ParseLog.Printf("term0: %+v %+v %+v %+v\n", $1, $2, $3, yylex)
                                                      $$ = $2
                                                      $$.Children = append($$.Children, $1)
                                                      $$.Children = append($$.Children, $3)
																										}
                    | factor                        {
																											log.ParseLog.Printf("term1: %+v %+v\n", $1, yylex)
                                                      $$ = $1
																										}
                    ;

mulop               : TIMES                         {
                                                      log.ParseLog.Printf("mulop0: %+v %+v\n", $1, yylex)
                                                      $$ = syntree.NewNode()
                                                      $$.NodeKind = syntree.EXPRESSION_KIND
                                                      $$.ExpKind = syntree.OP_KIND
                                                      $$.TokenType = syntree.TIMES
                                                    }
                    | OVER                          {
                                                      log.ParseLog.Printf("mulop1: %+v %+v\n", $1, yylex)
                                                      $$ = syntree.NewNode()
                                                      $$.NodeKind = syntree.EXPRESSION_KIND
                                                      $$.ExpKind = syntree.OP_KIND
                                                      $$.TokenType = syntree.OVER
                                                    }

factor              : LPAREN expression RPAREN      {
																											log.ParseLog.Printf("factor0: %+v %+v %+v %+v\n", $1, $2, $3, yylex)
                                                      $$ = $2
																										}
                    | var                           {
																											log.ParseLog.Printf("factor1: %+v %+v\n", $1, yylex)
                                                      $$ = $1
																										}
                    | call                          {
																											log.ParseLog.Printf("factor2: %+v %+v\n", $1, yylex)
                                                      $$ = $1
																										}
                    | NUM                           {
																											log.ParseLog.Printf("factor3: %+v %+v\n", $1, yylex)
                                                      $$ = syntree.NewNode()
                                                      $$.NodeKind = syntree.EXPRESSION_KIND
                                                      $$.ExpKind = syntree.CONST_KIND
                                                      $$.Value, _ = strconv.Atoi(yylex.(*Lexer).Text())
																										}
                    ;

call                : ID LPAREN args RPAREN         {
																											log.ParseLog.Printf("call0: %+v %+v %+v %+v %+v\n", $1, $2, $3, $4, yylex)
                                                      $$ = syntree.NewNode()
                                                      $$.NodeKind = syntree.EXPRESSION_KIND
                                                      $$.ExpKind = syntree.CALL_KIND
                                                      $$.Name = "TODO"
                                                      $$.Children = append($$.Children, $3)
																										}
                    ;

args                : args_list                     {
																											log.ParseLog.Printf("args0: %+v %+v\n", $1, yylex)
                                                      $$ = $1
																										}
                    | empty                         {
																											log.ParseLog.Printf("args1: %+v %+v\n", $1, yylex)
                                                      $$ = $1
																										}
                    ;

args_list           : args_list COMMA expression    {
                                                      log.ParseLog.Printf("args_list0: %+v %+v %+v %+v\n", $1, $2, $3, yylex)
                                                      t := $1
                                                      if t != nil {
                                                        for t.Sibling != nil {
                                                          t = t.Sibling
                                                        }
                                                        t.Sibling = $3
                                                        $$ = $1
                                                      } else {
                                                        $$ = $3
                                                      }
                                                    }
                    | expression                    {
                                                      log.ParseLog.Printf("args_list1: %+v %+v\n", $1, yylex)
                                                      $$ = $1
                                                    }
                    ;

empty               : /* empty */                   {
                                                      log.ParseLog.Printf("empty0: %+v\n", yylex)
                                                      $$ = nil
                                                    }
                    ;

%%
