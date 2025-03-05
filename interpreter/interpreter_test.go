package interpreter

import (
	"fmt"
	"testing"

	"github.com/vibe-lang/vibe/ast"
	"github.com/vibe-lang/vibe/lexer"
	"github.com/vibe-lang/vibe/parser"
)

func TestEvalIntegerExpression(t *testing.T) {
	tests := []struct {
		input    int
		expected int
	}{
		{5, 5},
		{10, 10},
	}

	for _, tt := range tests {
		// Create a program with a single number literal
		program := &ast.Program{
			Statements: []ast.Node{
				&ast.NumberLiteral{
					Value: float64(tt.input),
					IsInt: true,
				},
			},
		}

		interp := New()
		evaluated := interp.Eval(program)

		if !testIntegerValue(t, evaluated, tt.expected) {
			t.Errorf("Failed test for input: %d", tt.input)
		}
	}
}

func TestEvalBooleanExpression(t *testing.T) {
	tests := []struct {
		input    bool
		expected bool
	}{
		{true, true},
		{false, false},
	}

	for _, tt := range tests {
		// Create a program with a single boolean literal
		program := &ast.Program{
			Statements: []ast.Node{
				&ast.BooleanLiteral{
					Value: tt.input,
				},
			},
		}

		interp := New()
		evaluated := interp.Eval(program)

		if !testBooleanValue(t, evaluated, tt.expected) {
			t.Errorf("Failed test for input: %t", tt.input)
		}
	}
}

func TestIfElseExpressions(t *testing.T) {
	tests := []struct {
		condition bool
		consequence int
		alternative int
		expected int
	}{
		{true, 10, 20, 10},   // if true then 10 else 20
		{false, 10, 20, 20},  // if false then 10 else 20
	}

	for _, tt := range tests {
		// Create a program with an if statement directly
		program := &ast.Program{
			Statements: []ast.Node{
				&ast.IfStmt{
					Condition: &ast.BooleanLiteral{Value: tt.condition},
					Consequence: &ast.BlockStmt{
						Statements: []ast.Node{
							&ast.NumberLiteral{Value: float64(tt.consequence), IsInt: true},
						},
					},
					Alternative: &ast.BlockStmt{
						Statements: []ast.Node{
							&ast.NumberLiteral{Value: float64(tt.alternative), IsInt: true},
						},
					},
				},
			},
		}

		interp := New()
		evaluated := interp.Eval(program)

		if !testIntegerValue(t, evaluated, tt.expected) {
			t.Errorf("If-else test failed. Expected %d, got %v",
				tt.expected, evaluated.Inspect())
		}
	}
}

func TestReturnStatements(t *testing.T) {
	// Test a simple return statement that should work
	input := "return 5;"
	evaluated := testEval(input)

	if !testIntegerValue(t, evaluated, 5) {
		t.Errorf("Failed test for input: %s", input)
	}
}

func TestFunctionApplication(t *testing.T) {
	// Test a simple function application
	// Create a program with a function definition and a call
	program := &ast.Program{
		Statements: []ast.Node{
			// Define a function 'add' that takes two parameters and returns their sum
			&ast.FunctionDef{
				Name: "add",
				Parameters: []ast.Parameter{
					{Name: "a", Type: &ast.TypeAnnotation{TypeName: "int"}},
					{Name: "b", Type: &ast.TypeAnnotation{TypeName: "int"}},
				},
				ReturnType: &ast.TypeAnnotation{TypeName: "int"},
				Body: &ast.BlockStmt{
					Statements: []ast.Node{
						&ast.BinaryExpr{
							Left: &ast.Identifier{Name: "a"},
							Operator: "+",
							Right: &ast.Identifier{Name: "b"},
						},
					},
				},
			},
			// Call the function with arguments 5 and 7
			&ast.CallExpr{
				Function: &ast.Identifier{Name: "add"},
				Args: []ast.Node{
					&ast.NumberLiteral{Value: 5, IsInt: true},
					&ast.NumberLiteral{Value: 7, IsInt: true},
				},
			},
		},
	}

	interp := New()
	evaluated := interp.Eval(program)

	if !testIntegerValue(t, evaluated, 12) {
		t.Errorf("Function application test failed. Expected 12, got %v", evaluated.Inspect())
	}
}

func TestStringConcatenation(t *testing.T) {
	// For now, we'll just test a simple string
	input := `"Hello World!"`
	evaluated := testEval(input)

	str, ok := evaluated.(*StringValue)
	if !ok {
		t.Fatalf("Object is not String. got=%T (%+v)", evaluated, evaluated)
	}

	if str.Value != "Hello World!" {
		t.Errorf("String has wrong value. got=%q", str.Value)
	}
}

// For now, we'll skip the array test since we haven't fully implemented it yet

func TestTypeSystem(t *testing.T) {
	// Create a program with typed variables and a function call
	program := &ast.Program{
		Statements: []ast.Node{
			// Variable declaration: a: int = 5
			&ast.VariableDecl{
				Name: "a",
				TypeAnnotation: &ast.TypeAnnotation{TypeName: "int"},
				Value: &ast.NumberLiteral{Value: 5, IsInt: true},
			},

			// Function definition: identity(value: int) -> int
			&ast.FunctionDef{
				Name: "identity",
				Parameters: []ast.Parameter{
					{Name: "value", Type: &ast.TypeAnnotation{TypeName: "int"}},
				},
				ReturnType: &ast.TypeAnnotation{TypeName: "int"},
				Body: &ast.BlockStmt{
					Statements: []ast.Node{
						&ast.ReturnStmt{
							Value: &ast.Identifier{Name: "value"},
						},
					},
				},
			},

			// Function call: identity(50)
			&ast.CallExpr{
				Function: &ast.Identifier{Name: "identity"},
				Args: []ast.Node{
					&ast.NumberLiteral{Value: 50, IsInt: true},
				},
			},
		},
	}

	interp := New()
	evaluated := interp.Eval(program)

	// The result should be the value passed to the identity function
	if !testIntegerValue(t, evaluated, 50) {
		t.Errorf("Type system test failed. Expected result to be 50, got %v", evaluated.Inspect())
	}
}

// Helper functions

func testEval(input string) Value {
	l := lexer.New(input)
	p, errors := parser.Parse(l)

	// If there are parser errors, print them for debugging
	if len(errors) > 0 {
		fmt.Printf("Parser errors for input:\n")
		for _, err := range errors {
			fmt.Printf("  - %s\n", err)
		}
	}

	// Debug output removed for clarity

	interp := New()
	return interp.Eval(p)
}

func testIntegerValue(t *testing.T, obj Value, expected int) bool {
	result, ok := obj.(*IntegerValue)
	if !ok {
		t.Errorf("object is not IntegerValue. got=%T (%+v)", obj, obj)
		return false
	}
	if result.Value != int64(expected) {
		t.Errorf("object has wrong value. got=%d, want=%d", result.Value, expected)
		return false
	}
	return true
}

func testBooleanValue(t *testing.T, obj Value, expected bool) bool {
	result, ok := obj.(*BooleanValue)
	if !ok {
		t.Errorf("object is not BooleanValue. got=%T (%+v)", obj, obj)
		return false
	}
	if result.Value != expected {
		t.Errorf("object has wrong value. got=%t, want=%t", result.Value, expected)
		return false
	}
	return true
}

func testNilValue(t *testing.T, obj Value) bool {
	_, ok := obj.(*NilValue)
	if !ok {
		t.Errorf("object is not NilValue. got=%T (%+v)", obj, obj)
		return false
	}
	return true
}