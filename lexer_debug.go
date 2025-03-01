package main

import (
	"fmt"
	"os"

	"github.com/example/vibe/lexer"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run lexer_debug.go \"code to tokenize\"")
		os.Exit(1)
	}

	input := os.Args[1]
	l := lexer.New(input)

	fmt.Println("Tokens:")
	for {
		tok := l.NextToken()
		fmt.Printf("Type: %s, Literal: %s\n", tok.Type, tok.Literal)
		if tok.Type == lexer.EOF {
			break
		}
	}
}