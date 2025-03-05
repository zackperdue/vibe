package interpreter

import (
	"fmt"

	"github.com/vibe-lang/vibe/ast"
	"github.com/vibe-lang/vibe/object"
)

// NIL is a singleton nil value
var NIL = &NilValue{}

// Eval evaluates the given AST node and returns an object.Object
func Eval(node ast.Node, env *object.Environment) object.Object {
	switch node := node.(type) {
	// Program is a special case - we need to evaluate all statements
	case *ast.Program:
		fmt.Println("Evaluating program with", len(node.Statements), "statements")
		return evalProgram(node, env)

	// Statements
	case *ast.BlockStmt:
		return evalBlockStatement(node, env)
	case *ast.ExpressionStatement:
		return Eval(node.Expression, env)
	case *ast.VariableDecl:
		val := Eval(node.Value, env)
		if isError(val) {
			return val
		}
		env.Set(node.Name, val)
		return val
	case *ast.Assignment:
		val := Eval(node.Value, env)
		if isError(val) {
			return val
		}
		env.Set(node.Name, val)
		return val
	case *ast.ReturnStmt:
		val := Eval(node.Value, env)
		if isError(val) {
			return val
		}
		return &object.ReturnValue{Value: val}
	case *ast.IfStmt:
		return evalIfStatement(node, env)
	case *ast.ForStmt:
		return evalForStatement(node, env)

	// Function definition
	case *ast.FunctionDef:
		fmt.Println("Evaluating function definition:", node.Name)
		params := make([]*ast.Parameter, len(node.Parameters))
		for i, p := range node.Parameters {
			params[i] = &ast.Parameter{
				Name: p.Name,
				Type: p.Type,
			}
		}
		body := node.Body
		function := &object.Function{
			Parameters: params,
			Body:       body,
			Env:        env,
		}
		// Store the function in the environment
		env.Set(node.Name, function)
		fmt.Println("Added function", node.Name, "to environment")
		return function

	// Expressions
	case *ast.NumberLiteral:
		if node.IsInt {
			return &object.Integer{Value: int64(node.Value)}
		}
		return &object.Float{Value: node.Value}
	case *ast.StringLiteral:
		return &object.String{Value: node.Value}
	case *ast.BooleanLiteral:
		return nativeBoolToBooleanObject(node.Value)
	case *ast.BinaryExpr:
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return evalBinaryExpression(node.Operator, left, right)
	case *ast.Identifier:
		fmt.Println("Looking up identifier:", node.Name)
		return evalIdentifier(node, env)
	case *ast.ArrayLiteral:
		elements := evalExpressions(node.Elements, env)
		if len(elements) == 1 && isError(elements[0]) {
			return elements[0]
		}
		return &object.Array{Elements: elements}
	case *ast.IndexExpr:
		left := Eval(node.Array, env)
		if isError(left) {
			return left
		}
		index := Eval(node.Index, env)
		if isError(index) {
			return index
		}
		return evalIndexExpression(left, index)
	case *ast.CallExpr:
		fmt.Println("Evaluating call expression for function:", node.Function)
		function := Eval(node.Function, env)
		if isError(function) {
			return function
		}
		args := evalExpressions(node.Args, env)
		if len(args) == 1 && isError(args[0]) {
			return args[0]
		}
		return applyFunction(function, args)
	case *ast.UnaryExpr:
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return evalPrefixExpression(node.Operator, right)
	}

	return nil
}

// evalProgram evaluates all statements in a program with the same environment
func evalProgram(program *ast.Program, env *object.Environment) object.Object {
	var result object.Object

	for _, statement := range program.Statements {
		fmt.Println("Evaluating statement:", statement)
		result = Eval(statement, env)

		switch result := result.(type) {
		case *object.ReturnValue:
			return result.Value
		case *object.Error:
			return result
		}
	}

	return result
}

// evalBlockStatement evaluates all statements in a block with the same environment
func evalBlockStatement(block *ast.BlockStmt, env *object.Environment) object.Object {
	var result object.Object

	for _, statement := range block.Statements {
		result = Eval(statement, env)

		if result != nil {
			rt := result.Type()
			if rt == object.RETURN_VALUE_OBJ || rt == object.ERROR_OBJ {
				return result
			}
		}
	}

	return result
}

// evalExpressions evaluates a list of expressions with the same environment
func evalExpressions(exps []ast.Node, env *object.Environment) []object.Object {
	var result []object.Object

	for _, e := range exps {
		evaluated := Eval(e, env)
		if isError(evaluated) {
			return []object.Object{evaluated}
		}
		result = append(result, evaluated)
	}

	return result
}

// evalPrefixExpression evaluates a prefix expression with the same environment
func evalPrefixExpression(operator string, right object.Object) object.Object {
	switch operator {
	case "!":
		return evalBangOperatorExpression(right)
	case "-":
		return evalMinusPrefixOperatorExpression(right)
	default:
		return newError("unknown operator: %s%s", operator, right.Type())
	}
}

// evalBangOperatorExpression evaluates a bang operator expression with the same environment
func evalBangOperatorExpression(right object.Object) object.Object {
	switch right {
	case TRUE:
		return FALSE
	case FALSE:
		return TRUE
	case NULL:
		return TRUE
	default:
		return FALSE
	}
}

// evalMinusPrefixOperatorExpression evaluates a minus prefix operator expression with the same environment
func evalMinusPrefixOperatorExpression(right object.Object) object.Object {
	if right.Type() != object.INTEGER_OBJ {
		return newError("unknown operator: -%s", right.Type())
	}

	value := right.(*object.Integer).Value
	return &object.Integer{Value: -value}
}

// evalBinaryExpression evaluates a binary expression with the same environment
func evalBinaryExpression(operator string, left, right object.Object) object.Object {
	switch {
	case left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ:
		return evalIntegerBinaryExpression(operator, left, right)
	case left.Type() == object.STRING_OBJ && right.Type() == object.STRING_OBJ:
		return evalStringBinaryExpression(operator, left, right)
	case operator == "==":
		return nativeBoolToBooleanObject(left == right)
	case operator == "!=":
		return nativeBoolToBooleanObject(left != right)
	case left.Type() != right.Type():
		return newError("type mismatch: %s %s %s", left.Type(), operator, right.Type())
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

// evalIntegerBinaryExpression evaluates an integer binary expression with the same environment
func evalIntegerBinaryExpression(operator string, left, right object.Object) object.Object {
	leftVal := left.(*object.Integer).Value
	rightVal := right.(*object.Integer).Value

	switch operator {
	case "+":
		return &object.Integer{Value: leftVal + rightVal}
	case "-":
		return &object.Integer{Value: leftVal - rightVal}
	case "*":
		return &object.Integer{Value: leftVal * rightVal}
	case "/":
		return &object.Integer{Value: leftVal / rightVal}
	case "<":
		return nativeBoolToBooleanObject(leftVal < rightVal)
	case ">":
		return nativeBoolToBooleanObject(leftVal > rightVal)
	case "==":
		return nativeBoolToBooleanObject(leftVal == rightVal)
	case "!=":
		return nativeBoolToBooleanObject(leftVal != rightVal)
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

// evalStringBinaryExpression evaluates a string binary expression with the same environment
func evalStringBinaryExpression(operator string, left, right object.Object) object.Object {
	if operator != "+" {
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}

	leftVal := left.(*object.String).Value
	rightVal := right.(*object.String).Value
	return &object.String{Value: leftVal + rightVal}
}

// evalIfStatement evaluates an if statement with the same environment
func evalIfStatement(ie *ast.IfStmt, env *object.Environment) object.Object {
	condition := Eval(ie.Condition, env)
	if isError(condition) {
		return condition
	}

	if isObjectTruthy(condition) {
		return Eval(ie.Consequence, env)
	} else if ie.Alternative != nil {
		return Eval(ie.Alternative, env)
	} else {
		return NULL
	}
}

// isObjectTruthy checks if an object is truthy
func isObjectTruthy(obj object.Object) bool {
	switch obj {
	case NULL:
		return false
	case TRUE:
		return true
	case FALSE:
		return false
	default:
		return true
	}
}

// evalForStatement evaluates a for statement with the same environment
func evalForStatement(fs *ast.ForStmt, env *object.Environment) object.Object {
	fmt.Println("Evaluating for statement with iterator:", fs.Iterator)

	// Evaluate the iterable
	iterable := Eval(fs.Iterable, env)
	if isError(iterable) {
		return iterable
	}

	// Check if the iterable is an array
	arr, ok := iterable.(*object.Array)
	if !ok {
		return newError("for loop iterable is not an array: got %s", iterable.Type())
	}

	// Iterate over the array
	var result object.Object = NULL
	for _, element := range arr.Elements {
		// Set the iterator variable in the current environment
		// This is important - we want to update the current environment, not create a new one
		env.Set(fs.Iterator, element)
		fmt.Printf("Setting iterator %s = %v in environment\n", fs.Iterator, element)

		// Evaluate the body in the current environment
		result = Eval(fs.Body, env)

		// Check for return or error
		if result != nil {
			rt := result.Type()
			if rt == object.RETURN_VALUE_OBJ || rt == object.ERROR_OBJ {
				return result
			}
		}
	}

	return result
}

// evalIdentifier evaluates an identifier with the same environment
func evalIdentifier(node *ast.Identifier, env *object.Environment) object.Object {
	if val, ok := env.Get(node.Name); ok {
		fmt.Println("Found identifier:", node.Name, "with value:", val)
		return val
	}

	fmt.Println("Identifier not found:", node.Name)
	return newError("identifier not found: %s", node.Name)
}

// evalIndexExpression evaluates an index expression with the same environment
func evalIndexExpression(left, index object.Object) object.Object {
	switch {
	case left.Type() == object.ARRAY_OBJ && index.Type() == object.INTEGER_OBJ:
		return evalArrayIndexExpression(left, index)
	default:
		return newError("index operator not supported: %s", left.Type())
	}
}

// evalArrayIndexExpression evaluates an array index expression with the same environment
func evalArrayIndexExpression(array, index object.Object) object.Object {
	arrayObject := array.(*object.Array)
	idx := index.(*object.Integer).Value
	max := int64(len(arrayObject.Elements) - 1)

	if idx < 0 || idx > max {
		return NULL
	}

	return arrayObject.Elements[idx]
}

// applyFunction applies a function to a list of arguments with the same environment
func applyFunction(fn object.Object, args []object.Object) object.Object {
	function, ok := fn.(*object.Function)
	if !ok {
		return newError("not a function: %s", fn.Type())
	}

	extendedEnv := extendFunctionEnv(function, args)
	evaluated := Eval(function.Body, extendedEnv)
	return unwrapReturnValue(evaluated)
}

// extendFunctionEnv extends a function environment with a list of arguments
func extendFunctionEnv(fn *object.Function, args []object.Object) *object.Environment {
	env := object.NewEnclosedEnvironment(fn.Env)

	for paramIdx, param := range fn.Parameters {
		env.Set(param.Name, args[paramIdx])
	}

	return env
}

// unwrapReturnValue unwraps a return value from an evaluated object
func unwrapReturnValue(obj object.Object) object.Object {
	if returnValue, ok := obj.(*object.ReturnValue); ok {
		return returnValue.Value
	}

	return obj
}

// newError creates a new error object with the given format and arguments
func newError(format string, a ...interface{}) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, a...)}
}

// isError checks if an object is an error
func isError(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.ERROR_OBJ
	}
	return false
}

// Helper variables for boolean objects
var (
	NULL  = &object.Null{}
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
)

// nativeBoolToBooleanObject converts a native boolean to a boolean object
func nativeBoolToBooleanObject(input bool) *object.Boolean {
	if input {
		return TRUE
	}
	return FALSE
}