package parser_test

import (
	"testing"

	"github.com/vibe-lang/vibe/ast"
	"github.com/vibe-lang/vibe/lexer"
	"github.com/vibe-lang/vibe/parser"
)

func TestSimpleForLoop(t *testing.T) {
	input := `for x in [1] do end`

	l := lexer.New(input)
	program, errors := parser.Parse(l)

	if len(errors) > 0 {
		t.Fatalf("Parser encountered errors: %v", errors)
	}

	if len(program.Statements) != 1 {
		t.Fatalf("Program does not contain 1 statement. got=%d", len(program.Statements))
	}

	// Check the statement is a for loop
	forStmt, ok := program.Statements[0].(*ast.ForStmt)
	if !ok {
		t.Fatalf("Statement is not a ForStmt. got=%T", program.Statements[0])
	}

	if forStmt.Iterator != "x" {
		t.Errorf("Iterator is not 'x'. got=%s", forStmt.Iterator)
	}

	// Check that the iterable is an array literal
	arrayLiteral, ok := forStmt.Iterable.(*ast.ArrayLiteral)
	if !ok {
		t.Fatalf("Iterable is not an ArrayLiteral. got=%T", forStmt.Iterable)
	}

	if len(arrayLiteral.Elements) != 1 {
		t.Errorf("Array does not have 1 element. got=%d", len(arrayLiteral.Elements))
	}

	// Test that the body is empty
	if len(forStmt.Body.Statements) != 0 {
		t.Fatalf("Body is not empty. got=%d statements", len(forStmt.Body.Statements))
	}
}