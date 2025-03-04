package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/vibe-lang/vibe/interpreter"
	"github.com/vibe-lang/vibe/lexer"
	"github.com/vibe-lang/vibe/parser"
)

var debug bool = false

func main() {
	args := os.Args[1:]

	if len(args) == 0 {
		fmt.Println("Usage: vibe <filename> or vibe -i (for interactive mode)")
		fmt.Println("       vibe <filename> -d (for debug mode)")
		return
	}

	// Check for debug flag
	for i, arg := range args {
		if arg == "-d" || arg == "--debug" {
			debug = true
			// Remove the debug flag from args
			args = append(args[:i], args[i+1:]...)
			break
		}
	}

	if len(args) == 0 {
		fmt.Println("Usage: vibe <filename> or vibe -i (for interactive mode)")
		fmt.Println("       vibe <filename> -d (for debug mode)")
		return
	}

	if args[0] == "-i" {
		runInteractiveMode()
		return
	}

	filename := args[0]
	if !strings.HasSuffix(filename, ".vi") {
		filename = filename + ".vi"
	}

	source, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Printf("Error reading file: %s\n", err)
		return
	}

	runProgram(string(source))
}

func runInteractiveMode() {
	fmt.Println("Vibe Programming Language REPL")
	fmt.Println("Type 'exit' to quit")

	interp := interpreter.New()
	scanner := bufio.NewScanner(os.Stdin)

	for {
		// Track input buffer for multiline input
		var inputBuffer strings.Builder
		var isMultiline bool
		var blockCount int = 0

		fmt.Print(">> ")
		if !scanner.Scan() {
			break
		}

		line := scanner.Text()
		if line == "exit" {
			break
		}

		// Check if this line might start a multi-line block
		if containsBlockOpener(line) {
			isMultiline = true
			blockCount++
		}

		inputBuffer.WriteString(line)

		// If we're in a multiline context, keep collecting lines until the blocks are closed
		for isMultiline && blockCount > 0 {
			fmt.Print(".. ")
			if !scanner.Scan() {
				break
			}

			line = scanner.Text()
			if line == "exit" {
				break
			}

			// Track block openers and closers
			if containsBlockOpener(line) {
				blockCount++
			}
			if containsBlockCloser(line) {
				blockCount--
			}

			// Add the line to our buffer
			inputBuffer.WriteString("\n")
			inputBuffer.WriteString(line)
		}

		code := inputBuffer.String()

		// Create a lexer from the input
		l := lexer.New(code)

		// Parse the input
		program, errors := parser.Parse(l)

		if len(errors) > 0 {
			printParserErrors(errors)
			continue
		}

		// Evaluate the program
		result := interp.Eval(program)
		if result != nil {
			fmt.Printf("=> %s : %s\n", result.Inspect(), result.VibeType())
		}
	}
}

// Helper function to detect if a line contains a block opener token
func containsBlockOpener(line string) bool {
	// Check for block openers: for, if, function, class, while, etc.
	keywords := []string{"for", "if", "function", "class", "while", "do"}
	for _, keyword := range keywords {
		if strings.Contains(line, keyword+" ") || strings.HasSuffix(line, keyword) ||
		   strings.Contains(line, keyword+"\n") {
			return true
		}
	}
	// Also check for array literals that span multiple lines
	if strings.Contains(line, "[") && !strings.Contains(line, "]") {
		return true
	}
	return false
}

// Helper function to detect if a line contains a block closer token
func containsBlockCloser(line string) bool {
	// Check for 'end' keyword which closes most blocks
	if strings.TrimSpace(line) == "end" || strings.HasPrefix(strings.TrimSpace(line), "end ") ||
	   strings.HasSuffix(strings.TrimSpace(line), " end") || strings.Contains(line, " end ") {
		return true
	}
	// Also check for closing brackets for array literals
	if strings.Contains(line, "]") {
		return true
	}
	return false
}

func runProgram(source string) {
	// Create a lexer from the source code
	l := lexer.New(source)

	// Parse the input
	program, errors := parser.Parse(l)

	if len(errors) > 0 {
		printParserErrors(errors)
		return
	}

	if debug {
		fmt.Println("Program AST:")
		for i, stmt := range program.Statements {
			fmt.Printf("Statement %d: %s\n", i, stmt.String())
		}
		fmt.Println()
	}

	// Create an interpreter and evaluate the program
	interp := interpreter.New()
	result := interp.Eval(program)

	// The result is the last evaluated statement
	if result != nil && result.Type() != "NIL" {
		fmt.Printf("Result: %s : %s\n", result.Inspect(), result.VibeType())
	}
}

func printParserErrors(errors []string) {
	fmt.Println("Parser errors:")
	for _, err := range errors {
		fmt.Printf("\t%s\n", err)
	}
}