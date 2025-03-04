package interpreter_test

import (
	"testing"

	"github.com/vibe-lang/vibe/interpreter"
	"github.com/vibe-lang/vibe/lexer"
	"github.com/vibe-lang/vibe/object"
	"github.com/vibe-lang/vibe/parser"
)

// TestEvalIntegerExpression tests the evaluation of integer expressions
func TestEvalIntegerExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"5", 5},
		{"10", 10},
		{"-5", -5},
		{"-10", -10},
		{"5 + 5", 10},
		{"5 - 5", 0},
		{"5 * 5", 25},
		{"5 / 5", 1},
		{"5 + 5 + 5", 15},
		{"5 * 5 * 5", 125},
		{"5 * (5 + 5)", 50},
		{"(5 + 5) * 5", 50},
	}

	for _, tt := range tests {
		evaluated := testEval(t, tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

// TestEvalBooleanExpression tests the evaluation of boolean expressions
func TestEvalBooleanExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"true", true},
		{"false", false},
		{"5 > 3", true},
		{"5 < 3", false},
		{"5 == 5", true},
		{"5 != 5", false},
		{"5 != 6", true},
		{"5 == 6", false},
		{"true == true", true},
		{"false == false", true},
		{"true == false", false},
		{"true != false", true},
		{"false != true", true},
		{"(5 > 3) == true", true},
		{"(5 < 3) == true", false},
		{"(5 > 3) != false", true},
	}

	for _, tt := range tests {
		evaluated := testEval(t, tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}
}

// TestEvalIdentifier tests the evaluation of identifiers
func TestEvalIdentifier(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"let a = 5; a;", 5},
		{"let a = 5; let b = a; b;", 5},
		{"let a = 5; let b = a; let c = a + b; c;", 10},
	}

	for _, tt := range tests {
		evaluated := testEval(t, tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

// TestEvalIfExpression tests the evaluation of if expressions
func TestEvalIfExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"if (true) { 10 }", 10},
		{"if (false) { 10 }", nil},
		{"if (1) { 10 }", 10},
		{"if (1 < 2) { 10 }", 10},
		{"if (1 > 2) { 10 }", nil},
		{"if (1 > 2) { 10 } else { 20 }", 20},
		{"if (1 < 2) { 10 } else { 20 }", 10},
	}

	for _, tt := range tests {
		evaluated := testEval(t, tt.input)

		integer, ok := tt.expected.(int)
		if ok {
			testIntegerObject(t, evaluated, int64(integer))
		} else {
			testNilObject(t, evaluated)
		}
	}
}

// TestEvalArrayLiteral tests the evaluation of array literals
func TestEvalArrayLiteral(t *testing.T) {
	input := "[1, 2 * 2, 3 + 3]"

	evaluated := testEval(t, input)
	result, ok := evaluated.(*object.Array)
	if !ok {
		t.Fatalf("object is not Array. got=%T (%+v)", evaluated, evaluated)
	}

	if len(result.Elements) != 3 {
		t.Fatalf("array has wrong num of elements. got=%d", len(result.Elements))
	}

	testIntegerObject(t, result.Elements[0], 1)
	testIntegerObject(t, result.Elements[1], 4)
	testIntegerObject(t, result.Elements[2], 6)
}

// TestEvalIndexExpression tests the evaluation of index expressions
func TestEvalIndexExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"[1, 2, 3][0]", 1},
		{"[1, 2, 3][1]", 2},
		{"[1, 2, 3][2]", 3},
		{"let i = 0; [1][i];", 1},
		{"[1, 2, 3][1 + 1];", 3},
		{"let myArray = [1, 2, 3]; myArray[2];", 3},
		{"let myArray = [1, 2, 3]; myArray[0] + myArray[1] + myArray[2];", 6},
		{"let myArray = [1, 2, 3]; let i = myArray[0]; myArray[i]", 2},
	}

	for _, tt := range tests {
		evaluated := testEval(t, tt.input)
		integer, ok := tt.expected.(int)
		if ok {
			testIntegerObject(t, evaluated, int64(integer))
		} else {
			testNilObject(t, evaluated)
		}
	}
}

// TestEvalForLoop tests the evaluation of for loops
func TestEvalForLoop(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{`
		let sum = 0;
		for i in [1, 2, 3, 4, 5] {
			sum = sum + i;
		}
		sum
		`, 15},
		{`
		let sum = 0;
		let numbers = [1, 2, 3, 4, 5];
		for i in numbers {
			sum = sum + i;
		}
		sum
		`, 15},
	}

	for _, tt := range tests {
		evaluated := testEval(t, tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

// Helper functions
func testEval(t *testing.T, input string) object.Object {
	l := lexer.New(input)
	program, errors := parser.Parse(l)

	if len(errors) > 0 {
		t.Fatalf("parser errors: %v", errors)
	}

	env := object.NewEnvironment()
	return interpreter.Eval(program, env)
}

func testIntegerObject(t *testing.T, obj object.Object, expected int64) bool {
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

func testBooleanObject(t *testing.T, obj object.Object, expected bool) bool {
	result, ok := obj.(*object.Boolean)
	if !ok {
		t.Errorf("object is not Boolean. got=%T (%+v)", obj, obj)
		return false
	}
	if result.Value != expected {
		t.Errorf("object has wrong value. got=%t, want=%t", result.Value, expected)
		return false
	}
	return true
}

func testNilObject(t *testing.T, obj object.Object) bool {
	if obj != object.NULL {
		t.Errorf("object is not NULL. got=%T (%+v)", obj, obj)
		return false
	}
	return true
}