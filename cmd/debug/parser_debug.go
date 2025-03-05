package debug

import (
	"fmt"
	"os"

	"github.com/vibe-lang/vibe/lexer"
	"github.com/vibe-lang/vibe/parser"
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
	p := parser.New(l)

	program, errors := parser.Parse(l)

	fmt.Println("=== Parse Results ===")
	if len(errors) > 0 {
		fmt.Println("Parser Errors:")
		for _, err := range errors {
			fmt.Printf("  - %s\n", err)
		}
	} else {
		fmt.Println("Parsing successful!")
	}

	fmt.Printf("\nNumber of statements: %d\n", len(program.Statements))
	for i, stmt := range program.Statements {
		fmt.Printf("Statement %d: Type=%T, Value=%s\n", i, stmt, stmt.String())
	}
	fmt.Println("=== End Parse Results ===")
}