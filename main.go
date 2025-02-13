package main

import (
	"fmt"
	"jru-test/lexer"
	"os"
)

func main() {
	raw, err := os.ReadFile("main.jru")
	if err != nil {
		fmt.Printf("%v", err)
		os.Exit(1)
	}
	fmt.Printf("%s\n", lexer.Tokenize(string(raw)))
}
