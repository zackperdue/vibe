package parser

import (
	"testing"

	"github.com/vibe-lang/vibe/ast"
	"github.com/vibe-lang/vibe/lexer"
)

// TestLetStatements tests the parsing of let statements
func TestLetStatements(t *testing.T) {
	tests := []struct {
		input              string
		expectedIdentifier string
		expectedValue      interface{}
	}{
		{"let x = 5;", "x", 5},
		{"let y = true;", "y", true},
		{"let foobar = y;", "foobar", "y"},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program, err := p.Parse()

		if err != nil {
			t.Fatalf("parser error: %v", err)
		}

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain 1 statement. got=%d",
				len(program.Statements))
		}

		stmt := program.Statements[0]
		if !testLetStatement(t, stmt, tt.expectedIdentifier) {
			return
		}

		val := stmt.(*ast.LetStatement).Value
		if !testLiteralExpression(t, val, tt.expectedValue) {
			return
		}
	}
}

func testLetStatement(t *testing.T, s ast.Statement, name string) bool {
	if s.TokenLiteral() != "let" {
		t.Errorf("s.TokenLiteral not 'let'. got=%q", s.TokenLiteral())
		return false
	}

	letStmt, ok := s.(*ast.LetStatement)
	if !ok {
		t.Errorf("s not *ast.LetStatement. got=%T", s)
		return false
	}

	if letStmt.Name.Value != name {
		t.Errorf("letStmt.Name.Value not '%s'. got=%s", name, letStmt.Name.Value)
		return false
	}

	if letStmt.Name.TokenLiteral() != name {
		t.Errorf("letStmt.Name.TokenLiteral() not '%s'. got=%s",
			name, letStmt.Name.TokenLiteral())
		return false
	}

	return true
}

// TestReturnStatements tests the parsing of return statements
func TestReturnStatements(t *testing.T) {
	tests := []struct {
		input         string
		expectedValue interface{}
	}{
		{"return 5;", 5},
		{"return true;", true},
		{"return foobar;", "foobar"},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program, err := p.Parse()

		if err != nil {
			t.Fatalf("parser error: %v", err)
		}

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain 1 statement. got=%d",
				len(program.Statements))
		}

		stmt := program.Statements[0]
		returnStmt, ok := stmt.(*ast.ReturnStatement)
		if !ok {
			t.Fatalf("stmt not *ast.ReturnStatement. got=%T", stmt)
		}

		if returnStmt.TokenLiteral() != "return" {
			t.Fatalf("returnStmt.TokenLiteral not 'return', got %q",
				returnStmt.TokenLiteral())
		}

		if testLiteralExpression(t, returnStmt.ReturnValue, tt.expectedValue) {
			return
		}
	}
}

// TestForLoopStatement tests the parsing of for loop statements
func TestForLoopStatement(t *testing.T) {
	input := `for i in [1, 2, 3] do
		let x = i * 2;
	end`

	l := lexer.New(input)
	p := New(l)
	program, err := p.Parse()

	if err != nil {
		t.Fatalf("parser error: %v", err)
	}

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statement. got=%d",
			len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ForLoopStatement)
	if !ok {
		t.Fatalf("stmt is not ast.ForLoopStatement. got=%T", program.Statements[0])
	}

	if stmt.Iterator.String() != "i" {
		t.Errorf("iterator is not 'i'. got=%s", stmt.Iterator.String())
	}

	arrayLit, ok := stmt.Iterable.(*ast.ArrayLiteral)
	if !ok {
		t.Fatalf("stmt.Iterable is not ast.ArrayLiteral. got=%T", stmt.Iterable)
	}

	if len(arrayLit.Elements) != 3 {
		t.Errorf("array.Elements does not contain 3 elements. got=%d", len(arrayLit.Elements))
	}

	testIntegerLiteral(t, arrayLit.Elements[0], 1)
	testIntegerLiteral(t, arrayLit.Elements[1], 2)
	testIntegerLiteral(t, arrayLit.Elements[2], 3)

	bodyStmt, ok := stmt.Body.Statements[0].(*ast.LetStatement)
	if !ok {
		t.Fatalf("stmt.Body.Statements[0] is not ast.LetStatement. got=%T",
			stmt.Body.Statements[0])
	}

	if bodyStmt.Name.Value != "x" {
		t.Errorf("body variable name not 'x'. got=%s", bodyStmt.Name.Value)
	}
}