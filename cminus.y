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
                                                      root = $1
                                                      log.ParseLog.Printf("program0: %+v\n", root)
                                                    }
                    ;

declaration_list    : declaration_list declaration  {
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
                                                      log.ParseLog.Printf("declaration_list0: %+v\n", $$)
                                                    }
                    | declaration                   {
                                                      $$ = $1
                                                      log.ParseLog.Printf("declaration_list1: %+v\n", $$)
                                                    }
                    ;

declaration         : var_declaration               {
                                                      $$ = $1
                                                      log.ParseLog.Printf("declaration0: %+v\n", $$)
                                                    }
                    | fun_declaration               {
                                                      $$ = $1
                                                      log.ParseLog.Printf("declaration1: %+v\n", $$)
                                                    }
                    ;

var_declaration     : type_specifier ID SEMI        {
                                                      $$ = syntree.NewExpVarNode()
                                                      $$.(syntree.Location).SetPos(yylex.(*Lexer).Row(), yylex.(*Lexer).Col())
                                                      $$.(syntree.ExpType).SetExpType($1)
                                                      $$.(syntree.Name).SetName($2)
                                                      log.ParseLog.Printf("var_declaration0: %+v\n", $$)
                                                    }
                    | type_specifier ID LBRACKET NUM RBRACKET SEMI
                                                    {
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
                                                      $$ = syntree.INTEGER_TYPE
                                                      log.ParseLog.Printf("type_specifier0: %+v\n", $$)
                                                    }
                    | VOID                          {
                                                      $$ = syntree.VOID_TYPE
                                                      log.ParseLog.Printf("type_specifier1: %+v\n", $$)
                                                    }
                    ;

fun_declaration     : type_specifier ID LPAREN params RPAREN compound_stmt
                                                    {
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
                                                      $$ = $1
                                                      log.ParseLog.Printf("params0: %+v\n", $$)
																										}
                    | VOID                          {
                                                      $$ = syntree.NewExpParamNode()
                                                      $$.(syntree.Location).SetPos(yylex.(*Lexer).Row(), yylex.(*Lexer).Col())
                                                      $$.(syntree.ExpType).SetExpType(syntree.VOID_TYPE)
                                                      log.ParseLog.Printf("params1: %+v\n", $$)
																										}
                    ;

param_list          : param_list COMMA param        {
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
                                                      log.ParseLog.Printf("param_list0: %+v\n", $$)
																										}
                    | param                         {
                                                      $$ = $1
                                                      log.ParseLog.Printf("param_list1: %+v\n", $$)
																										}
                    ;

param               : type_specifier ID             {
                                                      $$ = syntree.NewExpParamNode()
                                                      $$.(syntree.Location).SetPos(yylex.(*Lexer).Row(), yylex.(*Lexer).Col())
                                                      $$.(syntree.ExpType).SetExpType($1)
                                                      $$.(syntree.Name).SetName($2)
                                                      log.ParseLog.Printf("param0: %+v\n", $$)
																										}
                    | type_specifier ID LBRACKET RBRACKET
                                                    {
                                                      $$ = syntree.NewExpParamArrayNode()
                                                      $$.(syntree.Location).SetPos(yylex.(*Lexer).Row(), yylex.(*Lexer).Col())
                                                      $$.(syntree.ExpType).SetExpType($1)
                                                      $$.(syntree.Name).SetName($2)
                                                      log.ParseLog.Printf("param1: %+v\n", $$)
																										}
                    ;

compound_stmt       : LBRACE local_declarations statement_list RBRACE
                                                    {
                                                      $$ = syntree.NewStmtCompoundNode()
                                                      $$.(syntree.Location).SetPos(yylex.(*Lexer).Row(), yylex.(*Lexer).Col())
                                                      $$.AddChild($2)
                                                      $$.AddChild($3)
                                                      log.ParseLog.Printf("compound_stmt0: %+v\n", $$)
																										}
                    ;

local_declarations  : local_declarations var_declaration
                                                    {
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
                                                      log.ParseLog.Printf("local_declarations0: %+v\n", $$)
																										}
                    | empty                         {
                                                      $$ = $1
                                                      log.ParseLog.Printf("local_declarations1: %+v\n", $$)
																										}
                    ;

statement_list      : statement_list statement      {
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
                                                      log.ParseLog.Printf("statement_list0: %+v\n", $$)
																										}
                    | empty                         {
                                                      $$ = $1
                                                      log.ParseLog.Printf("statement_list1: %+v\n", $$)
                                                    }
                    ;

statement           : expression_stmt               {
                                                      $$ = $1
                                                      log.ParseLog.Printf("statement0: %+v\n", $$)
																										}
                    | compound_stmt                 {
                                                      $$ = $1
                                                      log.ParseLog.Printf("statement1: %+v\n", $$)
																										}
                    | selection_stmt                {
                                                      $$ = $1
                                                      log.ParseLog.Printf("statement2: %+v\n", $$)
																										}
                    | iteration_stmt                {
                                                      $$ = $1
                                                      log.ParseLog.Printf("statement3: %+v\n", $$)
																										}
                    | return_stmt                   {
                                                      $$ = $1
                                                      log.ParseLog.Printf("statement4: %+v\n", $$)
																										}
                    ;

expression_stmt     : expression SEMI               {
                                                      $$ = $1
                                                      log.ParseLog.Printf("expression_stmt0: %+v\n", $$)
																										}
                    | SEMI                          {
                                                      $$ = nil
                                                      log.ParseLog.Printf("expression_stmt1: nil\n")
																										}
                    ;

selection_stmt      : IF LPAREN expression RPAREN statement %prec THEN
                                                    {
                                                      $$ = syntree.NewStmtSelectionNode()
                                                      $$.(syntree.Location).SetPos(yylex.(*Lexer).Row(), yylex.(*Lexer).Col())
                                                      $$.AddChild($3)
                                                      $$.AddChild($5)
                                                      log.ParseLog.Printf("selection_stmt0: %+v\n", $$)
																										}
                    | IF LPAREN expression RPAREN statement ELSE statement
                                                    {
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
                                                      $$ = syntree.NewStmtIterationNode()
                                                      $$.(syntree.Location).SetPos(yylex.(*Lexer).Row(), yylex.(*Lexer).Col())
                                                      $$.AddChild($3)
                                                      $$.AddChild($5)
                                                      log.ParseLog.Printf("iteration_stmt0: %+v\n", $$)
																										}
                    ;

return_stmt         : RETURN SEMI                   {
                                                      $$ = syntree.NewStmtReturnNode()
                                                      $$.(syntree.Location).SetPos(yylex.(*Lexer).Row(), yylex.(*Lexer).Col())
                                                      log.ParseLog.Printf("return_stmt0: %+v\n", $$)
																										}
                    | RETURN expression SEMI        {
                                                      $$ = syntree.NewStmtReturnNode()
                                                      $$.(syntree.Location).SetPos(yylex.(*Lexer).Row(), yylex.(*Lexer).Col())
                                                      $$.AddChild($2)
                                                      log.ParseLog.Printf("return_stmt1: %+v\n", $$)
                                                    }
                    ;

expression          : var ASSIGN expression         {
                                                      $$ = syntree.NewExpAssignNode()
                                                      $$.(syntree.Location).SetPos(yylex.(*Lexer).Row(), yylex.(*Lexer).Col())
                                                      $$.AddChild($1)
                                                      $$.AddChild($3)
                                                      log.ParseLog.Printf("expression0: %+v\n", $$)
                                                    }
                    | simple_expression             {
                                                      $$ = $1
                                                      log.ParseLog.Printf("expression1: %+v\n", $$)
                                                    }
                    ;

var                 : ID                            {
                                                      $$ = syntree.NewExpIdNode()
                                                      $$.(syntree.Location).SetPos(yylex.(*Lexer).Row(), yylex.(*Lexer).Col())
                                                      $$.(syntree.Name).SetName($1)
                                                      log.ParseLog.Printf("var0: %+v\n", $$)
                                                    }
                    | ID LBRACKET expression RBRACKET
                                                    {
                                                      $$ = syntree.NewExpIdArrayNode()
                                                      $$.(syntree.Location).SetPos(yylex.(*Lexer).Row(), yylex.(*Lexer).Col())
                                                      $$.(syntree.Name).SetName($1)
                                                      $$.AddChild($3)
                                                      log.ParseLog.Printf("var1: %+v\n", $$)
                                                    }
                    ;

simple_expression   : additive_expression relop additive_expression
                                                    {
                                                      $$ = syntree.NewExpOpNode()
                                                      $$.(syntree.Location).SetPos(yylex.(*Lexer).Row(), yylex.(*Lexer).Col())
                                                      $$.(syntree.TokType).SetTokType($2)
                                                      $$.AddChild($1)
                                                      $$.AddChild($3)
                                                      log.ParseLog.Printf("simple_expression0: %+v\n", $$)
                                                    }
                    | additive_expression           {
                                                      $$ = $1
                                                      log.ParseLog.Printf("simple_expression1: %+v\n", $$)
                                                    }
                    ;

relop               : LT                            {
                                                      $$ = syntree.LT
                                                      log.ParseLog.Printf("relop0: %+v\n", $$)
                                                    }
                    | LTE                           {
                                                      $$ = syntree.LTE
                                                      log.ParseLog.Printf("relop1: %+v\n", $$)
                                                    }
                    | GT                            {
                                                      $$ = syntree.GT
                                                      log.ParseLog.Printf("relop2: %+v\n", $$)
                                                    }
                    | GTE                           {
                                                      $$ = syntree.GTE
                                                      log.ParseLog.Printf("relop3: %+v\n", $$)
                                                    }
                    | EQ                            {
                                                      $$ = syntree.EQ
                                                      log.ParseLog.Printf("relop4: %+v\n", $$)
                                                    }
                    | NEQ                           {
                                                      $$ = syntree.NEQ
                                                      log.ParseLog.Printf("relop5: %+v\n", $$)
                                                    }

additive_expression : additive_expression addop term
                                                    {
                                                      $$ = syntree.NewExpOpNode()
                                                      $$.(syntree.Location).SetPos(yylex.(*Lexer).Row(), yylex.(*Lexer).Col())
                                                      $$.(syntree.TokType).SetTokType($2)
                                                      $$.AddChild($1)
                                                      $$.AddChild($3)
                                                      log.ParseLog.Printf("additive_expression0: %+v\n", $$)
																										}
                    | term                          {
                                                      $$ = $1
                                                      log.ParseLog.Printf("additive_expression1: %+v\n", $$)
																										}
                    ;

addop               : PLUS                          {
                                                      $$ = syntree.PLUS
                                                      log.ParseLog.Printf("addop0: %+v\n", $$)
                                                    }
                    | MINUS                         {
                                                      $$ = syntree.MINUS
                                                      log.ParseLog.Printf("addop1: %+v\n", $$)
                                                    }

term                : term mulop factor             {
                                                      $$ = syntree.NewExpOpNode()
                                                      $$.(syntree.Location).SetPos(yylex.(*Lexer).Row(), yylex.(*Lexer).Col())
                                                      $$.(syntree.TokType).SetTokType($2)
                                                      $$.AddChild($1)
                                                      $$.AddChild($3)
                                                      log.ParseLog.Printf("term0: %+v\n", $$)
																										}
                    | factor                        {
                                                      $$ = $1
                                                      log.ParseLog.Printf("term1: %+v\n", $$)
																										}
                    ;

mulop               : TIMES                         {
                                                      $$ = syntree.TIMES
                                                      log.ParseLog.Printf("mulop0: %+v\n", $$)
                                                    }
                    | OVER                          {
                                                      $$ = syntree.OVER
                                                      log.ParseLog.Printf("mulop1: %+v\n", $$)
                                                    }

factor              : LPAREN expression RPAREN      {
                                                      $$ = $2
                                                      log.ParseLog.Printf("factor0: %+v\n", $$)
																										}
                    | var                           {
                                                      $$ = $1
                                                      log.ParseLog.Printf("factor1: %+v\n", $$)
																										}
                    | call                          {
                                                      $$ = $1
                                                      log.ParseLog.Printf("factor2: %+v\n", $$)
																										}
                    | NUM                           {
                                                      $$ = syntree.NewExpConstNode()
                                                      $$.(syntree.Location).SetPos(yylex.(*Lexer).Row(), yylex.(*Lexer).Col())
                                                      v, _ := strconv.Atoi(yylex.(*Lexer).Text())
                                                      $$.(syntree.Value).SetValue(v)
                                                      log.ParseLog.Printf("factor3: %+v\n", $$)
																										}
                    ;

call                : ID LPAREN args RPAREN         {
                                                      $$ = syntree.NewExpCallNode()
                                                      $$.(syntree.Location).SetPos(yylex.(*Lexer).Row(), yylex.(*Lexer).Col())
                                                      $$.(syntree.Name).SetName($1)
                                                      $$.AddChild($3)
                                                      log.ParseLog.Printf("call0: %+v\n", $$)
																										}
                    ;

args                : args_list                     {
                                                      $$ = $1
                                                      log.ParseLog.Printf("args0: %+v\n", $$)
																										}
                    | empty                         {
                                                      $$ = $1
                                                      log.ParseLog.Printf("args1: %+v\n", $$)
																										}
                    ;

args_list           : args_list COMMA expression    {
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
                                                      log.ParseLog.Printf("args_list0: %+v\n", $$)
                                                    }
                    | expression                    {
                                                      $$ = $1
                                                      log.ParseLog.Printf("args_list1: %+v\n", $$)
                                                    }
                    ;

empty               : /* empty */                   {
                                                      $$ = nil
                                                      log.ParseLog.Printf("empty0: nil\n")
                                                    }
                    ;

%%
