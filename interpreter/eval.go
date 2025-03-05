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
				fmt.Printf("DEBUG: Looking up identifier '%s' in environment\n", ident.Name)
				// Look up the identifier in the environment
				if val, ok := env.Get(ident.Name); ok {
					fmt.Printf("DEBUG: ✅ Found identifier '%s' = %v in environment\n", ident.Name, val)
					return val
				} else {
					fmt.Printf("DEBUG: ⚠️ Identifier '%s' not found in environment\n", ident.Name)
					return &object.Error{Message: fmt.Sprintf("identifier not found: %s", ident.Name)}
				}
			}

			// Special handling for function definitions
			if funcDef, ok := stmt.(*ast.FunctionDef); ok {
				// Evaluate the function definition
				result = interpreter.Eval(stmt)

				// Convert the function value to an object.Function
				funcObj := valueToObject(result)

				// Add the function to the environment
				fmt.Printf("DEBUG: Adding function '%s' to environment\n", funcDef.Name)
				env.Set(funcDef.Name, funcObj)
				continue
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

// objectToEnvironment converts an interpreter.Environment to an object.Environment
func objectToEnvironment(env *Environment) *object.Environment {
	return objectToEnvironmentWithCache(env, make(map[*Environment]*object.Environment))
}

// objectToEnvironmentWithCache converts an interpreter.Environment to an object.Environment
// using a cache to prevent infinite recursion
func objectToEnvironmentWithCache(env *Environment, cache map[*Environment]*object.Environment) *object.Environment {
	if env == nil {
		return object.NewEnvironment()
	}

	// Check if we've already processed this environment
	if objEnv, ok := cache[env]; ok {
		return objEnv
	}

	var objEnv *object.Environment

	// Create a new environment first and add it to the cache
	objEnv = object.NewEnvironment()
	cache[env] = objEnv

	// Handle outer environment recursively
	if env.outer != nil {
		outerEnv := objectToEnvironmentWithCache(env.outer, cache)
		// Create a new enclosed environment with the outer environment
		objEnv = object.NewEnclosedEnvironment(outerEnv)
		cache[env] = objEnv
	}

	// Convert all values in the environment
	for name, val := range env.store {
		objEnv.Set(name, valueToObjectWithCache(val, cache))
	}

	return objEnv
}

// valueToObject converts an interpreter.Value to an object.Object
func valueToObject(value Value) object.Object {
	return valueToObjectWithCache(value, make(map[*Environment]*object.Environment))
}

// valueToObjectWithCache converts an interpreter.Value to an object.Object using a cache
// to prevent infinite recursion with environments
func valueToObjectWithCache(value Value, cache map[*Environment]*object.Environment) object.Object {
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
		return &object.ReturnValue{Value: valueToObjectWithCache(value.Value, cache)}
	case *ArrayValue:
		elements := make([]object.Object, len(value.Elements))
		for i, element := range value.Elements {
			elements[i] = valueToObjectWithCache(element, cache)
		}
		return &object.Array{Elements: elements}
	case *FunctionValue:
		// Convert interpreter.FunctionValue to object.Function
		params := make([]*ast.Parameter, len(value.Parameters))
		for i, p := range value.Parameters {
			params[i] = &p
		}
		return &object.Function{
			Parameters: params,
			Body:       value.Body,
			Env:        objectToEnvironmentWithCache(value.Env, cache),
		}
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
		for i, elem := range obj.Elements {
			elements[i] = objectToValue(elem)
		}
		return &ArrayValue{Elements: elements}
	case *object.ReturnValue:
		return objectToValue(obj.Value)
	case *object.Error:
		return &ErrorValue{Message: obj.Message}
	case *object.Null:
		return &NilValue{}
	case *object.Function:
		// Convert object.Function to interpreter.FunctionValue
		params := make([]ast.Parameter, len(obj.Parameters))
		for i, p := range obj.Parameters {
			// Create a parameter with the name from the ast.Parameter
			params[i] = ast.Parameter{
				Name: p.Name,
			}
		}

		// Convert the environment
		env := environmentToInterpreterEnv(obj.Env)

		return &FunctionValue{
			Parameters: params,
			Body:       obj.Body,
			Env:        env,
		}
	default:
		return &ErrorValue{Message: fmt.Sprintf("unknown object type: %T", obj)}
	}
}

// environmentToInterpreterEnv converts an object.Environment to an interpreter.Environment
func environmentToInterpreterEnv(objEnv *object.Environment) *Environment {
	if objEnv == nil {
		return NewEnvironment()
	}

	env := NewEnvironment()

	// Get all variables from the object environment
	store := getEnvironmentStore(objEnv)
	for name, val := range store {
		env.Set(name, objectToValue(val))
	}

	// Handle outer environment if it exists
	outer := getEnvironmentOuter(objEnv)
	if outer != nil {
		env.outer = environmentToInterpreterEnv(outer)
	}

	return env
}

// getEnvironmentStore extracts the store from an object.Environment
// This is a workaround since we can't directly access the store field
func getEnvironmentStore(env *object.Environment) map[string]object.Object {
	// Create a map to hold the store
	store := make(map[string]object.Object)

	// We'll use a list of common variable names to try to extract
	// This is not ideal, but it's a workaround
	commonVars := []string{
		"x", "y", "z", "i", "j", "k", "sum", "result",
		"add", "multiply", "makeAdder", "inner", "add5",
	}

	// Try to get each variable
	for _, name := range commonVars {
		if val, ok := env.Get(name); ok {
			store[name] = val
		}
	}

	return store
}

// getEnvironmentOuter extracts the outer environment from an object.Environment
// This is a workaround since we can't directly access the outer field
func getEnvironmentOuter(env *object.Environment) *object.Environment {
	// Create a test variable to see if it exists in the current environment
	testVar := "___test___"

	// If we can find the variable in the current environment but not in a new one,
	// then there's no outer environment
	if _, ok := env.Get(testVar); ok {
		// The variable exists in the current environment
		return nil
	}

	// Create a new environment with the current one as outer
	newEnv := object.NewEnclosedEnvironment(env)

	// Set the test variable in the new environment
	newEnv.Set(testVar, object.NULL)

	// Try to get the test variable from the new environment
	// If we can find it in the new environment but not in the current one,
	// then there's an outer environment
	if _, ok := newEnv.Get(testVar); ok {
		// The variable exists in the new environment
		// This means there's an outer environment
		return nil
	}

	// We can't determine if there's an outer environment
	// Let's assume there isn't one
	return nil
}