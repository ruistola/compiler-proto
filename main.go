package main

import (
	"fmt"
	"github.com/yassinebenaid/godump"
	"jru-test/lexer"
	"jru-test/parser"
	"jru-test/typechecker"
	"os"
	"time"
)

func main() {
	filename := "examples/arrayTypes.jru"
	sourceBytes, _ := os.ReadFile(filename)
	src := string(sourceBytes)

	fmt.Printf("Raw source (%s):\n--\n%s--\n", filename, src)

	totalDuration := time.Duration(0)

	startTokenization := time.Now()
	tokens := lexer.Tokenize(src)
	durationTokenization := time.Since(startTokenization)
	totalDuration += durationTokenization
	fmt.Printf("Tokenized %s in %v.\n\n", filename, durationTokenization)
	fmt.Printf("Tokens:\n%s\n\n", tokens)

	startParsing := time.Now()
	ast := parser.Parse(tokens)
	durationParsing := time.Since(startParsing)
	totalDuration += durationParsing
	fmt.Printf("Parsed %s in %v.\n\n", filename, durationParsing)

	fmt.Println("Parsed AST:")
	godump.Dump(ast)

	startTypeChecking := time.Now()
	errors := typechecker.Check(ast)
	durationTypeChecking := time.Since(startTypeChecking)
	totalDuration += durationTypeChecking
	if len(errors) == 0 {
		fmt.Println("0 errors.")
	} else {
		for _, err := range errors {
			fmt.Println(err)
		}
	}
	fmt.Printf("Type checked %s in %v.\n\n", filename, durationTypeChecking)

	fmt.Printf("Done in %v.\n", totalDuration)
}
