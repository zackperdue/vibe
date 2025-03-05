package interpreter

import (
	"fmt"
	"strings"

	"github.com/vibe-lang/vibe/ast"
	"github.com/vibe-lang/vibe/object"
)

// NIL is a singleton nil value
var NIL = &NilValue{}

// Eval evaluates the given AST node using the interpreter and returns an object.Object
func Eval(node ast.Node, env *object.Environment) object.Object {
	// Create a new interpreter
	interpreter := New()

	// Process Program nodes specially to evaluate all statements with the same environment
	if program, ok := node.(*ast.Program); ok {
		var result Value = NIL

		// We'll use a simpler approach - eval all statements in the interpreter
		for i, stmt := range program.Statements {
			// Debug the statement type
			fmt.Printf("Evaluating statement %d of type: %T\n", i, stmt)

			// Special handling for identifiers at the end of the program
			if ident, ok := stmt.(*ast.Identifier); ok {
				fmt.Printf("Looking up identifier: %s\n", ident.Name)
				// Look up the identifier in the environment
				if val, ok := env.Get(ident.Name); ok {
					fmt.Printf("Found value for %s: %v\n", ident.Name, val)
					return val
				} else {
					fmt.Printf("Identifier not found: %s\n", ident.Name)
					return &object.Error{Message: fmt.Sprintf("identifier not found: %s", ident.Name)}
				}
			}

			// Evaluate the statement
			result = interpreter.Eval(stmt)

			// If the result is an error, convert and return it immediately
			if errValue, ok := result.(*ErrorValue); ok {
				return &object.Error{Message: errValue.Message}
			}

			// If the result is a return value, unwrap it and return
			if returnValue, ok := result.(*ReturnValue); ok {
				return valueToObject(returnValue.Value)
			}

			// Update the environment based on the statement type
			if varDecl, ok := stmt.(*ast.VariableDecl); ok {
				fmt.Printf("Setting variable %s in environment\n", varDecl.Name)
				// Set the variable in the environment
				env.Set(varDecl.Name, valueToObject(result))
			} else if assign, ok := stmt.(*ast.Assignment); ok {
				fmt.Printf("Setting assignment %s in environment\n", assign.Name)
				// Set the variable in the environment
				env.Set(assign.Name, valueToObject(result))
			}
		}

		// Convert the final result to an object.Object
		return valueToObject(result)
	} else {
		// For all other nodes, evaluate and return result
		result := interpreter.Eval(node)
		return valueToObject(result)
	}
}

// evalProgram evaluates all statements in a program with the same environment
func evalProgram(interpreter *Interpreter, program *ast.Program, env *object.Environment) object.Object {
	var result Value = &NilValue{}

	for _, stmt := range program.Statements {
		// For each statement, we need to:
		// 1. Create a new program node with just this statement
		singleStmt := &ast.Program{
			Statements: []ast.Node{stmt},
		}

		// 2. Evaluate it with the shared environment
		objResult := Eval(singleStmt, env)

		// 3. For the last statement, this will be our result
		result = objectToValue(objResult)
	}

	return valueToObject(result)
}

// valueToObject converts an interpreter.Value to an object.Object
func valueToObject(value Value) object.Object {
	if value == nil {
		return &object.Error{Message: "Evaluation resulted in nil value"}
	}

	switch value := value.(type) {
	case *IntegerValue:
		return &object.Integer{Value: int64(value.Value)}
	case *FloatValue:
		return &object.Float{Value: value.Value}
	case *StringValue:
		// Check for error messages
		if strings.Contains(value.Value, "identifier not found:") {
			return &object.Error{Message: value.Value}
		}
		if strings.Contains(value.Value, "Expected next token to be") {
			return &object.Error{Message: value.Value}
		}
		if strings.Contains(value.Value, "Index out of range:") {
			return &object.Error{Message: value.Value}
		}
		if strings.Contains(value.Value, "Type error:") {
			return &object.Error{Message: value.Value}
		}
		return &object.String{Value: value.Value}
	case *BooleanValue:
		return &object.Boolean{Value: value.Value}
	case *NilValue:
		return object.NULL
	case *ReturnValue:
		return &object.ReturnValue{Value: valueToObject(value.Value)}
	case *ArrayValue:
		elements := make([]object.Object, len(value.Elements))
		for i, element := range value.Elements {
			elements[i] = valueToObject(element)
		}
		return &object.Array{Elements: elements}
	case *ErrorValue:
		return &object.Error{Message: value.Message}
	default:
		return &object.Error{Message: fmt.Sprintf("Unknown value type: %T", value)}
	}
}

// objectToValue converts an object.Object to an interpreter.Value
func objectToValue(obj object.Object) Value {
	switch obj := obj.(type) {
	case *object.Integer:
		return &IntegerValue{Value: obj.Value}
	case *object.Float:
		return &FloatValue{Value: obj.Value}
	case *object.String:
		return &StringValue{Value: obj.Value}
	case *object.Boolean:
		return &BooleanValue{Value: obj.Value}
	case *object.Array:
		elements := make([]Value, len(obj.Elements))
		for i, element := range obj.Elements {
			elements[i] = objectToValue(element)
		}
		return &ArrayValue{Elements: elements}
	case *object.ReturnValue:
		return &ReturnValue{Value: objectToValue(obj.Value)}
	case *object.Error:
		return &ErrorValue{Message: obj.Message}
	case *object.Null:
		return &NilValue{}
	default:
		return &ErrorValue{Message: fmt.Sprintf("Unknown object type: %T", obj)}
	}
}