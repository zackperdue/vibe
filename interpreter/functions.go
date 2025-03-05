package interpreter

import (
	"fmt"

	"github.com/vibe-lang/vibe/ast"
	"github.com/vibe-lang/vibe/types"
)

// Function evaluation methods

// evalFunctionDefinition evaluates a function definition
func (i *Interpreter) evalFunctionDefinition(node *ast.FunctionDef, env *Environment) Value {
	// Parse return type
	var returnType types.Type
	if node.ReturnType != nil {
		returnType = i.parseTypeAnnotation(node.ReturnType)
	} else {
		returnType = types.AnyType
	}

	// Create the function value with parameter types properly processed
	function := &FunctionValue{
		Name:           node.Name,
		Parameters:     node.Parameters, // Use the original parameters
		Body:           node.Body,
		ReturnType:     returnType,
		Env:            env,
	}

	// Add the function to the environment
	env.SetWithType(node.Name, function, function.VibeType())

	// Return the function value instead of nil
	return function
}

// evalCallExpression evaluates a function call expression
func (i *Interpreter) evalCallExpression(node *ast.CallExpr, env *Environment) Value {
	function := i.eval(node.Function, env)
	args := i.evalExpressions(node.Args, env)

	if fn, ok := function.(*FunctionValue); ok {
		// Check arity
		if len(fn.Parameters) != len(args) {
			return &StringValue{Value: fmt.Sprintf(
				"Error: wrong number of arguments: expected %d, got %d",
				len(fn.Parameters), len(args))}
		}

		// Create a new environment for the function
		newEnv := NewEnclosedEnvironment(fn.Env)

		// Bind arguments to parameters
		for paramIdx, param := range fn.Parameters {
			if paramIdx < len(args) {
				// Get the parameter type from the TypeAnnotation
				var paramType types.Type
				if param.Type != nil {
					paramType = i.parseTypeAnnotation(param.Type)
				} else {
					paramType = types.AnyType
				}

				// Type check the argument
				if !types.IsAssignable(args[paramIdx].VibeType(), paramType) {
					return &StringValue{Value: fmt.Sprintf(
						"Type error: Parameter '%s' of function '%s' expects %s, got %s",
						param.Name, fn.Name, paramType.String(), args[paramIdx].VibeType().String())}
				}

				// Bind the parameter
				newEnv.SetWithType(param.Name, args[paramIdx], paramType)
			} else {
				// Missing argument, use nil
				var paramType types.Type
				if param.Type != nil {
					paramType = i.parseTypeAnnotation(param.Type)
				} else {
					paramType = types.AnyType
				}
				newEnv.SetWithType(param.Name, &NilValue{}, paramType)
			}
		}

		// Evaluate the function body
		result := i.evalBlockStatement(fn.Body, newEnv)

		// Unwrap return value, if necessary
		if returnValue, ok := result.(*ReturnValue); ok {
			// Type check the return value
			if !types.IsAssignable(returnValue.Value.VibeType(), fn.ReturnType) {
				return &StringValue{Value: fmt.Sprintf(
					"Type error: Function '%s' returns %s, got %s",
					fn.Name, fn.ReturnType.String(), returnValue.Value.VibeType().String())}
			}
			return returnValue.Value
		}

		// Type check the return value
		if !types.IsAssignable(result.VibeType(), fn.ReturnType) {
			return &StringValue{Value: fmt.Sprintf(
				"Type error: Function '%s' returns %s, got %s",
				fn.Name, fn.ReturnType.String(), result.VibeType().String())}
		}

		return result
	} else if builtin, ok := function.(*BuiltinFunction); ok {
		// Check arity
		if len(args) != len(builtin.ParamTypes) {
			return &StringValue{Value: fmt.Sprintf(
				"Error: wrong number of arguments: function '%s' expects %d, got %d",
				builtin.Name, len(builtin.ParamTypes), len(args))}
		}

		// Type check arguments
		for i, arg := range args {
			if !types.IsAssignable(arg.VibeType(), builtin.ParamTypes[i]) {
				return &StringValue{Value: fmt.Sprintf(
					"Type error: Parameter %d of builtin function '%s' expects %s, got %s",
					i, builtin.Name, builtin.ParamTypes[i].String(), arg.VibeType().String())}
			}
		}

		return builtin.Fn(args)
	}

	return &StringValue{Value: fmt.Sprintf("Not a function: %s", function.Type())}
}

// evalExpressions evaluates a list of expressions
func (i *Interpreter) evalExpressions(
	exps []ast.Node,
	env *Environment,
) []Value {
	var result []Value

	for _, exp := range exps {
		evaluated := i.eval(exp, env)
		result = append(result, evaluated)
	}

	return result
}