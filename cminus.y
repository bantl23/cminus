%{
package main

import (
 "github.com/bantl23/cminus/log"
 "github.com/bantl23/cminus/syntree"
 "strconv"
)

var root syntree.Node
%}

%union {
  node syntree.Node
  exp  syntree.ExpressionType
  tok  syntree.TokenType
  str  string
}

%token <str> IF ELSE INT RETURN VOID WHILE ASSIGN SEMI COMMA
%token <str> LPAREN RPAREN LBRACE RBRACE LBRACKET RBRACKET
%token <str> ID NUM PLUS MINUS TIMES OVER LT LTE GT GTE EQ NEQ

%type <exp> type_specifier
%type <tok> addop mulop relop

%type <node> program declaration_list declaration var_declaration
%type <node> fun_declaration params param_list param
%type <node> compound_stmt local_declarations statement_list statement
%type <node> expression_stmt selection_stmt iteration_stmt return_stmt
%type <node> expression var simple_expression additive_expression
%type <node> term factor call args args_list empty

%nonassoc THEN
%nonassoc ELSE

%%

program             : declaration_list              {
                                                      log.ParseLog.Printf("program0: %+v %+v\n", $1, yylex)
                                                      root = $1
                                                    }
                    ;

declaration_list    : declaration_list declaration  {
                                                      log.ParseLog.Printf("declaration_list0: %+v %+v %+v\n", $1, $2, yylex)
                                                      t := $1
                                                      if t != nil {
                                                        for t.Sibling() != nil {
                                                          t = t.Sibling()
                                                        }
                                                        t.SetSibling($2)
                                                        $$ = $1
                                                      } else {
                                                        $$ = $2
                                                      }
                                                    }
                    | declaration                   {
                                                      log.ParseLog.Printf("declaration_list1: %+v %+v\n", $1, yylex)
                                                      $$ = $1
                                                    }
                    ;

declaration         : var_declaration               {
                                                      log.ParseLog.Printf("declaration0: %+v %+v\n", $1, yylex)
                                                      $$ = $1
                                                    }
                    | fun_declaration               {
                                                      log.ParseLog.Printf("declaration1: %+v %+v\n", $1, yylex)
                                                      $$ = $1
                                                    }
                    ;

var_declaration     : type_specifier ID SEMI        {
                                                      log.ParseLog.Printf("var_declaration0: %+v %+v %+v %+v\n", $1, $2, $3, yylex)
                                                      $$ = syntree.NewExpVarNode()
                                                      $$.(syntree.Location).SetPos(yylex.(*Lexer).Row(), yylex.(*Lexer).Col())
                                                      $$.(syntree.ExpType).SetExpType($1)
                                                      $$.(syntree.Name).SetName($2)
                                                      log.ParseLog.Printf("var_declaration0: %+v\n", $$)
                                                    }
                    | type_specifier ID LBRACKET NUM RBRACKET SEMI
                                                    {
                                                      log.ParseLog.Printf("var_declaration1: %+v %+v %+v %+v %+v %+v %+v\n", $1, $2, $3, $4, $5, $6, yylex)
                                                      $$ = syntree.NewExpVarArrayNode()
                                                      $$.(syntree.Location).SetPos(yylex.(*Lexer).Row(), yylex.(*Lexer).Col())
                                                      $$.(syntree.ExpType).SetExpType($1)
                                                      $$.(syntree.Name).SetName($2)
                                                      v, _ := strconv.Atoi($4)
                                                      $$.(syntree.Value).SetValue(v)
                                                      log.ParseLog.Printf("var_declaration1: %+v\n", $$)
                                                    }
                    ;

type_specifier      : INT                           {
                                                      log.ParseLog.Printf("type_specifier0: %+v %+v\n", $1, yylex)
                                                      $$ = syntree.INTEGER_TYPE
                                                      log.ParseLog.Printf("type_specifier0: %+v\n", $$)
                                                    }
                    | VOID                          {
                                                      log.ParseLog.Printf("type_specifier1: %+v %+v\n", $1, yylex)
                                                      $$ = syntree.VOID_TYPE
                                                      log.ParseLog.Printf("type_specifier1: %+v\n", $$)
                                                    }
                    ;

fun_declaration     : type_specifier ID LPAREN params RPAREN compound_stmt
                                                    {
                                                      log.ParseLog.Printf("fun_declaration0: %+v %+v %+v %+v %+v %+v %+v\n", $1, $2, $3, $4, $5, $6, yylex)
                                                      $$ = syntree.NewStmtFunctionNode()
                                                      $$.(syntree.Location).SetPos(yylex.(*Lexer).Row(), yylex.(*Lexer).Col())
                                                      $$.(syntree.ExpType).SetExpType($1)
                                                      $$.(syntree.Name).SetName($2)
                                                      $$.AddChild($4)
                                                      $$.AddChild($6)
                                                      log.ParseLog.Printf("fun_declaration0: %+v\n", $$)
                                                    }
                    ;

params              : param_list                    {
                                                      log.ParseLog.Printf("params0: %+v %+v\n", $1, yylex)
                                                      $$ = $1
																										}
                    | VOID                          {
                                                      log.ParseLog.Printf("params1: %+v %+v\n", $1, yylex)
                                                      $$ = syntree.NewExpParamNode()
                                                      $$.(syntree.Location).SetPos(yylex.(*Lexer).Row(), yylex.(*Lexer).Col())
                                                      $$.(syntree.ExpType).SetExpType(syntree.VOID_TYPE)
                                                      log.ParseLog.Printf("params1: %+v\n", $$)
																										}
                    ;

param_list          : param_list COMMA param        {
                                                      log.ParseLog.Printf("param_list0: %+v %+v %+v %+v\n", $1, $2, $3, yylex)
                                                      t := $1
                                                      if t != nil {
                                                        for t.Sibling() != nil {
                                                          t = t.Sibling()
                                                        }
                                                        t.SetSibling($3)
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
                                                      $$ = syntree.NewExpParamNode()
                                                      $$.(syntree.Location).SetPos(yylex.(*Lexer).Row(), yylex.(*Lexer).Col())
                                                      $$.(syntree.ExpType).SetExpType($1)
                                                      $$.(syntree.Name).SetName($2)
                                                      log.ParseLog.Printf("param0: %+v\n", $$)
																										}
                    | type_specifier ID LBRACKET RBRACKET
                                                    {
                                                      log.ParseLog.Printf("param1: %+v %+v %+v %+v %+v\n", $1, $2, $3, $3, yylex)
                                                      $$ = syntree.NewExpParamArrayNode()
                                                      $$.(syntree.Location).SetPos(yylex.(*Lexer).Row(), yylex.(*Lexer).Col())
                                                      $$.(syntree.ExpType).SetExpType($1)
                                                      $$.(syntree.Name).SetName($2)
                                                      log.ParseLog.Printf("param1: %+v\n", $$)
																										}
                    ;

compound_stmt       : LBRACE local_declarations statement_list RBRACE
                                                    {
                                                      log.ParseLog.Printf("compound_stmt0: %+v %+v %+v %+v %+v\n", $1, $2, $3, $4, yylex)
                                                      $$ = syntree.NewStmtCompoundNode()
                                                      $$.(syntree.Location).SetPos(yylex.(*Lexer).Row(), yylex.(*Lexer).Col())
                                                      $$.AddChild($2)
                                                      $$.AddChild($3)
                                                      log.ParseLog.Printf("compound_stmt0: %+v\n", $$)
																										}
                    ;

local_declarations  : local_declarations var_declaration
                                                    {
                                                      log.ParseLog.Printf("local_declarations0: %+v %+v %+v\n", $1, $2, yylex)
                                                      t := $1
                                                      if t != nil {
                                                        for t.Sibling() != nil {
                                                          t = t.Sibling()
                                                        }
                                                        t.SetSibling($2)
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
                                                        for t.Sibling() != nil {
                                                          t = t.Sibling()
                                                        }
                                                        t.SetSibling($2)
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
                                                      $$ = syntree.NewStmtSelectionNode()
                                                      $$.(syntree.Location).SetPos(yylex.(*Lexer).Row(), yylex.(*Lexer).Col())
                                                      $$.AddChild($3)
                                                      $$.AddChild($5)
                                                      log.ParseLog.Printf("selection_stmt0: %+v\n", $$)
																										}
                    | IF LPAREN expression RPAREN statement ELSE statement
                                                    {
                                                      log.ParseLog.Printf("selection_stmt1: %+v %+v %+v %+v %+v %+v %+v %+v\n", $1, $2, $3, $4, $5, $6, $7, yylex)
                                                      $$ = syntree.NewStmtSelectionNode()
                                                      $$.(syntree.Location).SetPos(yylex.(*Lexer).Row(), yylex.(*Lexer).Col())
                                                      $$.AddChild($3)
                                                      $$.AddChild($5)
                                                      $$.AddChild($7)
                                                      log.ParseLog.Printf("selection_stmt1: %+v\n", $$)
																										}
                    ;

iteration_stmt      : WHILE LPAREN expression RPAREN statement
                                                    {
                                                      log.ParseLog.Printf("iteration_stmt0: %+v %+v %+v %+v %+v %+v\n", $1, $2, $3, $4, $5, yylex)
                                                      $$ = syntree.NewStmtIterationNode()
                                                      $$.(syntree.Location).SetPos(yylex.(*Lexer).Row(), yylex.(*Lexer).Col())
                                                      $$.AddChild($3)
                                                      $$.AddChild($5)
                                                      log.ParseLog.Printf("iteration_stmt0: %+v\n", $$)
																										}
                    ;

return_stmt         : RETURN SEMI                   {
                                                      log.ParseLog.Printf("return_stmt0: %+v %+v %+v\n", $1, $2, yylex)
                                                      $$ = syntree.NewStmtReturnNode()
                                                      $$.(syntree.Location).SetPos(yylex.(*Lexer).Row(), yylex.(*Lexer).Col())
                                                      log.ParseLog.Printf("return_stmt0: %+v\n", $$)
																										}
                    | RETURN expression SEMI        {
                                                      log.ParseLog.Printf("return_stmt1: %+v %+v %+v %+v\n", $1, $2, $3, yylex)
                                                      $$ = syntree.NewStmtReturnNode()
                                                      $$.(syntree.Location).SetPos(yylex.(*Lexer).Row(), yylex.(*Lexer).Col())
                                                      $$.AddChild($2)
                                                      log.ParseLog.Printf("return_stmt1: %+v\n", $$)
                                                    }
                    ;

expression          : var ASSIGN expression         {
                                                      log.ParseLog.Printf("expression0: %+v %+v %+v %+v\n", $1, $2, $3, yylex)
                                                      $$ = syntree.NewExpAssignNode()
                                                      $$.(syntree.Location).SetPos(yylex.(*Lexer).Row(), yylex.(*Lexer).Col())
                                                      $$.AddChild($1)
                                                      $$.AddChild($3)
                                                      log.ParseLog.Printf("expression0: %+v\n", $$)
                                                    }
                    | simple_expression             {
                                                      log.ParseLog.Printf("expression1: %+v %+v\n", $1, yylex)
                                                      $$ = $1
                                                    }
                    ;

var                 : ID                            {
                                                      log.ParseLog.Printf("var0: %+v %+v\n", $1, yylex)
                                                      $$ = syntree.NewExpIdNode()
                                                      $$.(syntree.Location).SetPos(yylex.(*Lexer).Row(), yylex.(*Lexer).Col())
                                                      $$.(syntree.Name).SetName($1)
                                                      log.ParseLog.Printf("var0: %+v\n", $$)
                                                    }
                    | ID LBRACKET expression RBRACKET
                                                    {
                                                      log.ParseLog.Printf("var1: %+v %+v %+v %+v %+v\n", $1, $2, $3, $4, yylex)
                                                      $$ = syntree.NewExpIdArrayNode()
                                                      $$.(syntree.Location).SetPos(yylex.(*Lexer).Row(), yylex.(*Lexer).Col())
                                                      $$.(syntree.Name).SetName($1)
                                                      $$.AddChild($3)
                                                      log.ParseLog.Printf("var1: %+v\n", $$)
                                                    }
                    ;

simple_expression   : additive_expression relop additive_expression
                                                    {
                                                      log.ParseLog.Printf("simple_expression0: %+v %+v %+v %+v\n", $1, $2, $3, yylex)
                                                      $$ = syntree.NewExpOpNode()
                                                      $$.(syntree.Location).SetPos(yylex.(*Lexer).Row(), yylex.(*Lexer).Col())
                                                      $$.(syntree.TokType).SetTokType($2)
                                                      $$.AddChild($1)
                                                      $$.AddChild($3)
                                                      log.ParseLog.Printf("simple_expression0: %+v\n", $$)
                                                    }
                    | additive_expression           {
                                                      log.ParseLog.Printf("simple_expression6: %+v %+v\n", $1, yylex)
                                                      $$ = $1
                                                    }
                    ;

relop               : LT                            {
                                                      log.ParseLog.Printf("relop0: %+v %+v\n", $1, yylex)
                                                      $$ = syntree.LT
                                                      log.ParseLog.Printf("relop0: %+v\n", $$)
                                                    }
                    | LTE                           {
                                                      log.ParseLog.Printf("relop1: %+v %+v\n", $1, yylex)
                                                      $$ = syntree.LTE
                                                      log.ParseLog.Printf("relop1: %+v\n", $$)
                                                    }
                    | GT                            {
                                                      log.ParseLog.Printf("relop2: %+v %+v\n", $1, yylex)
                                                      $$ = syntree.GT
                                                      log.ParseLog.Printf("relop2: %+v\n", $$)
                                                    }
                    | GTE                           {
                                                      log.ParseLog.Printf("relop3: %+v %+v\n", $1, yylex)
                                                      $$ = syntree.GTE
                                                      log.ParseLog.Printf("relop3: %+v\n", $$)
                                                    }
                    | EQ                            {
                                                      log.ParseLog.Printf("relop4: %+v %+v\n", $1, yylex)
                                                      $$ = syntree.EQ
                                                      log.ParseLog.Printf("relop4: %+v\n", $$)
                                                    }
                    | NEQ                           {
                                                      log.ParseLog.Printf("relop5: %+v %+v\n", $1, yylex)
                                                      $$ = syntree.NEQ
                                                      log.ParseLog.Printf("relop5: %+v\n", $$)
                                                    }

additive_expression : additive_expression addop term
                                                    {
																											log.ParseLog.Printf("additive_expression0: %+v %+v %+v %+v\n", $1, $2, $3, yylex)
                                                      $$ = syntree.NewExpOpNode()
                                                      $$.(syntree.Location).SetPos(yylex.(*Lexer).Row(), yylex.(*Lexer).Col())
                                                      $$.(syntree.TokType).SetTokType($2)
                                                      $$.AddChild($1)
                                                      $$.AddChild($3)
                                                      log.ParseLog.Printf("additive_expression0: %+v\n", $$)
																										}
                    | term                          {
																											log.ParseLog.Printf("additive_expression1: %+v %+v\n", $1, yylex)
                                                      $$ = $1
																										}
                    ;

addop               : PLUS                          {
                                                      log.ParseLog.Printf("addop0: %+v %+v\n", $1, yylex)
                                                      $$ = syntree.PLUS
                                                      log.ParseLog.Printf("addop0: %+v\n", $$)
                                                    }
                    | MINUS                         {
                                                      log.ParseLog.Printf("addop1: %+v %+v\n", $1, yylex)
                                                      $$ = syntree.MINUS
                                                      log.ParseLog.Printf("addop1: %+v\n", $$)
                                                    }

term                : term mulop factor             {
																											log.ParseLog.Printf("term0: %+v %+v %+v %+v\n", $1, $2, $3, yylex)
                                                      $$ = syntree.NewExpOpNode()
                                                      $$.(syntree.Location).SetPos(yylex.(*Lexer).Row(), yylex.(*Lexer).Col())
                                                      $$.(syntree.TokType).SetTokType($2)
                                                      $$.AddChild($1)
                                                      $$.AddChild($3)
                                                      log.ParseLog.Printf("term0: %+v\n", $$)
																										}
                    | factor                        {
																											log.ParseLog.Printf("term1: %+v %+v\n", $1, yylex)
                                                      $$ = $1
																										}
                    ;

mulop               : TIMES                         {
                                                      log.ParseLog.Printf("mulop0: %+v %+v\n", $1, yylex)
                                                      $$ = syntree.TIMES
                                                      log.ParseLog.Printf("mulop0: %+v\n", $$)
                                                    }
                    | OVER                          {
                                                      log.ParseLog.Printf("mulop1: %+v %+v\n", $1, yylex)
                                                      $$ = syntree.OVER
                                                      log.ParseLog.Printf("mulop1: %+v\n", $$)
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
                                                      $$ = syntree.NewExpConstNode()
                                                      $$.(syntree.Location).SetPos(yylex.(*Lexer).Row(), yylex.(*Lexer).Col())
                                                      v, _ := strconv.Atoi(yylex.(*Lexer).Text())
                                                      $$.(syntree.Value).SetValue(v)
                                                      log.ParseLog.Printf("factor3: %+v %+v\n", $$)
																										}
                    ;

call                : ID LPAREN args RPAREN         {
																											log.ParseLog.Printf("call0: %+v %+v %+v %+v %+v\n", $1, $2, $3, $4, yylex)
                                                      $$ = syntree.NewExpCallNode()
                                                      $$.(syntree.Location).SetPos(yylex.(*Lexer).Row(), yylex.(*Lexer).Col())
                                                      $$.(syntree.Name).SetName($1)
                                                      $$.AddChild($3)
                                                      log.ParseLog.Printf("call0: %+v\n", $$)
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
                                                        for t.Sibling() != nil {
                                                          t = t.Sibling()
                                                        }
                                                        t.SetSibling($3)
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
