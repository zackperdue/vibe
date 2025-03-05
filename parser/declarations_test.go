package parser

import (
	"testing"

	"github.com/vibe-lang/vibe/ast"
	"github.com/vibe-lang/vibe/lexer"
)

// TestTypeDeclaration tests parsing of simple type declarations
func TestTypeDeclaration(t *testing.T) {
	input := `type StringAlias = string`

	l := lexer.New(input)
	p := New(l)
	program, err := p.Parse()

	if err != nil {
		t.Fatalf("parser error: %v", err)
	}

	if len(program.Statements) != 1 {
		t.Fatalf("program has not enough statements. got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.TypeDeclaration)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.TypeDeclaration. got=%T", program.Statements[0])
	}

	if stmt.Name != "StringAlias" {
		t.Errorf("stmt.Name not 'StringAlias'. got=%q", stmt.Name)
	}

	typeAnnotation, ok := stmt.TypeValue.(*ast.TypeAnnotation)
	if !ok {
		t.Fatalf("stmt.TypeValue is not ast.TypeAnnotation. got=%T", stmt.TypeValue)
	}

	if typeAnnotation.TypeName != "string" {
		t.Errorf("typeAnnotation.TypeName not 'string'. got=%q", typeAnnotation.TypeName)
	}
}

// TestGenericTypeDeclaration tests parsing of type declarations with generic types
func TestGenericTypeDeclaration(t *testing.T) {
	input := `type StringArray = Array<string>`

	l := lexer.New(input)
	p := New(l)
	program, err := p.Parse()

	if err != nil {
		t.Fatalf("parser error: %v", err)
	}

	if len(program.Statements) != 1 {
		t.Fatalf("program has not enough statements. got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.TypeDeclaration)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.TypeDeclaration. got=%T", program.Statements[0])
	}

	if stmt.Name != "StringArray" {
		t.Errorf("stmt.Name not 'StringArray'. got=%q", stmt.Name)
	}

	typeAnnotation, ok := stmt.TypeValue.(*ast.TypeAnnotation)
	if !ok {
		t.Fatalf("stmt.TypeValue is not ast.TypeAnnotation. got=%T", stmt.TypeValue)
	}

	if typeAnnotation.TypeName != "Array" {
		t.Errorf("typeAnnotation.TypeName not 'Array'. got=%q", typeAnnotation.TypeName)
	}

	if len(typeAnnotation.TypeParams) != 1 {
		t.Fatalf("length of typeAnnotation.TypeParams wrong. got=%d", len(typeAnnotation.TypeParams))
	}

	param, ok := typeAnnotation.TypeParams[0].(*ast.TypeAnnotation)
	if !ok {
		t.Fatalf("typeAnnotation.TypeParams[0] is not ast.TypeAnnotation. got=%T", typeAnnotation.TypeParams[0])
	}

	if param.TypeName != "string" {
		t.Errorf("param.TypeName not 'string'. got=%q", param.TypeName)
	}
}

// TestMultipleGenericTypeParameters tests parsing of types with multiple generic parameters
func TestMultipleGenericTypeParameters(t *testing.T) {
	input := `type Dictionary = Map<string, number>`

	l := lexer.New(input)
	p := New(l)
	program, err := p.Parse()

	if err != nil {
		t.Fatalf("parser error: %v", err)
	}

	stmt, ok := program.Statements[0].(*ast.TypeDeclaration)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.TypeDeclaration. got=%T", program.Statements[0])
	}

	typeAnnotation, ok := stmt.TypeValue.(*ast.TypeAnnotation)
	if !ok {
		t.Fatalf("stmt.TypeValue is not ast.TypeAnnotation. got=%T", stmt.TypeValue)
	}

	if len(typeAnnotation.TypeParams) != 2 {
		t.Fatalf("length of typeAnnotation.TypeParams wrong. got=%d", len(typeAnnotation.TypeParams))
	}

	stringParam, ok := typeAnnotation.TypeParams[0].(*ast.TypeAnnotation)
	if !ok || stringParam.TypeName != "string" {
		t.Fatalf("First type parameter should be 'string'. got=%v", typeAnnotation.TypeParams[0])
	}

	numberParam, ok := typeAnnotation.TypeParams[1].(*ast.TypeAnnotation)
	if !ok || numberParam.TypeName != "number" {
		t.Fatalf("Second type parameter should be 'number'. got=%v", typeAnnotation.TypeParams[1])
	}
}

// TestNestedGenericTypes tests parsing of types with nested generic parameters
func TestNestedGenericTypes(t *testing.T) {
	input := `type NestedArray = Array<Array<string>>`

	l := lexer.New(input)
	p := New(l)
	program, err := p.Parse()

	if err != nil {
		t.Fatalf("parser error: %v", err)
	}

	stmt, ok := program.Statements[0].(*ast.TypeDeclaration)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.TypeDeclaration. got=%T", program.Statements[0])
	}

	typeAnnotation, ok := stmt.TypeValue.(*ast.TypeAnnotation)
	if !ok {
		t.Fatalf("stmt.TypeValue is not ast.TypeAnnotation. got=%T", stmt.TypeValue)
	}

	if len(typeAnnotation.TypeParams) != 1 {
		t.Fatalf("length of typeAnnotation.TypeParams wrong. got=%d", len(typeAnnotation.TypeParams))
	}

	innerType, ok := typeAnnotation.TypeParams[0].(*ast.TypeAnnotation)
	if !ok || innerType.TypeName != "Array" {
		t.Fatalf("Inner type should be 'Array'. got=%v", typeAnnotation.TypeParams[0])
	}

	if len(innerType.TypeParams) != 1 {
		t.Fatalf("Inner Array should have 1 type parameter. got=%d", len(innerType.TypeParams))
	}

	stringParam, ok := innerType.TypeParams[0].(*ast.TypeAnnotation)
	if !ok || stringParam.TypeName != "string" {
		t.Fatalf("Innermost type parameter should be 'string'. got=%v", innerType.TypeParams[0])
	}
}

// TestVariableDeclarationWithTypeAnnotation tests parsing of variable declarations with type annotations
func TestVariableDeclarationWithTypeAnnotation(t *testing.T) {
	input := `a: string = "hello"`

	l := lexer.New(input)
	p := New(l)
	program, err := p.Parse()

	if err != nil {
		t.Fatalf("parser error: %v", err)
	}

	if len(program.Statements) != 1 {
		t.Fatalf("program has not enough statements. got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.VariableDecl)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.VariableDecl. got=%T", program.Statements[0])
	}

	if stmt.Name != "a" {
		t.Errorf("stmt.Name not 'a'. got=%q", stmt.Name)
	}

	typeAnnotation := stmt.TypeAnnotation
	if typeAnnotation == nil {
		t.Fatalf("stmt.TypeAnnotation is nil")
	}

	if typeAnnotation.TypeName != "string" {
		t.Errorf("typeAnnotation.TypeName not 'string'. got=%q", typeAnnotation.TypeName)
	}

	stringLiteral, ok := stmt.Value.(*ast.StringLiteral)
	if !ok {
		t.Fatalf("stmt.Value is not ast.StringLiteral. got=%T", stmt.Value)
	}

	if stringLiteral.Value != "hello" {
		t.Errorf("stringLiteral.Value not 'hello'. got=%q", stringLiteral.Value)
	}
}

// TestVariableDeclarationWithoutType tests parsing of variable declarations without type annotations
func TestVariableDeclarationWithoutType(t *testing.T) {
	input := `x = 5
	y = 10
	z = x + y`

	l := lexer.New(input)
	p := New(l)
	program, err := p.Parse()

	if err != nil {
		t.Fatalf("parser error: %v", err)
	}

	if len(program.Statements) != 3 {
		t.Fatalf("program does not have 3 statements. got=%d", len(program.Statements))
	}

	tests := []struct {
		expectedIdentifier string
	}{
		{"x"},
		{"y"},
		{"z"},
	}

	for i, tt := range tests {
		stmt := program.Statements[i]
		if !testVariableDeclaration(t, stmt, tt.expectedIdentifier) {
			return
		}
	}
}

func testVariableDeclaration(t *testing.T, s ast.Node, name string) bool {
	// The parser now creates Assignment nodes for simple assignments like x = 5
	if assignment, ok := s.(*ast.Assignment); ok {
		if assignment.Name != name {
			t.Errorf("assignment.Name not '%s'. got=%s", name, assignment.Name)
			return false
		}
		return true
	}

	// For variable declarations with type annotations, it creates VariableDecl nodes
	varDecl, ok := s.(*ast.VariableDecl)
	if !ok {
		t.Errorf("s not *ast.VariableDecl or *ast.Assignment. got=%T", s)
		return false
	}

	if varDecl.Name != name {
		t.Errorf("varDecl.Name not '%s'. got=%s", name, varDecl.Name)
		return false
	}

	return true
}