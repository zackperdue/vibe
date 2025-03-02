package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/example/vibe/lexer"
)

func main() {
	// Check if a file was provided
	if len(os.Args) < 2 {
		fmt.Println("Usage: lexer_test <filename>")
		os.Exit(1)
	}

	// Read the file
	filename := os.Args[1]
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Printf("Error reading file: %s\n", err)
		os.Exit(1)
	}

	// Create a lexer
	l := lexer.New(string(data))

	// Print all tokens
	for {
		tok := l.NextToken()
		fmt.Printf("Token: %s, Literal: %s, Line: %d, Column: %d\n",
			tok.Type, tok.Literal, tok.Line, tok.Column)

		if tok.Type == "EOF" {
			break
		}
	}
}