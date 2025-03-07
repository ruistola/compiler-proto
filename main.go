package main

import (
	"fmt"
	"jru-test/parser"
)

func main() {
	src := "2 + 3 * 4 - -1"
	fmt.Printf("Raw source: %s\n", src)
	ast := parser.Parse(src)
	fmt.Println("Parsed AST:")
	fmt.Printf("%s", ast)
}
