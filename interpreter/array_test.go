package interpreter

import (
	"testing"

	"github.com/vibe-lang/vibe/lexer"
	"github.com/vibe-lang/vibe/parser"
)

// evalForTest is a helper function for evaluating code in tests
func evalForTest(t *testing.T, input string) Value {
	l := lexer.New(input)
	program, errors := parser.Parse(l)
	if len(errors) > 0 {
		t.Fatalf("parser errors: %v", errors)
	}

	interp := New()
	return interp.Eval(program)
}

// TestArrayLiterals tests array literals
func TestArrayLiterals(t *testing.T) {
	input := "[1, 2 * 2, 3 + 3]"

	evaluated := evalForTest(t, input)
	result, ok := evaluated.(*ArrayValue)
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
	// Test directly with array values and index operations
	array := &ArrayValue{
		Elements: []Value{
			&IntegerValue{Value: 1},
			&IntegerValue{Value: 2},
			&IntegerValue{Value: 3},
		},
	}

	// Test index 0
	index0 := &IntegerValue{Value: 0}
	result0 := array.Index(index0)
	testInteger(t, result0, 1)

	// Test index 1
	index1 := &IntegerValue{Value: 1}
	result1 := array.Index(index1)
	testInteger(t, result1, 2)

	// Test index 2
	index2 := &IntegerValue{Value: 2}
	result2 := array.Index(index2)
	testInteger(t, result2, 3)
}

// Helper function to test integer values
func testInteger(t *testing.T, obj Value, expected int64) bool {
	result, ok := obj.(*IntegerValue)
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