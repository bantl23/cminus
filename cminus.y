%{
package main

import (
  "fmt"
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
                                                      fmt.Printf("program0: %+v %v\n", $1, currStr(yylex))
                                                    }
                    ;

declaration_list    : declaration_list declaration  {
                                                      fmt.Printf("declaration_list0: %+v %+v %v\n", $1, $2, currStr(yylex))
                                                    }
                    | declaration                   {
                                                      fmt.Printf("declaration_list1: %+v %v\n", $1, currStr(yylex))
                                                    }
                    ;

declaration         : var_declaration               {
                                                      fmt.Printf("declaration0: %+v %v\n", $1, currStr(yylex))
                                                    }
                    | fun_declaration               {
                                                      fmt.Printf("declaration1: %+v %v\n", $1, currStr(yylex))
                                                    }
                    ;

var_declaration     : type_specifier ID SEMI        {
                                                      fmt.Printf("var_declaration0: %+v %+v %+v %v\n", $1, $2, $3, currStr(yylex))
                                                    }
                    | type_specifier ID LBRACKET NUM RBRACKET SEMI
                                                    {
                                                      fmt.Printf("var_declaration1: %+v %+v %+v %+v %+v %+v %+v\n", $1, $2, $3, $4, $5, $6, currStr(yylex))
                                                    }
                    ;

type_specifier      : INT                           {
                                                      fmt.Printf("type_specifier0: %+v %+v\n", $1, currStr(yylex))
                                                    }
                    | VOID                          {
                                                      fmt.Printf("type_specifier1: %+v %+v\n", $1, currStr(yylex))
                                                    }
                    ;

fun_declaration     : type_specifier ID LPAREN params RPAREN compound_stmt
                                                    {
                                                      fmt.Printf("fun_declaration0: %+v %+v %+v %+v %+v %+v %+v\n", $1, $2, $3, $4, $5, $6, currStr(yylex))
                                                    }
                    ;

params              : param_list                    {
                                                      fmt.Printf("params0: %+v %+v\n", $1, currStr(yylex))
																										}
                    | VOID                          {
                                                      fmt.Printf("params1: %+v %+v\n", $1, currStr(yylex))
																										}
                    ;

param_list          : param_list COMMA param        {
                                                      fmt.Printf("param_list0: %+v %+v %+v %+v\n", $1, $2, $3, currStr(yylex))
																										}
                    | param                         {
                                                      fmt.Printf("param_list1: %+v %+v\n", $1, currStr(yylex))
																										}
                    ;

param               : type_specifier ID             {
                                                      fmt.Printf("param0: %+v %+v %+v\n", $1, $2, currStr(yylex))
																										}
                    | type_specifier ID LBRACKET RBRACKET
                                                    {
                                                      fmt.Printf("param1: %+v %+v %+v %+v %+v\n", $1, $2, $3, $3, currStr(yylex))
																										}
                    ;

compound_stmt       : LBRACE local_declarations statement_list RBRACE
                                                    {
                                                      fmt.Printf("compound_stmt0: %+v %+v %+v %+v %+v\n", $1, $2, $3, $4, currStr(yylex))
																										}
                    ;

local_declarations  : local_declarations var_declaration
                                                    {
                                                      fmt.Printf("local_declarations0: %+v %+v %+v\n", $1, $2, currStr(yylex))
																										}
                    | empty                         {
                                                      fmt.Printf("local_declarations1: %+v %+v\n", $1, currStr(yylex))
																										}
                    ;

statement_list      : statement_list statement      {
                                                      fmt.Printf("statement_list0: %+v %+v %+v\n", $1, $2, currStr(yylex))
																										}
                    | empty                         {
                                                      fmt.Printf("statement_list1: %+v %+v\n", $1, currStr(yylex))
                                                    }
                    ;

statement           : expression_stmt               {
                                                      fmt.Printf("statement0: %+v %+v", $1, currStr(yylex))
																										}
                    | compound_stmt                 {
                                                      fmt.Printf("statement1: %+v %+v", $1, currStr(yylex))
																										}
                    | selection_stmt                {
                                                      fmt.Printf("statement2: %+v %+v", $1, currStr(yylex))
																										}
                    | iteration_stmt                {
                                                      fmt.Printf("statement3: %+v %+v", $1, currStr(yylex))
																										}
                    | return_stmt                   {
                                                      fmt.Printf("statement4: %+v %+v", $1, currStr(yylex))
																										}
                    ;

expression_stmt     : expression SEMI               {
                                                      fmt.Printf("expression_stmt0: %+v %+v %+v", $1, $2, currStr(yylex))
																										}
                    | SEMI                          {
                                                      fmt.Printf("expression_stmt1: %+v %+v", $1, currStr(yylex))
																										}
                    ;

selection_stmt      : IF LPAREN expression RPAREN statement %prec THEN
                                                    {
                                                      fmt.Printf("selection_stmt0: %+v %+v %+v %+v %+v %+v\n", $1, $2, $3, $4, $5, currStr(yylex))
																										}
                    | IF LPAREN expression RPAREN statement ELSE statement
                                                    {
                                                      fmt.Printf("selection_stmt1: %+v %+v %+v %+v %+v %+v %+v %+v\n", $1, $2, $3, $4, $5, $6, $7, currStr(yylex))
																										}
                    ;

iteration_stmt      : WHILE LPAREN expression RPAREN statement
                                                    {
                                                      fmt.Printf("iteration_stmt0: %+v %+v %+v %+v %+v %+v\n", $1, $2, $3, $4, $5, currStr(yylex))
																										}
                    ;

return_stmt         : RETURN SEMI                   {
                                                      fmt.Printf("return_stmt0: %+v %+v %+v\n", $1, $2, currStr(yylex))
																										}
                    | RETURN expression SEMI        {
                                                      fmt.Printf("return_stmt1: %+v %+v %+v %+v\n", $1, $2, $3, currStr(yylex))
                                                    }
                    ;

expression          : var ASSIGN expression         {
                                                      fmt.Printf("expression0: %+v %+v %+v %+v\n", $1, $2, $3, currStr(yylex))
                                                    }
                    | simple_expression             {
                                                      fmt.Printf("expression1: %+v %+v\n", $1, currStr(yylex))
                                                    }
                    ;

var                 : ID                            {
                                                      fmt.Printf("var0: %+v %+v\n", $1, currStr(yylex))
                                                    }
                    | ID LBRACKET expression RBRACKET
                                                    {
                                                      fmt.Printf("var1: %+v %+v %+v %+v %+v\n", $1, $2, $3, $4, currStr(yylex))
                                                    }
                    ;

simple_expression   : additive_expression relop additive_expression
                                                    {
                                                      fmt.Printf("simple_expression0: %+v %+v %+v %+v\n", $1, $2, $3, currStr(yylex))
                                                    }
                    | additive_expression           {
                                                      fmt.Printf("simple_expression1: %+v %+v\n", $1, currStr(yylex))
                                                    }
                    ;

relop               : LTE                           {
                                                      fmt.Printf("relop0: %+v %+v\n", $1, currStr(yylex))
                                                    }
                    | LT                            {
                                                      fmt.Printf("relop1: %+v %+v\n", $1, currStr(yylex))
                                                    }
                    | GT                            {
                                                      fmt.Printf("relop2: %+v %+v\n", $1, currStr(yylex))
                                                    }
                    | GTE                           {
                                                      fmt.Printf("relop3: %+v %+v\n", $1, currStr(yylex))
                                                    }
                    | EQ                            {
                                                      fmt.Printf("relop4: %+v %+v\nv", $1, currStr(yylex))
                                                    }
                    | NEQ                           {
                                                      fmt.Printf("relop5: %+v %+v\n", $1, currStr(yylex))
                                                    }
                    ;

additive_expression : additive_expression addop term
                                                    {
																											fmt.Printf("additive_expression0: %+v %+v %+v %+v\n", $1, $2, $3, currStr(yylex))
																										}
                    | term                          {
																											fmt.Printf("additive_expression1: %+v %+v %+v %+v\n", $1, currStr(yylex))
																										}
                    ;

addop               : PLUS                          {
																											fmt.Printf("addop0: %+v %+v\n", $1, currStr(yylex))
																										}
                    | MINUS                         {
																											fmt.Printf("addop1: %+v %+v\n", $1, currStr(yylex))
																										}
                    ;

term                : term mulop factor             {
																											fmt.Printf("term0: %+v %+v %+v %+v\n", $1, $2, $3, currStr(yylex))
																										}
                    | factor                        {
																											fmt.Printf("term1: %+v %+v\n", $1, currStr(yylex))
																										}
                    ;

mulop               : TIMES                         {
																											fmt.Printf("mulop0: %+v %+v\n", $1, currStr(yylex))
																										}
                    | OVER                          {
																											fmt.Printf("mulop1: %+v %+v\n", $1, currStr(yylex))
																										}
                    ;

factor              : LPAREN expression RPAREN      {
																											fmt.Printf("factor0: %+v %+v %+v %+v\n", $1, $2, $3, currStr(yylex))
																										}
                    | var                           {
																											fmt.Printf("factor1: %+v %+v\n", $1, currStr(yylex))
																										}
                    | call                          {
																											fmt.Printf("factor2: %+v %+v\n", $1, currStr(yylex))
																										}
                    | NUM                           {
																											fmt.Printf("factor3: %+v %+v\n", $1, currStr(yylex))
																										}
                    ;

call                : ID LPAREN args RPAREN         {
																											fmt.Printf("call0: %+v %+v %+v %+v %+v\n", $1, $2, $3, $4, currStr(yylex))
																										}
                    ;

args                : args_list                     {
																											fmt.Printf("args0: %+v %+v\n", $1, currStr(yylex))
																										}
                    | empty                         {
																											fmt.Printf("args1: %+v %+v\n", $1, currStr(yylex))
																										}
                    ;

args_list           : args_list COMMA expression    {
                                                      fmt.Printf("args_list0: %+v %+v %+v %+v\n", $1, $2, $3, currStr(yylex))
                                                    }
                    | expression                    {
                                                      fmt.Printf("args_list1: %+v %+v\n", $1, currStr(yylex))
                                                    }
                    ;

empty               : /* empty */                   {
                                                      fmt.Printf("empty0: %+v\n", currStr(yylex))
                                                    }
                    ;

%%
func currText(y yyLexer) string {
  if len(y.(*Lexer).stack) > 0 {
    return y.(*Lexer).stack[0].s
  }
  return ""
}
func currLine(y yyLexer) int {
  if len(y.(*Lexer).stack) > 0 {
    return y.(*Lexer).stack[0].line
  }
  return 0
}
func currCol(y yyLexer) int {
  if len(y.(*Lexer).stack) > 0 {
    return y.(*Lexer).stack[0].column
  }
  return 0
}
func currStr(y yyLexer) string {
  return fmt.Sprintf("%+v [%+v:%+v]", currText(y), currLine(y), currCol(y))
}
