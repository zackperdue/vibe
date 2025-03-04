package integration_test

import (
	"testing"

	"github.com/vibe-lang/vibe/interpreter"
	"github.com/vibe-lang/vibe/lexer"
	"github.com/vibe-lang/vibe/object"
	"github.com/vibe-lang/vibe/parser"
)

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
		evaluated := testEval(t, tt.input)

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
			`
			let x = 5
			let y = 10
			let z = x + y
			z
			`,
			15,
		},
		{
			"If-else statement",
			`
			let x = 10
			let y = 0

			if x > 5 do
				y = 1
			else
				y = 2
			end

			y
			`,
			1,
		},
		{
			"For loop with array",
			`
			let sum = 0
			for i in [1, 2, 3, 4, 5] do
				sum = sum + i
			end
			sum
			`,
			15,
		},
		{
			"Function definition and call",
			`
			function add(a, b) do
				return a + b
			end

			add(5, 10)
			`,
			15,
		},
		{
			"Closure",
			`
			function makeAdder(x) do
				function(y) do
					return x + y
				end
			end

			let addFive = makeAdder(5)
			addFive(10)
			`,
			15,
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			evaluated := testEval(t, tt.input)

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

// TestErrorHandling tests error handling throughout the pipeline
func TestErrorHandling(t *testing.T) {
	tests := []struct {
		desc            string
		input           string
		expectedMessage string
	}{
		{
			"Type mismatch error",
			"5 + true",
			"type mismatch: INTEGER + BOOLEAN",
		},
		{
			"Unknown operator error",
			"5 * true",
			"type mismatch: INTEGER * BOOLEAN",
		},
		{
			"Unknown identifier error",
			"foobar",
			"identifier not found: foobar",
		},
		{
			"Function argument count mismatch",
			`
			function add(a, b) do
				return a + b
			end

			add(1)
			`,
			"wrong number of arguments: got=1, want=2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			evaluated := testEval(t, tt.input)

			errObj, ok := evaluated.(*object.Error)
			if !ok {
				t.Fatalf("expected error object, got=%T (%+v)", evaluated, evaluated)
			}

			if errObj.Message != tt.expectedMessage {
				t.Errorf("wrong error message. expected=%q, got=%q",
					tt.expectedMessage, errObj.Message)
			}
		})
	}
}

// TestComplexProgram tests a more complex program combining multiple features
func TestComplexProgram(t *testing.T) {
	input := `
	function map(arr, fn) do
		let result = []
		for element in arr do
			result = result + [fn(element)]
		end
		return result
	end

	function filter(arr, predicate) do
		let result = []
		for element in arr do
			if predicate(element) do
				result = result + [element]
			end
		end
		return result
	end

	let numbers = [1, 2, 3, 4, 5, 6, 7, 8, 9, 10]
	let isEven = function(x) do x % 2 == 0 end
	let square = function(x) do x * x end

	let evenNumbers = filter(numbers, isEven)
	let squaredEvenNumbers = map(evenNumbers, square)

	// Sum the squared even numbers
	let sum = 0
	for num in squaredEvenNumbers do
		sum = sum + num
	end

	sum
	`

	evaluated := testEval(t, input)
	// The sum of squares of even numbers from 1 to 10: 2² + 4² + 6² + 8² + 10² = 4 + 16 + 36 + 64 + 100 = 220
	testIntegerObject(t, evaluated, 220)
}

// Helper functions

func testEval(t *testing.T, input string) object.Object {
	l := lexer.New(input)
	p, errors := parser.Parse(l)
	if len(errors) != 0 {
		t.Fatalf("parser errors: %v", errors)
	}

	env := object.NewEnvironment()
	return interpreter.Eval(p, env)
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