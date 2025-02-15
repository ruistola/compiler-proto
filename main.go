package main

import (
	"fmt"
	"os"
)

func main() {
	raw, err := os.ReadFile("main.jru")
	if err != nil {
		fmt.Printf("%v", err)
		os.Exit(1)
	}
	fmt.Printf("%s\n", tokenize(string(raw)))
}
