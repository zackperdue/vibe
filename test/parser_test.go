package test

import (
	"testing"

	"github.com/example/vibe/ast"
	"github.com/example/vibe/lexer"
	"github.com/example/vibe/parser"
)

func TestTypeDeclaration(t *testing.T) {
	// Sample type declaration code - only the simple one without generics
	input := `type StringAlias = string`

	l := lexer.New(input)
	program, errors := parser.Parse(l)

	if len(errors) != 0 {
		t.Errorf("parser has %d errors", len(errors))
		for _, msg := range errors {
			t.Errorf("parser error: %q", msg)
		}
		t.FailNow()
	}

	// Debug: Print out all statements
	for i, stmt := range program.Statements {
		t.Logf("Statement %d: %T", i, stmt)
	}

	// Check that we only have one statement
	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain exactly 1 statement. got=%d",
			len(program.Statements))
	}

	// Check the type declaration
	stmt1, ok := program.Statements[0].(*ast.TypeDeclaration)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.TypeDeclaration. got=%T",
			program.Statements[0])
	}

	if stmt1.Name != "StringAlias" {
		t.Errorf("stmt1.Name not 'StringAlias'. got=%q", stmt1.Name)
	}

	typeValue1, ok := stmt1.TypeValue.(*ast.TypeAnnotation)
	if !ok {
		t.Fatalf("stmt1.TypeValue is not ast.TypeAnnotation. got=%T", stmt1.TypeValue)
	}

	if typeValue1.TypeName != "string" {
		t.Errorf("typeValue1.TypeName not 'string'. got=%q", typeValue1.TypeName)
	}
}

func TestGenericTypeDeclaration(t *testing.T) {
	// Sample generic type declaration code
	input := `type StringArray = Array<string>`

	l := lexer.New(input)
	program, errors := parser.Parse(l)

	if len(errors) != 0 {
		t.Errorf("parser has %d errors", len(errors))
		for _, msg := range errors {
			t.Errorf("parser error: %q", msg)
		}
		t.FailNow()
	}

	// Debug: Print out all statements
	for i, stmt := range program.Statements {
		t.Logf("Statement %d: %T", i, stmt)
	}

	// Check that we only have one statement
	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain exactly 1 statement. got=%d",
			len(program.Statements))
	}

	// Check the type declaration
	stmt, ok := program.Statements[0].(*ast.TypeDeclaration)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.TypeDeclaration. got=%T",
			program.Statements[0])
	}

	if stmt.Name != "StringArray" {
		t.Errorf("stmt.Name not 'StringArray'. got=%q", stmt.Name)
	}

	typeValue, ok := stmt.TypeValue.(*ast.TypeAnnotation)
	if !ok {
		t.Fatalf("stmt.TypeValue is not ast.TypeAnnotation. got=%T", stmt.TypeValue)
	}

	if typeValue.TypeName != "Array" {
		t.Errorf("typeValue.TypeName not 'Array'. got=%q", typeValue.TypeName)
	}

	if len(typeValue.TypeParams) != 1 {
		t.Fatalf("typeValue.TypeParams does not contain 1 type parameter. got=%d",
			len(typeValue.TypeParams))
	}

	typeParam, ok := typeValue.TypeParams[0].(*ast.TypeAnnotation)
	if !ok {
		t.Fatalf("typeValue.TypeParams[0] is not ast.TypeAnnotation. got=%T",
			typeValue.TypeParams[0])
	}

	if typeParam.TypeName != "string" {
		t.Errorf("typeParam.TypeName not 'string'. got=%q", typeParam.TypeName)
	}
}

func TestCombinedTypeDeclarations(t *testing.T) {
	// Use original syntax without spaces around angle brackets
	input := `type StringAlias = string
type StringArray = Array<string>`
	l := lexer.New(input)

	// Debug: Print all tokens
	t.Log("Tokens:")
	var token lexer.Token
	tokenCount := 0
	for {
		token = l.NextToken()
		t.Logf("Token %d: %s, Literal: %s, Line: %d, Column: %d",
			tokenCount, token.Type, token.Literal, token.Line, token.Column)
		tokenCount++
		if token.Type == lexer.EOF {
			break
		}
	}

	// Reset lexer for parsing
	l = lexer.New(input)
	program, errors := parser.Parse(l)

	if len(errors) != 0 {
		t.Errorf("parser has %d errors", len(errors))
		for _, msg := range errors {
			t.Errorf("parser error: %q", msg)
		}
		t.FailNow()
	}

	// Detailed output of statements
	t.Logf("Found %d statements", len(program.Statements))
	for i, stmt := range program.Statements {
		t.Logf("Statement %d: %T - %s", i, stmt, stmt.String())
	}

	if len(program.Statements) != 2 {
		t.Fatalf("program.Statements does not contain 2 statements. got=%d", len(program.Statements))
	}
}