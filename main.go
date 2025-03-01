package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/example/crystal/interpreter"
	"github.com/example/crystal/lexer"
	"github.com/example/crystal/parser"
	"github.com/example/crystal/types"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: crystal <filename> or crystal -i (for interactive mode)")
		os.Exit(1)
	}

	if os.Args[1] == "-i" {
		runREPL()
	} else {
		runFile(os.Args[1])
	}
}

func runREPL() {
	fmt.Println("Crystal Programming Language REPL")
	fmt.Println("Type 'exit' to quit")
	fmt.Println("Type checking is enabled - type annotations supported")

	interp := interpreter.New()
	var input string

	for {
		fmt.Print(">> ")
		fmt.Scanln(&input)

		if input == "exit" {
			break
		}

		// Ensure the input has a newline at the end
		if !strings.HasSuffix(input, "\n") {
			input += "\n"
		}

		// Create lexer
		l := lexer.New(input)

		// Parse input
		program, errors := parser.Parse(l)

		if len(errors) > 0 {
			for _, err := range errors {
				fmt.Printf("Error: %s\n", err)
			}
			continue
		}

		// Evaluate program
		result := interp.Eval(program)

		// Only print the result if it's not nil
		if result.Type() != "NIL" {
			fmt.Printf("=> %s (type: %s)\n", result.Inspect(), formatType(result.CrystalType()))
		}
	}
}

func runFile(filename string) {
	// Read file
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Printf("Error reading file: %s\n", err)
		os.Exit(1)
	}

	// Create lexer
	l := lexer.New(string(data))

	// Parse input
	program, errors := parser.Parse(l)

	if len(errors) > 0 {
		fmt.Printf("Parser errors in %s:\n", filename)
		for _, err := range errors {
			fmt.Printf("  - %s\n", err)
		}
		os.Exit(1)
	}

	// Create interpreter and evaluate program
	interp := interpreter.New()
	result := interp.Eval(program)

	// If the result is an error, print it and exit
	if errValue, ok := result.(*interpreter.StringValue); ok {
		if strings.HasPrefix(errValue.Value, "Type error:") {
			fmt.Printf("Type error in %s: %s\n", filename, errValue.Value)
			os.Exit(1)
		}
	}
}

// Helper function to format Crystal types for display
func formatType(t types.Type) string {
	switch t := t.(type) {
	case types.BasicType:
		return string(t)
	case types.ArrayType:
		return fmt.Sprintf("Array<%s>", formatType(t.ElementType))
	case types.FunctionType:
		params := []string{}
		for _, p := range t.ParameterTypes {
			params = append(params, formatType(p))
		}
		return fmt.Sprintf("(%s) -> %s", strings.Join(params, ", "), formatType(t.ReturnType))
	case types.UnionType:
		types := []string{}
		for _, t := range t.Types {
			types = append(types, formatType(t))
		}
		return strings.Join(types, " | ")
	default:
		return "unknown"
	}
}