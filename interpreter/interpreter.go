package interpreter

import (
	"fmt"
	"strconv"

	"github.com/vibe-lang/vibe/ast"
	"github.com/vibe-lang/vibe/types"
)

// Interpreter is the main interpreter for the Vibe language
type Interpreter struct {
	env           *Environment
	currentModule string
}

// New creates a new interpreter
func New() *Interpreter {
	env := NewEnvironment()
	i := &Interpreter{env: env, currentModule: "main"}

	// Register built-in functions
	env.RegisterBuiltin("print", func(args []Value) Value {
		for _, arg := range args {
			fmt.Println(arg.Inspect())
		}
		return &NilValue{}
	}, []types.Type{types.AnyType}, types.NilType)

	env.RegisterBuiltin("len", func(args []Value) Value {
		if len(args) != 1 {
			return &StringValue{Value: "Error: wrong number of arguments: expected 1, got " + strconv.Itoa(len(args))}
		}

		switch arg := args[0].(type) {
		case *StringValue:
			return &IntegerValue{Value: int64(len(arg.Value))}
		case *ArrayValue:
			return &IntegerValue{Value: int64(len(arg.Elements))}
		default:
			return &StringValue{Value: "Type error: len requires a string or array argument"}
		}
	}, []types.Type{types.AnyType}, types.IntType)

	env.RegisterBuiltin("int", func(args []Value) Value {
		if len(args) != 1 {
			return &StringValue{Value: "Error: wrong number of arguments: expected 1, got " + strconv.Itoa(len(args))}
		}

		switch arg := args[0].(type) {
		case *StringValue:
			i, err := strconv.ParseInt(arg.Value, 10, 64)
			if err != nil {
				return &StringValue{Value: "Type error: cannot convert string to int"}
			}
			return &IntegerValue{Value: i}
		case *FloatValue:
			return &IntegerValue{Value: int64(arg.Value)}
		case *IntegerValue:
			return arg
		default:
			return &StringValue{Value: "Type error: int requires a string, float, or int argument"}
		}
	}, []types.Type{types.AnyType}, types.IntType)

	env.RegisterBuiltin("float", func(args []Value) Value {
		if len(args) != 1 {
			return &StringValue{Value: "Error: wrong number of arguments: expected 1, got " + strconv.Itoa(len(args))}
		}

		switch arg := args[0].(type) {
		case *StringValue:
			f, err := strconv.ParseFloat(arg.Value, 64)
			if err != nil {
				return &StringValue{Value: "Type error: cannot convert string to float"}
			}
			return &FloatValue{Value: f}
		case *IntegerValue:
			return &FloatValue{Value: float64(arg.Value)}
		case *FloatValue:
			return arg
		default:
			return &StringValue{Value: "Type error: float requires a string, int, or float argument"}
		}
	}, []types.Type{types.AnyType}, types.FloatType)

	env.RegisterBuiltin("str", func(args []Value) Value {
		if len(args) != 1 {
			return &StringValue{Value: "Error: wrong number of arguments: expected 1, got " + strconv.Itoa(len(args))}
		}

		return &StringValue{Value: valueToString(args[0])}
	}, []types.Type{types.AnyType}, types.StringType)

	return i
}

// valueToString converts a Value to a string
func valueToString(v Value) string {
	switch v := v.(type) {
	case *StringValue:
		return v.Value
	case *IntegerValue:
		return strconv.FormatInt(v.Value, 10)
	case *FloatValue:
		return strconv.FormatFloat(v.Value, 'f', -1, 64)
	case *BooleanValue:
		return strconv.FormatBool(v.Value)
	case *NilValue:
		return "nil"
	default:
		return v.Inspect()
	}
}

// Eval evaluates an AST node
func (i *Interpreter) Eval(node ast.Node) Value {
	return i.eval(node, i.env)
}

// eval evaluates an AST node in a given environment
func (i *Interpreter) eval(node ast.Node, env *Environment) Value {
	if node == nil {
		return &ErrorValue{Message: "Cannot evaluate nil node"}
	}

	switch node := node.(type) {
	case *ast.Program:
		return i.evalProgram(node, env)
	case *ast.BlockStmt:
		return i.evalBlockStatement(node, env)
	case *ast.ExpressionStatement:
		return i.eval(node.Expression, env)
	case *ast.NumberLiteral:
		if node.IsInt {
			return &IntegerValue{Value: int64(node.Value)}
		}
		return &FloatValue{Value: node.Value}
	case *ast.StringLiteral:
		return &StringValue{Value: node.Value}
	case *ast.BooleanLiteral:
		return &BooleanValue{Value: node.Value}
	case *ast.NilLiteral:
		return &NilValue{}
	case *ast.Identifier:
		return i.evalIdentifier(node, env)
	case *ast.VariableDecl:
		return i.evalVariableDeclaration(node, env)
	case *ast.ReturnStmt:
		return i.evalReturnStatement(node, env)
	case *ast.IfStmt:
		return i.evalIfStatement(node, env)
	case *ast.ForStmt:
		return i.evalForStatement(node, env)
	case *ast.ClassDef:
		return i.evalClassDefinition(node, env)
	case *ast.RequireStmt:
		return i.evalImportStatement(node, env)
	case *ast.UnaryExpr:
		return i.evalPrefixExpression(node, env)
	case *ast.BinaryExpr:
		return i.evalInfixExpression(node, env)
	case *ast.Assignment:
		return i.evalAssignmentExpression(node, env)
	case *ast.ArrayLiteral:
		return i.evalArrayLiteral(node, env)
	case *ast.IndexExpr:
		return i.evalIndexExpression(node, env)
	case *ast.DotExpr:
		return i.evalMemberExpression(node, env)
	case *ast.MethodCall:
		return i.evalMethodCallExpression(node, env)
	case *ast.FunctionDef:
		return i.evalFunctionDefinition(node, env)
	case *ast.CallExpr:
		function := i.eval(node.Function, env)
		if _, isError := function.(*ErrorValue); isError {
			return function
		}

		args := []Value{}
		for _, arg := range node.Args {
			evaluated := i.eval(arg, env)
			if _, isError := evaluated.(*ErrorValue); isError {
				return evaluated
			}
			args = append(args, evaluated)
		}

		// Handle built-in functions
		if builtin, ok := function.(*BuiltinFunction); ok {
			return builtin.Fn(args)
		}

		// Handle user-defined functions
		if fn, ok := function.(*FunctionValue); ok {
			// Create a new environment for the function
			functionEnv := NewEnclosedEnvironment(fn.Env)

			// Bind arguments to parameters
			for i, param := range fn.Parameters {
				if i < len(args) {
					functionEnv.Set(param.Name, args[i])
				} else {
					// Missing argument, use nil
					functionEnv.Set(param.Name, &NilValue{})
				}
			}

			// Evaluate the function body
			result := i.evalBlockStatement(fn.Body, functionEnv)

			// Unwrap return value
			if returnValue, ok := result.(*ReturnValue); ok {
				return returnValue.Value
			}

			return result
		}

		return &ErrorValue{Message: fmt.Sprintf("Not a function: %s", function.Type())}
	default:
		return &ErrorValue{Message: fmt.Sprintf("Unknown node type: %T", node)}
	}
}

// evalProgram evaluates a program
func (i *Interpreter) evalProgram(program *ast.Program, env *Environment) Value {
	var result Value
	result = &NilValue{}

	for _, statement := range program.Statements {
		result = i.eval(statement, env)

		// If we hit a return statement, unwrap it and return the value
		if returnValue, ok := result.(*ReturnValue); ok {
			return returnValue.Value
		}
	}

	return result
}

// evalBlockStatement evaluates a block statement
func (i *Interpreter) evalBlockStatement(block *ast.BlockStmt, env *Environment) Value {
	var result Value
	result = &NilValue{}

	// Create a new environment for the block's scope
	blockEnv := NewEnclosedEnvironment(env)

	for _, statement := range block.Statements {
		result = i.eval(statement, blockEnv)

		// If we hit a return statement, break out of the block
		if _, ok := result.(*ReturnValue); ok {
			return result
		}
	}

	return result
}

// evalIdentifier evaluates an identifier
func (i *Interpreter) evalIdentifier(node *ast.Identifier, env *Environment) Value {
	if val, ok := env.Get(node.Name); ok {
		// Found the variable, return its value
		return val
	}

	// Variable not found, return an error with more context
	return &ErrorValue{Message: fmt.Sprintf("identifier not found: %s", node.Name)}
}

// evalVariableDeclaration evaluates a variable declaration
func (i *Interpreter) evalVariableDeclaration(node *ast.VariableDecl, env *Environment) Value {
	var value Value
	if node.Value != nil {
		value = i.eval(node.Value, env)

		// Check if the value is an error
		if _, isError := value.(*ErrorValue); isError {
			return value
		}
	} else {
		// If no value is provided, initialize with nil
		value = &NilValue{}
	}

	if node.TypeAnnotation != nil {
		// Parse the type annotation
		varType := i.parseTypeAnnotation(node.TypeAnnotation)

		// Check that the value is compatible with the declared type
		if !types.IsAssignable(value.VibeType(), varType) {
			return &ErrorValue{Message: fmt.Sprintf("Type error: Cannot assign value of type %s to variable of type %s",
				value.VibeType().String(), varType.String())}
		}

		// Set with type check
		err := env.SetWithType(node.Name, value, varType)
		if err != nil {
			return &ErrorValue{Message: err.Error()}
		}
	} else {
		// No type annotation, infer from the value
		err := env.Set(node.Name, value)
		if err != nil {
			return &ErrorValue{Message: err.Error()}
		}
	}

	return value // Return the value instead of nil
}

// evalReturnStatement evaluates a return statement
func (i *Interpreter) evalReturnStatement(node *ast.ReturnStmt, env *Environment) Value {
	var value Value

	if node.Value != nil {
		value = i.eval(node.Value, env)
	} else {
		value = &NilValue{}
	}

	return &ReturnValue{Value: value}
}

// evalIfStatement evaluates an if statement
func (i *Interpreter) evalIfStatement(node *ast.IfStmt, env *Environment) Value {
	condition := i.eval(node.Condition, env)

	// Check if the condition is true
	if isTruthy(condition) {
		return i.eval(node.Consequence, env)
	}

	// Check else branch
	if node.Alternative != nil {
		return i.eval(node.Alternative, env)
	}

	return &NilValue{}
}

// evalForStatement evaluates a for statement
func (i *Interpreter) evalForStatement(node *ast.ForStmt, env *Environment) Value {
	// Create a new environment for the loop
	loopEnv := NewEnclosedEnvironment(env)

	// For loop with iterator and iterable
	// Evaluate the iterable
	iterable := i.eval(node.Iterable, env)

	// Handle iteration based on the type of iterable
	switch iter := iterable.(type) {
	case *ArrayValue:
		for _, element := range iter.Elements {
			// Bind the current element to the iterator variable
			err := loopEnv.Set(node.Iterator, element)
			if err != nil {
				return &ErrorValue{Message: err.Error()}
			}

			// Execute the body
			result := i.eval(node.Body, loopEnv)
			if returnValue, ok := result.(*ReturnValue); ok {
				return returnValue
			}
		}
	default:
		return &ErrorValue{Message: fmt.Sprintf("Type error: for loop requires an array to iterate over, got %s", iter.Type())}
	}

	return &NilValue{}
}

// evalClassDefinition evaluates a class definition
func (i *Interpreter) evalClassDefinition(node *ast.ClassDef, env *Environment) Value {
	// Create a new class value
	class := &ClassValue{
		Name:       node.Name,
		Properties: make(map[string]types.Type),
		Methods:    make(map[string]*FunctionValue),
	}

	// Add the class to the environment
	env.Set(node.Name, class)

	return &NilValue{}
}

// evalImportStatement evaluates an import statement
func (i *Interpreter) evalImportStatement(node *ast.RequireStmt, env *Environment) Value {
	// Not implemented yet
	return &NilValue{}
}

// evalPrefixExpression evaluates a prefix expression
func (i *Interpreter) evalPrefixExpression(node *ast.UnaryExpr, env *Environment) Value {
	right := i.eval(node.Right, env)

	switch node.Operator {
	case "-":
		// Handle numeric negation
		switch right := right.(type) {
		case *IntegerValue:
			return &IntegerValue{Value: -right.Value}
		case *FloatValue:
			return &FloatValue{Value: -right.Value}
		default:
			return &StringValue{Value: fmt.Sprintf("unknown operator: -%s", right.Type())}
		}
	case "!":
		// Handle logical negation
		return &BooleanValue{Value: !isTruthy(right)}
	default:
		return &StringValue{Value: fmt.Sprintf("unknown operator: %s%s", node.Operator, right.Type())}
	}
}

// evalInfixExpression evaluates an infix expression
func (i *Interpreter) evalInfixExpression(node *ast.BinaryExpr, env *Environment) Value {
	left := i.eval(node.Left, env)
	if _, isError := left.(*ErrorValue); isError {
		return left
	}

	right := i.eval(node.Right, env)
	if _, isError := right.(*ErrorValue); isError {
		return right
	}

	operator := node.Operator

	switch {
	case left.Type() == "INTEGER" && right.Type() == "INTEGER":
		return evalIntegerInfixExpression(operator, left, right)
	case left.Type() == "INTEGER" && right.Type() == "FLOAT" || left.Type() == "FLOAT" && right.Type() == "INTEGER" || left.Type() == "FLOAT" && right.Type() == "FLOAT":
		return evalNumberInfixExpression(operator, left, right)
	case left.Type() == "STRING" && right.Type() == "STRING":
		return evalStringInfixExpression(operator, left, right)
	case operator == "==":
		return &BooleanValue{Value: left.Inspect() == right.Inspect()}
	case operator == "!=":
		return &BooleanValue{Value: left.Inspect() != right.Inspect()}
	default:
		return &ErrorValue{Message: fmt.Sprintf("unknown operator: %s %s %s", left.Type(), operator, right.Type())}
	}
}

// evalIntegerInfixExpression evaluates an infix expression with integer operands
func evalIntegerInfixExpression(operator string, left, right Value) Value {
	leftVal := left.(*IntegerValue).Value
	rightVal := right.(*IntegerValue).Value

	switch operator {
	case "+":
		return &IntegerValue{Value: leftVal + rightVal}
	case "-":
		return &IntegerValue{Value: leftVal - rightVal}
	case "*":
		return &IntegerValue{Value: leftVal * rightVal}
	case "/":
		if rightVal == 0 {
			return &StringValue{Value: "Error: division by zero"}
		}
		return &IntegerValue{Value: leftVal / rightVal}
	case "<":
		return &BooleanValue{Value: leftVal < rightVal}
	case ">":
		return &BooleanValue{Value: leftVal > rightVal}
	case "<=":
		return &BooleanValue{Value: leftVal <= rightVal}
	case ">=":
		return &BooleanValue{Value: leftVal >= rightVal}
	case "==":
		return &BooleanValue{Value: leftVal == rightVal}
	case "!=":
		return &BooleanValue{Value: leftVal != rightVal}
	default:
		return &StringValue{Value: fmt.Sprintf("unknown operator: INTEGER %s INTEGER", operator)}
	}
}

// evalNumberInfixExpression evaluates an infix expression with numeric operands
func evalNumberInfixExpression(operator string, left, right Value) Value {
	var leftVal, rightVal float64

	// Convert left to float64
	if left.Type() == "INTEGER" {
		leftVal = float64(left.(*IntegerValue).Value)
	} else {
		leftVal = left.(*FloatValue).Value
	}

	// Convert right to float64
	if right.Type() == "INTEGER" {
		rightVal = float64(right.(*IntegerValue).Value)
	} else {
		rightVal = right.(*FloatValue).Value
	}

	switch operator {
	case "+":
		return &FloatValue{Value: leftVal + rightVal}
	case "-":
		return &FloatValue{Value: leftVal - rightVal}
	case "*":
		return &FloatValue{Value: leftVal * rightVal}
	case "/":
		if rightVal == 0 {
			return &StringValue{Value: "Error: division by zero"}
		}
		return &FloatValue{Value: leftVal / rightVal}
	case "<":
		return &BooleanValue{Value: leftVal < rightVal}
	case ">":
		return &BooleanValue{Value: leftVal > rightVal}
	case "<=":
		return &BooleanValue{Value: leftVal <= rightVal}
	case ">=":
		return &BooleanValue{Value: leftVal >= rightVal}
	case "==":
		return &BooleanValue{Value: leftVal == rightVal}
	case "!=":
		return &BooleanValue{Value: leftVal != rightVal}
	default:
		return &StringValue{Value: fmt.Sprintf("unknown operator: NUMBER %s NUMBER", operator)}
	}
}

// evalStringInfixExpression evaluates an infix expression with string operands
func evalStringInfixExpression(operator string, left, right Value) Value {
	leftVal := left.(*StringValue).Value
	rightVal := right.(*StringValue).Value

	switch operator {
	case "+":
		return &StringValue{Value: leftVal + rightVal}
	case "==":
		return &BooleanValue{Value: leftVal == rightVal}
	case "!=":
		return &BooleanValue{Value: leftVal != rightVal}
	default:
		return &StringValue{Value: fmt.Sprintf("unknown operator: STRING %s STRING", operator)}
	}
}

// evalAssignmentExpression evaluates an assignment expression
func (i *Interpreter) evalAssignmentExpression(node *ast.Assignment, env *Environment) Value {
	// Check if the value is nil
	if node.Value == nil {
		return &ErrorValue{Message: "Cannot evaluate nil node"}
	}

	value := i.eval(node.Value, env)

	// Check if the value is an error
	if _, isError := value.(*ErrorValue); isError {
		return value
	}

	// Set the value in the environment
	err := env.Set(node.Name, value)
	if err != nil {
		return &ErrorValue{Message: err.Error()}
	}

	return value
}

// evalArrayLiteral evaluates an array literal
func (i *Interpreter) evalArrayLiteral(node *ast.ArrayLiteral, env *Environment) Value {
	elements := make([]Value, 0, len(node.Elements))

	for _, element := range node.Elements {
		evaluated := i.eval(element, env)
		elements = append(elements, evaluated)
	}

	return &ArrayValue{Elements: elements}
}

// evalIndexExpression evaluates an index expression
func (i *Interpreter) evalIndexExpression(node *ast.IndexExpr, env *Environment) Value {
	array := i.eval(node.Array, env)
	index := i.eval(node.Index, env)

	// Check if we have a valid array
	if array.Type() != "ARRAY" {
		return &StringValue{Value: fmt.Sprintf("index operator not supported: %s", array.Type())}
	}

	// Check if we have a valid index
	if index.Type() != "INTEGER" {
		return &StringValue{Value: fmt.Sprintf("Type error: array index must be INTEGER, got %s", index.Type())}
	}

	arrayValue := array.(*ArrayValue)
	indexValue := index.(*IntegerValue)

	// Bounds check
	if indexValue.Value < 0 || indexValue.Value >= int64(len(arrayValue.Elements)) {
		return &NilValue{} // Return nil for out-of-bounds index
	}

	return arrayValue.Elements[indexValue.Value]
}

// evalMemberExpression evaluates a member expression
func (i *Interpreter) evalMemberExpression(node *ast.DotExpr, env *Environment) Value {
	// Not implemented yet
	return &NilValue{}
}

// evalMethodCallExpression evaluates a method call expression
func (i *Interpreter) evalMethodCallExpression(node *ast.MethodCall, env *Environment) Value {
	// Not implemented yet
	return &NilValue{}
}

// parseTypeAnnotation parses a type annotation
func (i *Interpreter) parseTypeAnnotation(node *ast.TypeAnnotation) types.Type {
	switch node.TypeName {
	case "int":
		return types.IntType
	case "float":
		return types.FloatType
	case "string":
		return types.StringType
	case "bool":
		return types.BoolType
	case "any":
		return types.AnyType
	case "Array":
		if node.GenericType != nil {
			elemType := i.parseTypeAnnotation(node.GenericType)
			return types.ArrayType{ElementType: elemType}
		} else if len(node.TypeParams) > 0 {
			if typeAnnotation, ok := node.TypeParams[0].(*ast.TypeAnnotation); ok {
				elemType := i.parseTypeAnnotation(typeAnnotation)
				return types.ArrayType{ElementType: elemType}
			}
		}
		// Default to Array of any
		return types.ArrayType{ElementType: types.AnyType}
	default:
		// Unknown type, default to any
		return types.AnyType
	}
}

// isTruthy returns true if the value is truthy
func isTruthy(obj Value) bool {
	switch obj := obj.(type) {
	case *BooleanValue:
		return obj.Value
	case *NilValue:
		return false
	case *IntegerValue:
		return obj.Value != 0
	case *FloatValue:
		return obj.Value != 0
	case *StringValue:
		return obj.Value != ""
	default:
		return true
	}
}

// ErrorValue represents an error value
type ErrorValue struct {
	Message string
}

func (e *ErrorValue) Type() string { return "ERROR" }
func (e *ErrorValue) Inspect() string { return "ERROR: " + e.Message }
func (e *ErrorValue) VibeType() types.Type { return types.AnyType }