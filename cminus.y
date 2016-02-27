%{
package main

import (
 "github.com/bantl23/cminus/log"
)

var savedName string
%}

%union {
  str string
}

%token <str> IF ELSE INT RETURN VOID WHILE PLUS MINUS TIMES
%token <str> OVER LT LTE GT GTE EQ NEQ ASSIGN SEMI COMMA
%token <str> LPAREN RPAREN LBRACE RBRACE LBRACKET RBRACKET
%token <str> ID NUM

%type <str> program declaration_list declaration var_declaration
%type <str> type_specifier fun_declaration params param_list param
%type <str> compound_stmt local_declarations statement_list statement
%type <str> expression_stmt selection_stmt iteration_stmt return_stmt
%type <str> expression var simple_expression relop additive_expression addop
%type <str> term mulop factor call args args_list empty

%nonassoc THEN
%nonassoc ELSE

%%

program             : declaration_list              {
                                                      log.Trace.Printf("program0: %+v %v\n", $1, yylex)
                                                    }
                    ;

declaration_list    : declaration_list declaration  {
                                                      log.Trace.Printf("declaration_list0: %+v %+v %v\n", $1, $2, yylex)
                                                    }
                    | declaration                   {
                                                      log.Trace.Printf("declaration_list1: %+v %v\n", $1, yylex)
                                                    }
                    ;

declaration         : var_declaration               {
                                                      log.Trace.Printf("declaration0: %+v %v\n", $1, yylex)
                                                    }
                    | fun_declaration               {
                                                      log.Trace.Printf("declaration1: %+v %v\n", $1, yylex)
                                                    }
                    ;

var_declaration     : type_specifier ID SEMI        {
                                                      log.Trace.Printf("var_declaration0: %+v %+v %+v %v\n", $1, $2, $3, yylex)
                                                    }
                    | type_specifier ID LBRACKET NUM RBRACKET SEMI
                                                    {
                                                      log.Trace.Printf("var_declaration1: %+v %+v %+v %+v %+v %+v %+v\n", $1, $2, $3, $4, $5, $6, yylex)
                                                    }
                    ;

type_specifier      : INT                           {
                                                      log.Trace.Printf("type_specifier0: %+v %+v\n", $1, yylex)
                                                    }
                    | VOID                          {
                                                      log.Trace.Printf("type_specifier1: %+v %+v\n", $1, yylex)
                                                    }
                    ;

fun_declaration     : type_specifier ID LPAREN params RPAREN compound_stmt
                                                    {
                                                      log.Trace.Printf("fun_declaration0: %+v %+v %+v %+v %+v %+v %+v\n", $1, $2, $3, $4, $5, $6, yylex)
                                                    }
                    ;

params              : param_list                    {
                                                      log.Trace.Printf("params0: %+v %+v\n", $1, yylex)
																										}
                    | VOID                          {
                                                      log.Trace.Printf("params1: %+v %+v\n", $1, yylex)
																										}
                    ;

param_list          : param_list COMMA param        {
                                                      log.Trace.Printf("param_list0: %+v %+v %+v %+v\n", $1, $2, $3, yylex)
																										}
                    | param                         {
                                                      log.Trace.Printf("param_list1: %+v %+v\n", $1, yylex)
																										}
                    ;

param               : type_specifier ID             {
                                                      log.Trace.Printf("param0: %+v %+v %+v\n", $1, $2, yylex)
																										}
                    | type_specifier ID LBRACKET RBRACKET
                                                    {
                                                      log.Trace.Printf("param1: %+v %+v %+v %+v %+v\n", $1, $2, $3, $3, yylex)
																										}
                    ;

compound_stmt       : LBRACE local_declarations statement_list RBRACE
                                                    {
                                                      log.Trace.Printf("compound_stmt0: %+v %+v %+v %+v %+v\n", $1, $2, $3, $4, yylex)
																										}
                    ;

local_declarations  : local_declarations var_declaration
                                                    {
                                                      log.Trace.Printf("local_declarations0: %+v %+v %+v\n", $1, $2, yylex)
																										}
                    | empty                         {
                                                      log.Trace.Printf("local_declarations1: %+v %+v\n", $1, yylex)
																										}
                    ;

statement_list      : statement_list statement      {
                                                      log.Trace.Printf("statement_list0: %+v %+v %+v\n", $1, $2, yylex)
																										}
                    | empty                         {
                                                      log.Trace.Printf("statement_list1: %+v %+v\n", $1, yylex)
                                                    }
                    ;

statement           : expression_stmt               {
                                                      log.Trace.Printf("statement0: %+v %+v\n", $1, yylex)
																										}
                    | compound_stmt                 {
                                                      log.Trace.Printf("statement1: %+v %+v\n", $1, yylex)
																										}
                    | selection_stmt                {
                                                      log.Trace.Printf("statement2: %+v %+v\n", $1, yylex)
																										}
                    | iteration_stmt                {
                                                      log.Trace.Printf("statement3: %+v %+v\n", $1, yylex)
																										}
                    | return_stmt                   {
                                                      log.Trace.Printf("statement4: %+v %+v\n", $1, yylex)
																										}
                    ;

expression_stmt     : expression SEMI               {
                                                      log.Trace.Printf("expression_stmt0: %+v %+v %+v\n", $1, $2, yylex)
																										}
                    | SEMI                          {
                                                      log.Trace.Printf("expression_stmt1: %+v %+v\n", $1, yylex)
																										}
                    ;

selection_stmt      : IF LPAREN expression RPAREN statement %prec THEN
                                                    {
                                                      log.Trace.Printf("selection_stmt0: %+v %+v %+v %+v %+v %+v\n", $1, $2, $3, $4, $5, yylex)
																										}
                    | IF LPAREN expression RPAREN statement ELSE statement
                                                    {
                                                      log.Trace.Printf("selection_stmt1: %+v %+v %+v %+v %+v %+v %+v %+v\n", $1, $2, $3, $4, $5, $6, $7, yylex)
																										}
                    ;

iteration_stmt      : WHILE LPAREN expression RPAREN statement
                                                    {
                                                      log.Trace.Printf("iteration_stmt0: %+v %+v %+v %+v %+v %+v\n", $1, $2, $3, $4, $5, yylex)
																										}
                    ;

return_stmt         : RETURN SEMI                   {
                                                      log.Trace.Printf("return_stmt0: %+v %+v %+v\n", $1, $2, yylex)
																										}
                    | RETURN expression SEMI        {
                                                      log.Trace.Printf("return_stmt1: %+v %+v %+v %+v\n", $1, $2, $3, yylex)
                                                    }
                    ;

expression          : var ASSIGN expression         {
                                                      log.Trace.Printf("expression0: %+v %+v %+v %+v\n", $1, $2, $3, yylex)
                                                    }
                    | simple_expression             {
                                                      log.Trace.Printf("expression1: %+v %+v\n", $1, yylex)
                                                    }
                    ;

var                 : ID                            {
                                                      log.Trace.Printf("var0: %+v %+v\n", $1, yylex)
                                                    }
                    | ID LBRACKET expression RBRACKET
                                                    {
                                                      log.Trace.Printf("var1: %+v %+v %+v %+v %+v\n", $1, $2, $3, $4, yylex)
                                                    }
                    ;

simple_expression   : additive_expression relop additive_expression
                                                    {
                                                      log.Trace.Printf("simple_expression0: %+v %+v %+v %+v\n", $1, $2, $3, yylex)
                                                    }
                    | additive_expression           {
                                                      log.Trace.Printf("simple_expression1: %+v %+v\n", $1, yylex)
                                                    }
                    ;

relop               : LTE                           {
                                                      log.Trace.Printf("relop0: %+v %+v\n", $1, yylex)
                                                    }
                    | LT                            {
                                                      log.Trace.Printf("relop1: %+v %+v\n", $1, yylex)
                                                    }
                    | GT                            {
                                                      log.Trace.Printf("relop2: %+v %+v\n", $1, yylex)
                                                    }
                    | GTE                           {
                                                      log.Trace.Printf("relop3: %+v %+v\n", $1, yylex)
                                                    }
                    | EQ                            {
                                                      log.Trace.Printf("relop4: %+v %+v\n", $1, yylex)
                                                    }
                    | NEQ                           {
                                                      log.Trace.Printf("relop5: %+v %+v\n", $1, yylex)
                                                    }
                    ;

additive_expression : additive_expression addop term
                                                    {
																											log.Trace.Printf("additive_expression0: %+v %+v %+v %+v\n", $1, $2, $3, yylex)
																										}
                    | term                          {
																											log.Trace.Printf("additive_expression1: %+v %+v\n", $1, yylex)
																										}
                    ;

addop               : PLUS                          {
																											log.Trace.Printf("addop0: %+v %+v\n", $1, yylex)
																										}
                    | MINUS                         {
																											log.Trace.Printf("addop1: %+v %+v\n", $1, yylex)
																										}
                    ;

term                : term mulop factor             {
																											log.Trace.Printf("term0: %+v %+v %+v %+v\n", $1, $2, $3, yylex)
																										}
                    | factor                        {
																											log.Trace.Printf("term1: %+v %+v\n", $1, yylex)
																										}
                    ;

mulop               : TIMES                         {
																											log.Trace.Printf("mulop0: %+v %+v\n", $1, yylex)
																										}
                    | OVER                          {
																											log.Trace.Printf("mulop1: %+v %+v\n", $1, yylex)
																										}
                    ;

factor              : LPAREN expression RPAREN      {
																											log.Trace.Printf("factor0: %+v %+v %+v %+v\n", $1, $2, $3, yylex)
																										}
                    | var                           {
																											log.Trace.Printf("factor1: %+v %+v\n", $1, yylex)
																										}
                    | call                          {
																											log.Trace.Printf("factor2: %+v %+v\n", $1, yylex)
																										}
                    | NUM                           {
																											log.Trace.Printf("factor3: %+v %+v\n", $1, yylex)
																										}
                    ;

call                : ID LPAREN args RPAREN         {
																											log.Trace.Printf("call0: %+v %+v %+v %+v %+v\n", $1, $2, $3, $4, yylex)
																										}
                    ;

args                : args_list                     {
																											log.Trace.Printf("args0: %+v %+v\n", $1, yylex)
																										}
                    | empty                         {
																											log.Trace.Printf("args1: %+v %+v\n", $1, yylex)
																										}
                    ;

args_list           : args_list COMMA expression    {
                                                      log.Trace.Printf("args_list0: %+v %+v %+v %+v\n", $1, $2, $3, yylex)
                                                    }
                    | expression                    {
                                                      log.Trace.Printf("args_list1: %+v %+v\n", $1, yylex)
                                                    }
                    ;

empty               : /* empty */                   {
                                                      log.Trace.Printf("empty0: %+v\n", yylex)
                                                    }
                    ;

%%
