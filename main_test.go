package main

import (
	"github.com/ruistola/compiler-proto/lexer"
	"github.com/ruistola/compiler-proto/parser"
	"os"
	"path/filepath"
	"testing"
)

func TestParse(t *testing.T) {
	// Find all .jru files in the project directory
	files, err := filepath.Glob("examples/*.jru")
	if err != nil {
		t.Fatalf("Failed to find source files : %v", err)
	}

	if len(files) == 0 {
		t.Fatal("No .jru files found for testing")
	}

	// Test parsing each file
	for _, filename := range files {
		t.Run(filename, func(t *testing.T) {
			sourceBytes, err := os.ReadFile(filename)
			if err != nil {
				t.Fatalf("Failed to read file %s : %v", filename, err)
			}

			src := string(sourceBytes)

			// Tokenize
			tokens := lexer.Tokenize(src)

			// Parse tokens to AST
			defer func() {
				if r := recover(); r != nil {
					t.Fatalf("Parsing failed for %s : %v", filename, r)
				}
			}()

			parser.Parse(tokens)

			// If we got here without panic, the test passes
		})
	}
}
