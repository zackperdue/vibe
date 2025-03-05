package interpreter

import (
	"testing"

	"github.com/vibe-lang/vibe/lexer"
	"github.com/vibe-lang/vibe/object"
	"github.com/vibe-lang/vibe/parser"
)

// testEval2 is a helper function for evaluating code in tests
func testEval2(t *testing.T, input string) object.Object {
	l := lexer.New(input)
	program, errors := parser.Parse(l)
	if len(errors) > 0 {
		t.Fatalf("parser errors: %v", errors)
	}

	env := object.NewEnvironment()
	return Eval(program, env)
}

// TestFunctionObject tests function objects
func TestFunctionObject(t *testing.T) {
	input := "def f(x: int): int do x + 2 end"

	evaluated := testEval2(t, input)
	fn, ok := evaluated.(*object.Function)
	if !ok {
		t.Fatalf("object is not Function. got=%T (%+v)", evaluated, evaluated)
	}

	if len(fn.Parameters) != 1 {
		t.Fatalf("function has wrong parameters. Parameters=%+v",
			fn.Parameters)
	}

	if fn.Parameters[0].String() != "x: int" {
		t.Fatalf("parameter is not 'x: int'. got=%q", fn.Parameters[0])
	}

	expectedBody := "(x + 2)"

	if fn.Body.String() != expectedBody {
		t.Fatalf("body is not %q. got=%q", expectedBody, fn.Body.String())
	}
}

// TestEnhancedFunctionApplication tests function application with additional test cases
func TestEnhancedFunctionApplication(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"identity = def(x: int): int do x end; identity(5);", 5},
		{"identity = def(x: int): int do return x end; identity(5);", 5},
		{"double = def(x: int): int do x * 2 end; double(5);", 10},
		{"add = def(x: int, y: int): int do x + y end; add(5, 5);", 10},
		{"add = def(x: int, y: int): int do x + y end; add(5 + 5, add(5, 5));", 20},
		{"def f(x: int): int do x end; f(5)", 5},
	}

	for _, tt := range tests {
		validateIntegerObject(t, testEval2(t, tt.input), tt.expected)
	}
}

// TestClosures tests closures
func TestClosures(t *testing.T) {
	input := `
let newAdder = function(x) do
  function(y) do x + y end
end;

let addTwo = newAdder(2);
addTwo(2);`

	validateIntegerObject(t, testEval2(t, input), 4)
}

// Helper function to validate integer objects in this file
func validateIntegerObject(t *testing.T, obj object.Object, expected int64) bool {
	result, ok := obj.(*object.Integer)
	if !ok {
		t.Errorf("object is not Integer. got=%T (%+v)", obj, obj)
		return false
	}
	if result.Value != expected {
		t.Errorf("object has wrong value. got=%d, want=%d", result.Value, expected)
		return false
	}
	return true
}