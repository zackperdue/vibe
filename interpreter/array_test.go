package interpreter

import (
	"testing"

	"github.com/vibe-lang/vibe/lexer"
	"github.com/vibe-lang/vibe/object"
	"github.com/vibe-lang/vibe/parser"
)

// evalForTest is a helper function for evaluating code in tests
func evalForTest(t *testing.T, input string) object.Object {
	l := lexer.New(input)
	program, errors := parser.Parse(l)
	if len(errors) > 0 {
		t.Fatalf("parser errors: %v", errors)
	}

	env := object.NewEnvironment()
	return Eval(program, env)
}

// TestArrayLiterals tests array literals
func TestArrayLiterals(t *testing.T) {
	input := "[1, 2 * 2, 3 + 3]"

	evaluated := evalForTest(t, input)
	result, ok := evaluated.(*object.Array)
	if !ok {
		t.Fatalf("object is not Array. got=%T (%+v)", evaluated, evaluated)
	}

	if len(result.Elements) != 3 {
		t.Fatalf("array has wrong num of elements. got=%d",
			len(result.Elements))
	}

	testInteger(t, result.Elements[0], 1)
	testInteger(t, result.Elements[1], 4)
	testInteger(t, result.Elements[2], 6)
}

// TestArrayIndexExpressions tests array index expressions
func TestArrayIndexExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{
			"[1, 2, 3][0]",
			1,
		},
		{
			"[1, 2, 3][1]",
			2,
		},
		{
			"[1, 2, 3][2]",
			3,
		},
		{
			"i: int = 0; [1][i];",
			1,
		},
		{
			"[1, 2, 3][1 + 1];",
			3,
		},
		{
			"myArray: array = [1, 2, 3]; myArray[2];",
			3,
		},
		{
			"myArray: array = [1, 2, 3]; myArray[0] + myArray[1] + myArray[2];",
			6,
		},
		{
			"myArray: array = [1, 2, 3]; i: int = myArray[0]; myArray[i]",
			2,
		},
	}

	for _, tt := range tests {
		evaluated := evalForTest(t, tt.input)

		switch expected := tt.expected.(type) {
		case int:
			testInteger(t, evaluated, int64(expected))
		default:
			t.Errorf("Unknown expected type: %T", expected)
		}
	}
}

// Helper function to test integer values
func testInteger(t *testing.T, obj object.Object, expected int64) bool {
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