package main

import (
	"fmt"
	"os"

	"github.com/vibe-lang/vibe/lexer"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Usage: %s <filename>\n", os.Args[0])
		os.Exit(1)
	}

	filename := os.Args[1]
	content, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("Error reading file %s: %v\n", filename, err)
		os.Exit(1)
	}

	l := lexer.New(string(content))

	fmt.Println("=== Token Stream ===")
	for tok := l.NextToken(); tok.Type != lexer.EOF; tok = l.NextToken() {
		fmt.Printf("Token: Type=%s, Literal=%q, Line=%d, Col=%d\n",
			tok.Type, tok.Literal, tok.Line, tok.Column)
	}
	fmt.Println("=== End Token Stream ===")
}