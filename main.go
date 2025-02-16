package main

import (
	"fmt"
)

func main() {
	src := "2 + 3 * 4 - 1"
	fmt.Printf("Raw source: %s\n", src)
	ast := parse(src)
	fmt.Printf("Parsed into AST: %s\n", ast)
}
