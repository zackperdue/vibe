package cmd

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/chzyer/readline"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/vibe-lang/vibe/internal/ui"
	"github.com/vibe-lang/vibe/interpreter"
	"github.com/vibe-lang/vibe/lexer"
	"github.com/vibe-lang/vibe/parser"
)

// interactiveCmd represents the interactive command
var interactiveCmd = &cobra.Command{
	Use:     "interactive",
	Aliases: []string{"i"},
	Short:   "Start an interactive REPL session",
	Long:    `Start an interactive Read-Eval-Print Loop (REPL) session for the Vibe programming language.`,
	Run: func(cmd *cobra.Command, args []string) {
		runInteractiveMode()
	},
}

func init() {
	rootCmd.AddCommand(interactiveCmd)
}

func runInteractiveMode() {
	ui.HeadingColor.Println("Vibe Programming Language REPL")
	fmt.Println("Type 'exit' to quit")

	// Setup readline with history file
	home, err := os.UserHomeDir()
	if err != nil {
		ui.ErrorColor.Printf("Error getting home directory: %s\n", err)
		return
	}

	historyDir := filepath.Join(home, ".vibe")
	historyFile := filepath.Join(historyDir, "history")

	// Ensure history directory exists
	if _, err := os.Stat(historyDir); os.IsNotExist(err) {
		if err := os.MkdirAll(historyDir, 0755); err != nil {
			ui.ErrorColor.Printf("Error creating history directory: %s\n", err)
			// Continue without persistence if we can't create the directory
		}
	}

	// The actual prompt shown by readline
	mainPrompt := ui.PromptColor.Sprint(">> ")
	continuationPrompt := ui.ContinuationColor.Sprint(".. ")

	// Configure readline
	rl, err := readline.NewEx(&readline.Config{
		Prompt:                 mainPrompt,
		HistoryFile:            historyFile,
		HistoryLimit:           viper.GetInt("history.maxSize"),
		InterruptPrompt:        "^C",
		EOFPrompt:              "exit",
		DisableAutoSaveHistory: true, // Disable automatic history saving
	})
	if err != nil {
		ui.ErrorColor.Printf("Error setting up readline: %s\n", err)
		return
	}
	defer rl.Close()

	interp := interpreter.New()

	for {
		// Track input buffer for multiline input
		var inputBuffer strings.Builder
		var isMultiline bool
		var blockCount int = 0

		// Get the first line
		line, err := rl.Readline()
		if err == readline.ErrInterrupt {
			// If the user hits Ctrl+C, we'll start a new prompt
			continue
		} else if err == io.EOF {
			// If the user hits Ctrl+D, we'll exit
			break
		}

		// Check if we should exit
		if line == "exit" {
			break
		}

		// Manually save to history if it's not empty and not 'exit'
		if strings.TrimSpace(line) != "" && line != "exit" {
			rl.SaveHistory(line)
		}

		// Check if this line might start a multi-line block
		if containsBlockOpener(line) {
			isMultiline = true
			blockCount++
		}

		inputBuffer.WriteString(line)

		// If we're in a multiline context, keep collecting lines until the blocks are closed
		for isMultiline && blockCount > 0 {
			// Change the prompt for continuation lines
			rl.SetPrompt(continuationPrompt)

			line, err := rl.Readline()
			if err != nil {
				break
			}

			if line == "exit" {
				break
			}

			// Manually save to history if it's not empty and not 'exit'
			if strings.TrimSpace(line) != "" && line != "exit" {
				rl.SaveHistory(line)
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

		// Reset the prompt for the next input
		rl.SetPrompt(mainPrompt)

		code := inputBuffer.String()

		// Skip empty input
		if strings.TrimSpace(code) == "" {
			continue
		}

		// Create a lexer from the input
		l := lexer.New(code)

		// Parse the input
		program, errors := parser.Parse(l)

		if len(errors) > 0 {
			ui.PrintParserErrors(errors)
			continue
		}

		// Evaluate the program with panic recovery
		result := func() (res interpreter.Value) {
			// Recover from panics to prevent REPL from crashing
			defer func() {
				if r := recover(); r != nil {
					ui.ErrorColor.Printf("Evaluation error: %v\n", r)
					res = nil
				}
			}()
			return interp.Eval(program)
		}()

		if result != nil {
			// Check if the result is an error value by seeing if it starts with "Type error:"
			resultStr := result.Inspect()
			if strings.HasPrefix(resultStr, "Type error:") || result.Type() == "ERROR" {
				// For errors, just print the error message without the type
				ui.ErrorColor.Printf("=> %s\n", resultStr)
			} else {
				// For non-errors, print both the value and its type
				ui.ResultColor.Printf("=> %s ", resultStr)
				ui.TypeColor.Printf(": %s\n", result.VibeType())
			}
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