package interpreter

import (
	"testing"

	"github.com/vibe-lang/vibe/lexer"
	"github.com/vibe-lang/vibe/object"
	"github.com/vibe-lang/vibe/parser"
)

// evalTest is a helper function for testing the evaluation of expressions
func evalTest(t *testing.T, input string) object.Object {
	l := lexer.New(input)
	// We don't need to create a parser instance if using parser.Parse directly
	program, errors := parser.Parse(l)
	if len(errors) > 0 {
		t.Fatalf("parser errors: %v", errors)
	}

	env := object.NewEnvironment()
	return Eval(program, env)
}

// TestBasicExpressionEvaluation tests the complete workflow from lexing to parsing to evaluation
func TestBasicExpressionEvaluation(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		// Integer expressions
		{"5", 5},
		{"10", 10},
		{"-5", -5},
		{"5 + 5 + 5 + 5 - 10", 10},
		{"2 * 2 * 2 * 2 * 2", 32},
		{"5 * 2 + 10", 20},
		{"5 + 2 * 10", 25},
		{"20 + 2 * -10", 0},
		{"50 / 2 * 2 + 10", 60},
		{"2 * (5 + 10)", 30},
		{"(5 + 10 * 2 + 15 / 3) * 2 + -10", 50},

		// Boolean expressions
		{"true", true},
		{"false", false},
		{"1 < 2", true},
		{"1 > 2", false},
		{"1 == 1", true},
		{"1 != 1", false},
		{"true == true", true},
		{"true != false", true},

		// String expressions
		{`"Hello World!"`, "Hello World!"},
		{`"Hello" + " " + "World!"`, "Hello World!"},
	}

	for _, tt := range tests {
		evaluated := evalTest(t, tt.input)

		switch expected := tt.expected.(type) {
		case int:
			testIntegerObject(t, evaluated, int64(expected))
		case bool:
			testBooleanObject(t, evaluated, expected)
		case string:
			testStringObject(t, evaluated, expected)
		}
	}
}

// TestStatementEvaluation tests complete statement evaluation
func TestStatementEvaluation(t *testing.T) {
	tests := []struct {
		desc     string
		input    string
		expected interface{}
	}{
		{
			"Variable assignment and retrieval",
			`x: int = 5
			y: int = 10
			z: int = x + y
			z`,
			15,
		},
		{
			"Variable assignment and retrieval without type annotation",
			`x = 5
			y = 10
			z = x + y
			z`,
			15,
		},
		{
			"If-else statement",
			`x: int = 10
			y: int = 0

			if x > 5 do
				y = 1
			else
				y = 2
			end

			y`,
			1,
		},
		{
			"For loop with array",
			`sum: int = 0
			for i in [1, 2, 3, 4, 5] do
				sum = sum + i
			end
			sum`,
			15,
		},
		{
			"Function definition and call",
			`def add(x: int, y: int): int do
				x + y
			end
			add(2, 5)`,
			7,
		},
		{
			"Closure",
			`def makeAdder(x: int): function do
				def inner(y: int): int do
					x + y
				end
				inner
			end
			adder = makeAdder(5)
			adder(10)`,
			15,
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			evaluated := evalTest(t, tt.input)

			switch expected := tt.expected.(type) {
			case int:
				testIntegerObject(t, evaluated, int64(expected))
			case bool:
				testBooleanObject(t, evaluated, expected)
			case string:
				testStringObject(t, evaluated, expected)
			}
		})
	}
}

// Test helper functions for validating returned objects
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

func testStringObject(t *testing.T, obj object.Object, expected string) bool {
	result, ok := obj.(*object.String)
	if !ok {
		t.Errorf("object is not String. got=%T (%+v)", obj, obj)
		return false
	}
	if result.Value != expected {
		t.Errorf("object has wrong value. got=%q, want=%q", result.Value, expected)
		return false
	}
	return true
}