package main

import (
	"fmt"
	"io"
	"os"
	"bufio/scanner"
)

func main() {
	l := bufio.
	l := NewLexer(os.Stdin)
	yyParse(l)
	fmt.Println("hello")
}
