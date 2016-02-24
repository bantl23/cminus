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
                                                      fmt.Printf("%+v", $1)
                                                    }
                    ;

declaration_list    : declaration_list declaration  {
                                                      fmt.Printf("%+v %+v", $1, $2)
                                                    }
                    | declaration                   {
                                                      fmt.Printf("%+v", $1)
                                                    }
                    ;

declaration         : var_declaration               {
                                                      fmt.Printf("%+v", $1)
                                                    }
                    | fun_declaration               {
                                                      fmt.Printf("%+v", $1)
                                                    }
                    ;

var_declaration     : type_specifier ID SEMI        {
                                                      fmt.Printf("%+v %+v %+v", $1, $2, $3)
                                                    }
                    | type_specifier ID LBRACKET NUM RBRACKET SEMI
                                                    {
                                                      fmt.Printf("%+v %+v %+v %+v %+v %+v", $1, $2, $3, $4, $5, $6)
                                                    }
                    ;

type_specifier      : INT                           {
                                                      fmt.Printf("%+v", $1)
                                                    }
                    | VOID                          {
                                                      fmt.Printf("%+v", $1)
                                                    }
                    ;

fun_declaration     : type_specifier ID LPAREN params RPAREN compound_stmt
                                                    {
                                                      fmt.Printf("%+v %+v %+v %+v %+v %+v", $1, $2, $3, $4, $5, $6)
                                                    }
                    ;

params              : param_list                    {
																										}
                    | VOID                          {
																										}
                    ;

param_list          : param_list COMMA param        {
																										}
                    | param                         {
																										}
                    ;

param               : type_specifier ID             {
																										}
                    | type_specifier ID LBRACKET RBRACKET
                                                    {
																										}
                    ;

compound_stmt       : LBRACE local_declarations statement_list RBRACE
                                                    {
																										}
                    ;

local_declarations  : local_declarations var_declaration
                                                    {
																										}
                    | empty                         {
																										}
                    ;

statement_list      : statement_list statement      {
																										}
                    | empty
                    ;

statement           : expression_stmt               {
																										}
                    | compound_stmt                 {
																										}
                    | selection_stmt                {
																										}
                    | iteration_stmt                {
																										}
                    | return_stmt                   {
																										}
                    ;

expression_stmt     : expression SEMI               {
																										}
                    | SEMI                          {
																										}
                    ;

selection_stmt      : IF LPAREN expression RPAREN statement %prec THEN
                                                    {
																										}
                    | IF LPAREN expression RPAREN statement ELSE statement
                                                    {
																										}
                    ;

iteration_stmt      : WHILE LPAREN expression RPAREN statement
                                                    {
																										}
                    ;

return_stmt         : RETURN SEMI                   {
																										}
                    | RETURN expression SEMI        {}
                    ;

expression          : var ASSIGN expression         {}
                    | simple_expression             {}
                    ;

var                 : ID                            {}
                    | ID LBRACKET expression RBRACKET
                                                    {}
                    ;

simple_expression   : additive_expression relop additive_expression
                                                    {}
                    | additive_expression           {}
                    ;

relop               : LTE                           {}
                    | LT                            {}
                    | GT                            {}
                    | GTE                           {}
                    | EQ                            {}
                    | NEQ                           {}
                    ;

additive_expression : additive_expression addop term
                                                    {
																											fmt.Printf("%+v %+v %+v", $1, $2, $3)
																										}
                    | term                          {
																											fmt.Printf("%+v %+v %+v", $1, $2, $3)
																										}
                    ;

addop               : PLUS                          {
																											fmt.Printf("%+v", $1)
																										}
                    | MINUS                         {
																											fmt.Printf("%+v", $1)
																										}
                    ;

term                : term mulop factor             {
																											fmt.Printf("%+v %+v %+v", $1, $2, $3)
																										}
                    | factor                        {
																											fmt.Printf("%+v", $1)
																										}
                    ;

mulop               : TIMES                         {
																											fmt.Printf("%+v", $1)
																										}
                    | OVER                          {
																											fmt.Printf("%+v", $1)
																										}
                    ;

factor              : LPAREN expression RPAREN      {
																											fmt.Printf("%+v %+v %+v", $1, $2, $3)
																										}
                    | var                           {
																											fmt.Printf("%+v", $1)
																										}
                    | call                          {
																											fmt.Printf("%+v", $1)
																										}
                    | NUM                           {
																											fmt.Printf("%+v", $1)
																										}
                    ;

call                : ID LPAREN args RPAREN         {
																											fmt.Printf("%+v %+v %+v %+v", $1, $2, $3, $4)
																										}
                    ;

args                : args_list                     {
																											fmt.Printf("%+v", $1)
																										}
                    | empty                         {
																											fmt.Printf("%+v", $1)
																										}
                    ;

args_list           : args_list COMMA expression    {
                                                      fmt.Printf("%+v %+v %+v", $1, $2, $3)
                                                    }
                    | expression                    {
                                                      fmt.Printf("%+v", $1)
                                                    }
                    ;

empty               : /* empty */                   {
                                                      fmt.Printf("%+v", $1)
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
