// Package ui provides user interface utilities for the Vibe language
package ui

import (
	"github.com/fatih/color"
)

// Define colors for different UI elements
var (
	// PromptColor is used for the main REPL prompt
	PromptColor = color.New(color.FgCyan, color.Bold)

	// ContinuationColor is used for the continuation prompt in multi-line input
	ContinuationColor = color.New(color.FgCyan)

	// ResultColor is used for the result value in the REPL
	ResultColor = color.New(color.FgGreen)

	// TypeColor is used for type information in the REPL
	TypeColor = color.New(color.FgMagenta)

	// ErrorColor is used for error messages
	ErrorColor = color.New(color.FgRed)

	// WarningColor is used for warning messages
	WarningColor = color.New(color.FgYellow)

	// HeadingColor is used for headings and titles
	HeadingColor = color.New(color.FgHiBlue, color.Bold)
)

// PrintError prints an error message in the error color
func PrintError(format string, a ...interface{}) {
	ErrorColor.Printf(format, a...)
}

// PrintParserErrors prints parser errors in the error color
func PrintParserErrors(errors []string) {
	ErrorColor.Println("Parser errors:")
	for _, err := range errors {
		ErrorColor.Printf("\t%s\n", err)
	}
}