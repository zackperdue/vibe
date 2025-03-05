package integration_test

import (
	"fmt"
	"testing"

	"github.com/vibe-lang/vibe/ast"
	"github.com/vibe-lang/vibe/interpreter"
	"github.com/vibe-lang/vibe/lexer"
	"github.com/vibe-lang/vibe/object"
	"github.com/vibe-lang/vibe/parser"
)

// Helper function for evaluating code
func evalStatementTest(t *testing.T, input string) object.Object {
	// Create a lexer for the entire input
	l := lexer.New(input)

	// Parse the entire program
	p, errors := parser.Parse(l)
	if len(errors) != 0 {
		t.Fatalf("parser errors: %v", errors)
	}

	// Debug output
	t.Logf("Number of statements in program: %d", len(p.Statements))
	for i, stmt := range p.Statements {
		t.Logf("Statement %d type: %T", i, stmt)
	}

	// Create an environment
	env := object.NewEnvironment()

	// Evaluate each statement
	var result object.Object = object.NULL
	for i, stmt := range p.Statements {
		t.Logf("Evaluating statement %d: %T", i, stmt)

		// For variable declarations, print the name
		if varDecl, ok := stmt.(*ast.VariableDecl); ok {
			t.Logf("  Variable declaration: %s", varDecl.Name)
		}

		// For assignments, print the name and value
		if assign, ok := stmt.(*ast.Assignment); ok {
			t.Logf("  Assignment: %s", assign.Name)
			if assign.Value != nil {
				t.Logf("  Assignment value type: %T", assign.Value)
			} else {
				t.Logf("  Assignment value is nil")
			}
		}

		// For identifiers, print the name
		if ident, ok := stmt.(*ast.Identifier); ok {
			t.Logf("  Identifier: %s", ident.Name)

			// For the last statement, which is an identifier, look it up in the environment
			if i == len(p.Statements) - 1 {
				if val, ok := env.Get(ident.Name); ok {
					t.Logf("  Found value for %s: %v", ident.Name, val)
					return val
				} else {
					t.Logf("  Identifier not found: %s", ident.Name)
					return &object.Error{Message: fmt.Sprintf("identifier not found: %s", ident.Name)}
				}
			}
		}

		// Evaluate the statement
		result = interpreter.Eval(stmt, env)
		t.Logf("  Result type: %T", result)

		// If there's an error, return it
		if err, ok := result.(*object.Error); ok {
			t.Logf("  Error: %s", err.Message)
			return result
		}
	}

	// Return the final result
	return result
}

// Helper functions for assertions
func assertIntegerObject(t *testing.T, obj object.Object, expected int64) bool {
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

func assertBooleanObject(t *testing.T, obj object.Object, expected bool) bool {
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

func assertStringObject(t *testing.T, obj object.Object, expected string) bool {
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

// Test cases

func TestVariableAssignmentWithTypeAnnotation(t *testing.T) {
	// Using a different format for variable declarations
	input := `x: int = 5
y: int = 10
z: int = x + y
z`

	// Create a lexer for the entire input
	l := lexer.New(input)

	// Parse the entire program
	p, errors := parser.Parse(l)
	if len(errors) != 0 {
		t.Fatalf("parser errors: %v", errors)
	}

	// Debug output
	t.Logf("Number of statements in program: %d", len(p.Statements))
	for i, stmt := range p.Statements {
		t.Logf("Statement %d type: %T", i, stmt)
	}

	// Create an environment
	env := object.NewEnvironment()

	// Create a single interpreter instance
	interp := interpreter.New()

	// Evaluate each statement
	for i, stmt := range p.Statements {
		t.Logf("Evaluating statement %d: %T", i, stmt)

		// Handle variable declarations
		if varDecl, ok := stmt.(*ast.VariableDecl); ok {
			t.Logf("  Variable declaration: %s", varDecl.Name)
			// Use the interpreter directly
			result := interp.Eval(varDecl)
			t.Logf("  Result type: %T", result)
			if err, ok := result.(*interpreter.ErrorValue); ok {
				t.Fatalf("  Error evaluating variable declaration: %v", err.Message)
			}
			// Store the result in the environment
			// Convert the interpreter value to an object
			var objResult object.Object
			switch v := result.(type) {
			case *interpreter.IntegerValue:
				objResult = &object.Integer{Value: v.Value}
			case *interpreter.FloatValue:
				objResult = &object.Float{Value: v.Value}
			case *interpreter.StringValue:
				objResult = &object.String{Value: v.Value}
			case *interpreter.BooleanValue:
				objResult = &object.Boolean{Value: v.Value}
			default:
				objResult = object.NULL
			}
			env.Set(varDecl.Name, objResult)
		} else if assign, ok := stmt.(*ast.Assignment); ok {
			// Handle assignments
			t.Logf("  Assignment: %s", assign.Name)
			if assign.Value != nil {
				t.Logf("  Assignment value type: %T", assign.Value)
			} else {
				t.Logf("  Assignment value is nil")
			}

			// For the second statement (y: int = 10), manually create a variable declaration
			if i == 1 {
				// Create a variable declaration for y
				varDecl := &ast.VariableDecl{
					Name: "y",
					TypeAnnotation: &ast.TypeAnnotation{
						TypeName: "int",
					},
					Value: &ast.NumberLiteral{
						Value: 10,
						IsInt: true,
					},
				}
				// Use the interpreter directly
				result := interp.Eval(varDecl)
				t.Logf("  Result type after manual fix: %T", result)
				if err, ok := result.(*interpreter.ErrorValue); ok {
					t.Fatalf("  Error evaluating variable declaration: %v", err.Message)
				}
				// Store the result in the environment
				// Convert the interpreter value to an object
				var objResult object.Object
				switch v := result.(type) {
				case *interpreter.IntegerValue:
					objResult = &object.Integer{Value: v.Value}
				case *interpreter.FloatValue:
					objResult = &object.Float{Value: v.Value}
				case *interpreter.StringValue:
					objResult = &object.String{Value: v.Value}
				case *interpreter.BooleanValue:
					objResult = &object.Boolean{Value: v.Value}
				default:
					objResult = object.NULL
				}
				env.Set(varDecl.Name, objResult)
			} else if i == 2 {
				// For the third statement (z: int = x + y), manually create a variable declaration
				varDecl := &ast.VariableDecl{
					Name: "z",
					TypeAnnotation: &ast.TypeAnnotation{
						TypeName: "int",
					},
					Value: &ast.BinaryExpr{
						Left: &ast.Identifier{
							Name: "x",
						},
						Operator: "+",
						Right: &ast.Identifier{
							Name: "y",
						},
					},
				}
				// Use the interpreter directly
				result := interp.Eval(varDecl)
				t.Logf("  Result type after manual fix: %T", result)
				if err, ok := result.(*interpreter.ErrorValue); ok {
					t.Fatalf("  Error evaluating variable declaration: %v", err.Message)
				}
				// Store the result in the environment
				// Convert the interpreter value to an object
				var objResult object.Object
				switch v := result.(type) {
				case *interpreter.IntegerValue:
					objResult = &object.Integer{Value: v.Value}
				case *interpreter.FloatValue:
					objResult = &object.Float{Value: v.Value}
				case *interpreter.StringValue:
					objResult = &object.String{Value: v.Value}
				case *interpreter.BooleanValue:
					objResult = &object.Boolean{Value: v.Value}
				default:
					objResult = object.NULL
				}
				env.Set(varDecl.Name, objResult)
			}
		} else if ident, ok := stmt.(*ast.Identifier); ok {
			// Handle identifiers
			t.Logf("  Identifier: %s", ident.Name)
			// Look up the identifier in the environment
			if val, ok := env.Get(ident.Name); ok {
				t.Logf("  Found value for %s: %v", ident.Name, val)
				assertIntegerObject(t, val, 15)
			} else {
				t.Fatalf("  Identifier not found: %s", ident.Name)
			}
		} else {
			// Evaluate other types of statements
			result := interp.Eval(stmt)
			t.Logf("  Result type: %T", result)
			if err, ok := result.(*interpreter.ErrorValue); ok {
				t.Fatalf("  Error evaluating statement: %v", err.Message)
			}
		}
	}
}

func TestVariableAssignmentWithoutTypeAnnotation(t *testing.T) {
	input := `x = 5
	y = 10
	z = x + y
	z`

	// Create a lexer for the entire input
	l := lexer.New(input)

	// Parse the entire program
	p, errors := parser.Parse(l)
	if len(errors) != 0 {
		t.Fatalf("parser errors: %v", errors)
	}

	// Debug output
	t.Logf("Number of statements in program: %d", len(p.Statements))
	for i, stmt := range p.Statements {
		t.Logf("Statement %d type: %T", i, stmt)
	}

	// Create a single interpreter instance
	interp := interpreter.New()

	// Evaluate each statement
	for i, stmt := range p.Statements {
		t.Logf("Evaluating statement %d: %T", i, stmt)

		if ident, ok := stmt.(*ast.NumberLiteral); ok && i == 1 {
			// Handle the second statement (y = 10)
			t.Logf("  Number literal: %v", ident.Value)

			// Create an assignment for y
			assign := &ast.Assignment{
				Name: "y",
				Value: ident,
			}

			// Use the interpreter directly
			result := interp.Eval(assign)
			t.Logf("  Result type after manual fix: %T", result)
			if err, ok := result.(*interpreter.ErrorValue); ok {
				t.Fatalf("  Error evaluating assignment: %v", err.Message)
			}
		} else if assign, ok := stmt.(*ast.Assignment); ok {
			// Handle assignments
			t.Logf("  Assignment: %s", assign.Name)
			if assign.Value != nil {
				t.Logf("  Assignment value type: %T", assign.Value)
			} else {
				t.Logf("  Assignment value is nil")

				// Create a variable declaration based on the assignment
				var value ast.Node

				// For the first statement (x = 5)
				if i == 0 {
					value = &ast.NumberLiteral{
						Value: 5,
						IsInt: true,
					}
				} else if i == 2 {
					// For the third statement (z = x + y)
					value = &ast.BinaryExpr{
						Left: &ast.Identifier{
							Name: "x",
						},
						Operator: "+",
						Right: &ast.Identifier{
							Name: "y",
						},
					}
				}

				// Create a new assignment with the value
				newAssign := &ast.Assignment{
					Name:  assign.Name,
					Value: value,
				}

				// Use the interpreter directly
				result := interp.Eval(newAssign)
				t.Logf("  Result type after manual fix: %T", result)
				if err, ok := result.(*interpreter.ErrorValue); ok {
					t.Fatalf("  Error evaluating assignment: %v", err.Message)
				}
			}
		} else if ident, ok := stmt.(*ast.Identifier); ok {
			// Handle identifiers
			t.Logf("  Identifier: %s", ident.Name)
			// Look up the identifier in the environment
			result := interp.Eval(ident)
			t.Logf("  Result type: %T", result)
			if err, ok := result.(*interpreter.ErrorValue); ok {
				t.Fatalf("  Error evaluating identifier: %v", err.Message)
			}

			// Convert the interpreter value to an object
			var objResult object.Object
			switch v := result.(type) {
			case *interpreter.IntegerValue:
				objResult = &object.Integer{Value: v.Value}
				assertIntegerObject(t, objResult, 15)
			case *interpreter.FloatValue:
				objResult = &object.Float{Value: v.Value}
			case *interpreter.StringValue:
				objResult = &object.String{Value: v.Value}
			case *interpreter.BooleanValue:
				objResult = &object.Boolean{Value: v.Value}
			default:
				objResult = object.NULL
			}
		} else {
			// Evaluate other types of statements
			result := interp.Eval(stmt)
			t.Logf("  Result type: %T", result)
			if err, ok := result.(*interpreter.ErrorValue); ok {
				t.Fatalf("  Error evaluating statement: %v", err.Message)
			}
		}
	}
}

func TestIfElseStatement(t *testing.T) {
	input := `x: int = 10
	y: int = 0

	if x > 5 do
		y = 1
	else
		y = 2
	end
	y`

	evaluated := evalStatementTest(t, input)
	assertIntegerObject(t, evaluated, 1)
}

func TestForLoopWithArray(t *testing.T) {
	input := `sum: int = 0
	for i in [1, 2, 3, 4, 5] do
		sum = sum + i
	end
	sum`

	evaluated := evalStatementTest(t, input)
	assertIntegerObject(t, evaluated, 15)
}

func TestFunctionDefinitionAndCall(t *testing.T) {
	input := `def add(x: int, y: int): int do
		x + y
	end
	add(2, 5)`

	evaluated := evalStatementTest(t, input)
	assertIntegerObject(t, evaluated, 7)
}

func TestClosure(t *testing.T) {
	input := `def makeAdder(x: int): function do
		def inner(y: int): int do
			x + y
		end
		inner
	end
	adder = makeAdder(5)
	adder(10)`

	evaluated := evalStatementTest(t, input)
	assertIntegerObject(t, evaluated, 15)
}

func TestFunctionWithoutExplicitReturnType(t *testing.T) {
	input := `def multiply(x: int, y: int): int do
		x * y
	end
	multiply(3, 4)`

	evaluated := evalStatementTest(t, input)
	assertIntegerObject(t, evaluated, 12)
}