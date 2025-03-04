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

	// Convert object.Environment to interpreter.Environment
	// Note: In a real implementation, we would need to properly convert all variables
	// from the object.Environment to the interpreter.Environment
	// For now, we'll just use the default environment
	// interpreterEnv := interpreter.env

	// Evaluate the node
	result := interpreter.Eval(node)

	// Convert the result to an object.Object
	return valueToObject(result)
}

// valueToObject converts an interpreter.Value to an object.Object
func valueToObject(value Value) object.Object {
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
	default:
		return &object.Error{Message: fmt.Sprintf("Unknown value type: %T", value)}
	}
}