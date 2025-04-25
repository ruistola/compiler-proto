package main

import (
	"fmt"
	"github.com/yassinebenaid/godump"
	"jru-test/lexer"
	"jru-test/parser"
	"os"
	"time"
)

func main() {
	filename := "examples/ifStatement.jru"
	sourceBytes, _ := os.ReadFile(filename)
	src := string(sourceBytes)

	fmt.Printf("Raw source (%s):\n--\n%s--\n", filename, src)

	startTokenization := time.Now()
	tokens := lexer.Tokenize(src)
	durationTokenization := time.Since(startTokenization)
	fmt.Printf("Tokenized %s in %v.\n\n", filename, durationTokenization)
	fmt.Printf("Tokens:\n%s\n\n", tokens)

	startParsing := time.Now()
	ast := parser.Parse(tokens)
	durationParsing := time.Since(startParsing)
	fmt.Printf("Parsed %s in %v.\n\n", filename, durationParsing)

	fmt.Println("Parsed AST:")
	godump.Dump(ast)
	fmt.Printf("Done in %v.\n", durationTokenization+durationParsing)
}
