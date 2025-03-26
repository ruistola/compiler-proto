package main

import (
	"fmt"
	"jru-test/lexer"
	"jru-test/parser"
	"os"
	"time"
)

func main() {
	filename := "source1.jru"
	sourceBytes, _ := os.ReadFile(filename)
	src := string(sourceBytes)

	start := time.Now()
	tokens := lexer.Tokenize(src)
	ast := parser.Parse(tokens)
	duration := time.Since(start)
	fmt.Printf("Parsed %s in %v\n\n", filename, duration)

	fmt.Printf("Raw source:\n%s\n", src)

	fmt.Printf("Tokens:\n%s\n\n", tokens)

	fmt.Printf("Parsed AST:\n%s\n", ast)
}
