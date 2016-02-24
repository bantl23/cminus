package main

import (
	"os"
)

func main() {
	yyParse(NewLexer(os.Stdin))
}
