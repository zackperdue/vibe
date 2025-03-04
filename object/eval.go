package object

import (
	"fmt"

	"github.com/vibe-lang/vibe/ast"
	"github.com/vibe-lang/vibe/lexer"
	"github.com/vibe-lang/vibe/parser"
)

var (
	NULL  = &Null{}
	TRUE  = &Boolean{Value: true}
	FALSE = &Boolean{Value: false}
)

// Eval evaluates the input string and returns an Object
func Eval(input string) Object {
	l := lexer.New(input)
	program, errors := parser.Parse(l)

	if len(errors) != 0 {
		errorMessages := ""
		for _, msg := range errors {
			errorMessages += msg + "\n"
		}
		return &Error{Message: errorMessages}
	}

	// Create a new environment
	env := NewEnvironment()

	// Evaluate each statement in the program
	var result Object = NULL
	for _, stmt := range program.Statements {
		result = evalNode(stmt, env)

		// If it's a return value, unwrap it
		if returnValue, ok := result.(*ReturnValue); ok {
			return returnValue.Value
		}

		// If it's an error, stop evaluation
		if isError(result) {
			return result
		}
	}

	return result
}

// evalNode evaluates an AST node
func evalNode(node ast.Node, env *Environment) Object {
	if node == nil {
		// Return NULL for nil nodes instead of an error
		// This allows the program to continue rather than fail
		return NULL
	}

	switch node := node.(type) {
	// Statements
	case *ast.Program:
		return evalProgram(node, env)
	case *ast.BlockStmt:
		return evalBlockStatement(node, env)
	case *ast.ReturnStmt:
		val := evalNode(node.Value, env)
		if isError(val) {
			return val
		}
		return &ReturnValue{Value: val}
	case *ast.IfStmt:
		return evalIfStatement(node, env)
	case *ast.ForStmt:
		return evalForStatement(node, env)
	case *ast.Assignment:
		val := evalNode(node.Value, env)
		if isError(val) {
			return val
		}
		env.Set(node.Name, val)
		return val

	// Expressions
	case *ast.NumberLiteral:
		if node.IsInt {
			return &Integer{Value: int64(node.Value)}
		}
		return &Float{Value: node.Value}
	case *ast.StringLiteral:
		return &String{Value: node.Value}
	case *ast.BooleanLiteral:
		return nativeBoolToBooleanObject(node.Value)
	case *ast.NilLiteral:
		return NULL
	case *ast.ArrayLiteral:
		elements := evalExpressions(node.Elements, env)
		if len(elements) == 1 && isError(elements[0]) {
			return elements[0]
		}
		return &Array{Elements: elements}
	case *ast.Identifier:
		return evalIdentifier(node, env)
	case *ast.BinaryExpr:
		left := evalNode(node.Left, env)
		if isError(left) {
			return left
		}
		right := evalNode(node.Right, env)
		if isError(right) {
			return right
		}
		return evalBinaryExpression(node.Operator, left, right)
	case *ast.UnaryExpr:
		right := evalNode(node.Right, env)
		if isError(right) {
			return right
		}
		return evalUnaryExpression(node.Operator, right)
	case *ast.IndexExpr:
		left := evalNode(node.Array, env)
		if isError(left) {
			return left
		}
		index := evalNode(node.Index, env)
		if isError(index) {
			return index
		}
		return evalIndexExpression(left, index)
	}

	return &Error{Message: fmt.Sprintf("Unknown node type: %T", node)}
}

// evalProgram evaluates a program
func evalProgram(program *ast.Program, env *Environment) Object {
	var result Object = NULL

	for _, statement := range program.Statements {
		result = evalNode(statement, env)

		switch result := result.(type) {
		case *ReturnValue:
			return result.Value
		case *Error:
			return result
		}
	}

	return result
}

// evalBlockStatement evaluates a block statement
func evalBlockStatement(block *ast.BlockStmt, env *Environment) Object {
	var result Object = NULL

	for _, statement := range block.Statements {
		result = evalNode(statement, env)

		if result != nil {
			rt := result.Type()
			if rt == RETURN_VALUE_OBJ || rt == ERROR_OBJ {
				return result
			}
		}
	}

	return result
}

// evalIfStatement evaluates an if statement
func evalIfStatement(ifStmt *ast.IfStmt, env *Environment) Object {
	condition := evalNode(ifStmt.Condition, env)
	if isError(condition) {
		return condition
	}

	if isTruthy(condition) {
		return evalBlockStatement(ifStmt.Consequence, env)
	} else if ifStmt.Alternative != nil {
		return evalBlockStatement(ifStmt.Alternative, env)
	} else {
		return NULL
	}
}

// evalForStatement evaluates a for statement
func evalForStatement(forStmt *ast.ForStmt, env *Environment) Object {
	iterable := evalNode(forStmt.Iterable, env)
	if isError(iterable) {
		return iterable
	}

	if array, ok := iterable.(*Array); ok {
		return evalForLoop(env, forStmt.Iterator, array, forStmt.Body)
	}

	return &Error{Message: fmt.Sprintf("for loop iterable must be an array, got %s", iterable.Type())}
}

// evalForLoop evaluates a for loop
func evalForLoop(env *Environment, iterator string, iterable *Array, body ast.Node) Object {
	// First, ensure the iterator variable exists in the outer environment
	// and initialize it to NULL
	env.Set(iterator, NULL)

	// Create a new environment for the loop that inherits from the outer environment
	// This way, the loop body can access variables from the outer scope
	// but variables defined inside the loop won't leak out
	loopEnv := NewEnclosedEnvironment(env)

	// No elements to iterate over, return NULL
	if len(iterable.Elements) == 0 {
		// Make sure the iterator is defined in the loop environment too
		loopEnv.Set(iterator, NULL)
		return NULL
	}

	// Iterate over each element
	for _, item := range iterable.Elements {
		// Set the iterator variable in both environments
		// This is crucial - the outer environment needs it for access after the loop
		// and the loop environment needs it for access during the loop
		env.Set(iterator, item)
		loopEnv.Set(iterator, item)

		// Evaluate the loop body
		var result Object
		if blockStmt, ok := body.(*ast.BlockStmt); ok {
			result = evalBlockStatement(blockStmt, loopEnv)
		} else {
			result = evalNode(body, loopEnv)
		}

		// Handle return or error from within the loop body
		if result != nil && (result.Type() == RETURN_VALUE_OBJ || result.Type() == ERROR_OBJ) {
			return result
		}
	}

	// Loop completed successfully
	return NULL
}

// evalExpressions evaluates a list of expressions
func evalExpressions(exps []ast.Node, env *Environment) []Object {
	var result []Object

	for _, e := range exps {
		evaluated := evalNode(e, env)
		if isError(evaluated) {
			return []Object{evaluated}
		}
		result = append(result, evaluated)
	}

	return result
}

// evalIdentifier evaluates an identifier
func evalIdentifier(node *ast.Identifier, env *Environment) Object {
	if val, ok := env.Get(node.Name); ok {
		return val
	}

	return &Error{Message: fmt.Sprintf("identifier not found: %s", node.Name)}
}

// evalBinaryExpression evaluates a binary expression
func evalBinaryExpression(operator string, left, right Object) Object {
	switch {
	case left.Type() == INTEGER_OBJ && right.Type() == INTEGER_OBJ:
		return evalIntegerBinaryExpression(operator, left, right)
	case left.Type() == STRING_OBJ && right.Type() == STRING_OBJ:
		return evalStringBinaryExpression(operator, left, right)
	case operator == "==":
		return nativeBoolToBooleanObject(left == right)
	case operator == "!=":
		return nativeBoolToBooleanObject(left != right)
	case left.Type() != right.Type():
		return &Error{Message: fmt.Sprintf("type mismatch: %s %s %s", left.Type(), operator, right.Type())}
	default:
		return &Error{Message: fmt.Sprintf("unknown operator: %s %s %s", left.Type(), operator, right.Type())}
	}
}

// evalIntegerBinaryExpression evaluates a binary expression with integer operands
func evalIntegerBinaryExpression(operator string, left, right Object) Object {
	leftVal := left.(*Integer).Value
	rightVal := right.(*Integer).Value

	switch operator {
	case "+":
		return &Integer{Value: leftVal + rightVal}
	case "-":
		return &Integer{Value: leftVal - rightVal}
	case "*":
		return &Integer{Value: leftVal * rightVal}
	case "/":
		return &Integer{Value: leftVal / rightVal}
	case "<":
		return nativeBoolToBooleanObject(leftVal < rightVal)
	case ">":
		return nativeBoolToBooleanObject(leftVal > rightVal)
	case "==":
		return nativeBoolToBooleanObject(leftVal == rightVal)
	case "!=":
		return nativeBoolToBooleanObject(leftVal != rightVal)
	default:
		return &Error{Message: fmt.Sprintf("unknown operator: %s %s %s", left.Type(), operator, right.Type())}
	}
}

// evalStringBinaryExpression evaluates a binary expression with string operands
func evalStringBinaryExpression(operator string, left, right Object) Object {
	if operator != "+" {
		return &Error{Message: fmt.Sprintf("unknown operator: %s %s %s", left.Type(), operator, right.Type())}
	}

	leftVal := left.(*String).Value
	rightVal := right.(*String).Value
	return &String{Value: leftVal + rightVal}
}

// evalUnaryExpression evaluates a unary expression
func evalUnaryExpression(operator string, right Object) Object {
	switch operator {
	case "!":
		return evalBangOperatorExpression(right)
	case "-":
		return evalMinusPrefixOperatorExpression(right)
	default:
		return &Error{Message: fmt.Sprintf("unknown operator: %s%s", operator, right.Type())}
	}
}

// evalBangOperatorExpression evaluates a bang operator expression
func evalBangOperatorExpression(right Object) Object {
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

// evalMinusPrefixOperatorExpression evaluates a minus prefix operator expression
func evalMinusPrefixOperatorExpression(right Object) Object {
	if right.Type() != INTEGER_OBJ {
		return &Error{Message: fmt.Sprintf("unknown operator: -%s", right.Type())}
	}

	value := right.(*Integer).Value
	return &Integer{Value: -value}
}

// evalIndexExpression evaluates an index expression
func evalIndexExpression(left, index Object) Object {
	switch {
	case left.Type() == ARRAY_OBJ && index.Type() == INTEGER_OBJ:
		return evalArrayIndexExpression(left, index)
	default:
		return &Error{Message: fmt.Sprintf("index operator not supported: %s", left.Type())}
	}
}

// evalArrayIndexExpression evaluates an array index expression
func evalArrayIndexExpression(array, index Object) Object {
	arrayObject := array.(*Array)
	idx := index.(*Integer).Value
	max := int64(len(arrayObject.Elements) - 1)

	if idx < 0 || idx > max {
		return NULL
	}

	return arrayObject.Elements[idx]
}

// Helper functions

// isTruthy determines if an object is truthy
func isTruthy(obj Object) bool {
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

// nativeBoolToBooleanObject converts a native bool to a Boolean object
func nativeBoolToBooleanObject(input bool) *Boolean {
	if input {
		return TRUE
	}
	return FALSE
}

// isError checks if an object is an error
func isError(obj Object) bool {
	if obj != nil {
		return obj.Type() == ERROR_OBJ
	}
	return false
}