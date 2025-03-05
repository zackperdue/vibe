package cmd

import (
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/vibe-lang/vibe/internal/ui"
	"github.com/vibe-lang/vibe/interpreter"
	"github.com/vibe-lang/vibe/lexer"
	"github.com/vibe-lang/vibe/parser"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run [file]",
	Short: "Run a Vibe program file",
	Long:  `Execute a Vibe program file with the .vi extension.`,
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			ui.ErrorColor.Println("Error: You must specify a file to run")
			return
		}

		filename := args[0]
		if !strings.HasSuffix(filename, ".vi") {
			filename = filename + ".vi"
		}

		source, err := os.ReadFile(filename)
		if err != nil {
			ui.ErrorColor.Printf("Error reading file: %s\n", err)
			return
		}

		runProgram(string(source), filename)
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}

func runProgram(source string, filename string) {
	// Create a lexer from the source code
	l := lexer.New(source)

	// Parse the input
	program, errors := parser.Parse(l)

	if len(errors) > 0 {
		ui.PrintParserErrors(errors)
		return
	}

	// Create an interpreter and evaluate the program
	interp := interpreter.New()

	// Evaluate the program with panic recovery
	result := func() (res interpreter.Value) {
		// Recover from panics to prevent program from crashing
		defer func() {
			if r := recover(); r != nil {
				ui.ErrorColor.Printf("Evaluation error: %v\n", r)
				res = nil
			}
		}()
		return interp.Eval(program)
	}()

	// The result is the last evaluated statement
	if result != nil && result.Type() != "NIL" {
		// Check if the result is an error value by seeing if it starts with "Type error:"
		resultStr := result.Inspect()
		if strings.HasPrefix(resultStr, "Type error:") ||
		   strings.HasPrefix(resultStr, "Error:") ||
		   result.Type() == "ERROR" {
			// For errors, show the error with file location
			ui.ErrorColor.Printf("Error in file %s: %s\n", filename, resultStr)
		} else {
			// For non-errors, print both the value and its type
			ui.ResultColor.Printf("Result: %s ", resultStr)
			ui.TypeColor.Printf(": %s\n", result.VibeType())
		}
	}
}