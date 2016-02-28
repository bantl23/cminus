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
                                                      log.ParseLog.Printf("program0: %+v %v\n", $1, yylex)
                                                    }
                    ;

declaration_list    : declaration_list declaration  {
                                                      log.ParseLog.Printf("declaration_list0: %+v %+v %v\n", $1, $2, yylex)
                                                    }
                    | declaration                   {
                                                      log.ParseLog.Printf("declaration_list1: %+v %v\n", $1, yylex)
                                                    }
                    ;

declaration         : var_declaration               {
                                                      log.ParseLog.Printf("declaration0: %+v %v\n", $1, yylex)
                                                    }
                    | fun_declaration               {
                                                      log.ParseLog.Printf("declaration1: %+v %v\n", $1, yylex)
                                                    }
                    ;

var_declaration     : type_specifier ID SEMI        {
                                                      log.ParseLog.Printf("var_declaration0: %+v %+v %+v %v\n", $1, $2, $3, yylex)
                                                    }
                    | type_specifier ID LBRACKET NUM RBRACKET SEMI
                                                    {
                                                      log.ParseLog.Printf("var_declaration1: %+v %+v %+v %+v %+v %+v %+v\n", $1, $2, $3, $4, $5, $6, yylex)
                                                    }
                    ;

type_specifier      : INT                           {
                                                      log.ParseLog.Printf("type_specifier0: %+v %+v\n", $1, yylex)
                                                    }
                    | VOID                          {
                                                      log.ParseLog.Printf("type_specifier1: %+v %+v\n", $1, yylex)
                                                    }
                    ;

fun_declaration     : type_specifier ID LPAREN params RPAREN compound_stmt
                                                    {
                                                      log.ParseLog.Printf("fun_declaration0: %+v %+v %+v %+v %+v %+v %+v\n", $1, $2, $3, $4, $5, $6, yylex)
                                                    }
                    ;

params              : param_list                    {
                                                      log.ParseLog.Printf("params0: %+v %+v\n", $1, yylex)
																										}
                    | VOID                          {
                                                      log.ParseLog.Printf("params1: %+v %+v\n", $1, yylex)
																										}
                    ;

param_list          : param_list COMMA param        {
                                                      log.ParseLog.Printf("param_list0: %+v %+v %+v %+v\n", $1, $2, $3, yylex)
																										}
                    | param                         {
                                                      log.ParseLog.Printf("param_list1: %+v %+v\n", $1, yylex)
																										}
                    ;

param               : type_specifier ID             {
                                                      log.ParseLog.Printf("param0: %+v %+v %+v\n", $1, $2, yylex)
																										}
                    | type_specifier ID LBRACKET RBRACKET
                                                    {
                                                      log.ParseLog.Printf("param1: %+v %+v %+v %+v %+v\n", $1, $2, $3, $3, yylex)
																										}
                    ;

compound_stmt       : LBRACE local_declarations statement_list RBRACE
                                                    {
                                                      log.ParseLog.Printf("compound_stmt0: %+v %+v %+v %+v %+v\n", $1, $2, $3, $4, yylex)
																										}
                    ;

local_declarations  : local_declarations var_declaration
                                                    {
                                                      log.ParseLog.Printf("local_declarations0: %+v %+v %+v\n", $1, $2, yylex)
																										}
                    | empty                         {
                                                      log.ParseLog.Printf("local_declarations1: %+v %+v\n", $1, yylex)
																										}
                    ;

statement_list      : statement_list statement      {
                                                      log.ParseLog.Printf("statement_list0: %+v %+v %+v\n", $1, $2, yylex)
																										}
                    | empty                         {
                                                      log.ParseLog.Printf("statement_list1: %+v %+v\n", $1, yylex)
                                                    }
                    ;

statement           : expression_stmt               {
                                                      log.ParseLog.Printf("statement0: %+v %+v\n", $1, yylex)
																										}
                    | compound_stmt                 {
                                                      log.ParseLog.Printf("statement1: %+v %+v\n", $1, yylex)
																										}
                    | selection_stmt                {
                                                      log.ParseLog.Printf("statement2: %+v %+v\n", $1, yylex)
																										}
                    | iteration_stmt                {
                                                      log.ParseLog.Printf("statement3: %+v %+v\n", $1, yylex)
																										}
                    | return_stmt                   {
                                                      log.ParseLog.Printf("statement4: %+v %+v\n", $1, yylex)
																										}
                    ;

expression_stmt     : expression SEMI               {
                                                      log.ParseLog.Printf("expression_stmt0: %+v %+v %+v\n", $1, $2, yylex)
																										}
                    | SEMI                          {
                                                      log.ParseLog.Printf("expression_stmt1: %+v %+v\n", $1, yylex)
																										}
                    ;

selection_stmt      : IF LPAREN expression RPAREN statement %prec THEN
                                                    {
                                                      log.ParseLog.Printf("selection_stmt0: %+v %+v %+v %+v %+v %+v\n", $1, $2, $3, $4, $5, yylex)
																										}
                    | IF LPAREN expression RPAREN statement ELSE statement
                                                    {
                                                      log.ParseLog.Printf("selection_stmt1: %+v %+v %+v %+v %+v %+v %+v %+v\n", $1, $2, $3, $4, $5, $6, $7, yylex)
																										}
                    ;

iteration_stmt      : WHILE LPAREN expression RPAREN statement
                                                    {
                                                      log.ParseLog.Printf("iteration_stmt0: %+v %+v %+v %+v %+v %+v\n", $1, $2, $3, $4, $5, yylex)
																										}
                    ;

return_stmt         : RETURN SEMI                   {
                                                      log.ParseLog.Printf("return_stmt0: %+v %+v %+v\n", $1, $2, yylex)
																										}
                    | RETURN expression SEMI        {
                                                      log.ParseLog.Printf("return_stmt1: %+v %+v %+v %+v\n", $1, $2, $3, yylex)
                                                    }
                    ;

expression          : var ASSIGN expression         {
                                                      log.ParseLog.Printf("expression0: %+v %+v %+v %+v\n", $1, $2, $3, yylex)
                                                    }
                    | simple_expression             {
                                                      log.ParseLog.Printf("expression1: %+v %+v\n", $1, yylex)
                                                    }
                    ;

var                 : ID                            {
                                                      log.ParseLog.Printf("var0: %+v %+v\n", $1, yylex)
                                                    }
                    | ID LBRACKET expression RBRACKET
                                                    {
                                                      log.ParseLog.Printf("var1: %+v %+v %+v %+v %+v\n", $1, $2, $3, $4, yylex)
                                                    }
                    ;

simple_expression   : additive_expression relop additive_expression
                                                    {
                                                      log.ParseLog.Printf("simple_expression0: %+v %+v %+v %+v\n", $1, $2, $3, yylex)
                                                    }
                    | additive_expression           {
                                                      log.ParseLog.Printf("simple_expression1: %+v %+v\n", $1, yylex)
                                                    }
                    ;

relop               : LTE                           {
                                                      log.ParseLog.Printf("relop0: %+v %+v\n", $1, yylex)
                                                    }
                    | LT                            {
                                                      log.ParseLog.Printf("relop1: %+v %+v\n", $1, yylex)
                                                    }
                    | GT                            {
                                                      log.ParseLog.Printf("relop2: %+v %+v\n", $1, yylex)
                                                    }
                    | GTE                           {
                                                      log.ParseLog.Printf("relop3: %+v %+v\n", $1, yylex)
                                                    }
                    | EQ                            {
                                                      log.ParseLog.Printf("relop4: %+v %+v\n", $1, yylex)
                                                    }
                    | NEQ                           {
                                                      log.ParseLog.Printf("relop5: %+v %+v\n", $1, yylex)
                                                    }
                    ;

additive_expression : additive_expression addop term
                                                    {
																											log.ParseLog.Printf("additive_expression0: %+v %+v %+v %+v\n", $1, $2, $3, yylex)
																										}
                    | term                          {
																											log.ParseLog.Printf("additive_expression1: %+v %+v\n", $1, yylex)
																										}
                    ;

addop               : PLUS                          {
																											log.ParseLog.Printf("addop0: %+v %+v\n", $1, yylex)
																										}
                    | MINUS                         {
																											log.ParseLog.Printf("addop1: %+v %+v\n", $1, yylex)
																										}
                    ;

term                : term mulop factor             {
																											log.ParseLog.Printf("term0: %+v %+v %+v %+v\n", $1, $2, $3, yylex)
																										}
                    | factor                        {
																											log.ParseLog.Printf("term1: %+v %+v\n", $1, yylex)
																										}
                    ;

mulop               : TIMES                         {
																											log.ParseLog.Printf("mulop0: %+v %+v\n", $1, yylex)
																										}
                    | OVER                          {
																											log.ParseLog.Printf("mulop1: %+v %+v\n", $1, yylex)
																										}
                    ;

factor              : LPAREN expression RPAREN      {
																											log.ParseLog.Printf("factor0: %+v %+v %+v %+v\n", $1, $2, $3, yylex)
																										}
                    | var                           {
																											log.ParseLog.Printf("factor1: %+v %+v\n", $1, yylex)
																										}
                    | call                          {
																											log.ParseLog.Printf("factor2: %+v %+v\n", $1, yylex)
																										}
                    | NUM                           {
																											log.ParseLog.Printf("factor3: %+v %+v\n", $1, yylex)
																										}
                    ;

call                : ID LPAREN args RPAREN         {
																											log.ParseLog.Printf("call0: %+v %+v %+v %+v %+v\n", $1, $2, $3, $4, yylex)
																										}
                    ;

args                : args_list                     {
																											log.ParseLog.Printf("args0: %+v %+v\n", $1, yylex)
																										}
                    | empty                         {
																											log.ParseLog.Printf("args1: %+v %+v\n", $1, yylex)
																										}
                    ;

args_list           : args_list COMMA expression    {
                                                      log.ParseLog.Printf("args_list0: %+v %+v %+v %+v\n", $1, $2, $3, yylex)
                                                    }
                    | expression                    {
                                                      log.ParseLog.Printf("args_list1: %+v %+v\n", $1, yylex)
                                                    }
                    ;

empty               : /* empty */                   {
                                                      log.ParseLog.Printf("empty0: %+v\n", yylex)
                                                    }
                    ;

%%
